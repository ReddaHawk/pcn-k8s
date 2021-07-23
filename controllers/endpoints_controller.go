/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EndpointsReconciler reconciles a Endpoints object
type EndpointsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=endpoints/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=endpoints/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Endpoints object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *EndpointsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	endpoint := &corev1.Endpoints{}
	if err := r.Client.Get(ctx, req.NamespacedName, endpoint); err != nil {
		// add some debug information if it's not a NotFound error

		removeEndpoint(req.NamespacedName.String())

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	msg := fmt.Sprintf("received reconcile request for Endpoint %q (namespace: %q)", endpoint.GetName(), endpoint.GetNamespace())
	log.Log.Info(msg)
	// is object marked for deletion?
	if !endpoint.DeletionTimestamp.IsZero() {
		log.Log.Info("Endpoint marked for deletion")
		//removeEndpoint(req.NamespacedName.String())
		return ctrl.Result{}, nil
	}

	r.addService(req.NamespacedName.String(),endpoint)
	return ctrl.Result{}, nil
}

func removeEndpoint(namespacedName string) {
	
}

func (r* EndpointsReconciler) addService(namespacedName string, endpoint *corev1.Endpoints) {

	if _, ok := Services[namespacedName]; ok {
		updateService(namespacedName, endpoint)
		return
	}
	service, err := r.parseService(endpoint)
	if err != nil {
		return
	}
	Services[namespacedName] = service
	printService(service)
	addNodeService(service)

}

func (r* EndpointsReconciler) updateService(namespacedName string, endpoint *corev1.Endpoints) {
	if _, ok := Services[namespacedName]; !ok {
		r.addService(namespacedName,endpoint)
		return
	}
	log.Log.Info("Updating service "+endpoint.Name)
	serviceOld := Services[namespacedName]
	serviceNew, err := r.parseService(endpoint)

	if err != nil {
		return
	}

	// get list of added and removed ports
	deletedPorts, addedPorts := getServicePortsDiff(serviceOld, serviceNew)
	if len(deletedPorts) > 0 {
		log.Debugf("ports deleted:")
	}
	for _, i := range deletedPorts {
		log.Debugf("--%d:%s", i.Port, i.Proto)
		delNodeServicePort(serviceNew, i)
	}

	if len(addedPorts) > 0 {
		log.Debugf("ports added:")
	}
	for _, i := range addedPorts {
		log.Debugf("--%d:%s", i.Port, i.Proto)
		addNodeServicePort(serviceNew, i)
	}

	// for all ports, get list update backends
	for portOldKey, portOldValue := range serviceOld.Ports {
		portNewValue, ok := serviceNew.Ports[portOldKey]
		if !ok {
			continue
		}
		deletedBackends, addedBackends := getServicePortBackendsDiff(portOldValue, portNewValue)
		if len(deletedBackends) > 0 {
			log.Log.Info("addressed deleted:")
		}
		for _, i := range deletedBackends {
			log.Log.Info(fmt.Sprintf("--%s:%d", i.IP, i.Port))
			delNodeServicePortBackend(serviceNew, portOldValue, i)
		}

		if len(addedBackends) > 0 {
			log.Log.Info("addressed added:")
		}
		for _, i := range addedBackends {
			log.Log.Info(fmt.Sprintf("--%s:%d", i.IP, i.Port))
			addNodeServicePortBackend(serviceNew, portOldValue, i)
		}
	}

	services[uid] = serviceNew
}


func (r* EndpointsReconciler) parseService(endpoint *corev1.Endpoints) (Service, error){
	uid := endpoint.UID
	name := endpoint.Name
	namespace := endpoint.Namespace
	serviceName := types.NamespacedName{Namespace: namespace,Name: name}
	svc := &corev1.Service{}
	if err := r.Client.Get(context.Background(), serviceName, svc); err != nil {
		log.Log.Error(err,"Error getting svc ")
		return Service{}, err
	}
	type_ := svc.Spec.Type
	vip := svc.Spec.ClusterIP
	externalTrafficPolicy := svc.Spec.ExternalTrafficPolicy

	servicePorts := make(map[ServicePortKey]ServicePort)
	// parse different ports
	for _, port := range svc.Spec.Ports {
		port_ := ServicePort{Name: port.Name,
			Port:     port.Port,
			Nodeport: port.NodePort,
			Proto:    string(port.Protocol)}

		port_.Backends = make(map[Backend]bool)

		addresses := make(map[string]bool)
		endpoindPorts := make(map[int32]bool)

		// get endpoints that implement this port
		for _, subset := range endpoint.Subsets {
			// get IP addresses
			for _, addr := range subset.Addresses {
				addresses[addr.IP] = true
			}

			// get ports
			for _, endpointPort := range subset.Ports {
				if port.Name == "" || port.Name == endpointPort.Name {
					endpoindPorts[endpointPort.Port] = true
				}
			}
		}

		// perform address x ports (cartesian product)
		for address := range addresses {
			for endpointPort := range endpoindPorts {
				port_.Backends[Backend{address, endpointPort}] = true
			}
		}

		servicePorts[ServicePortKey{port_.Port, port_.Proto}] = port_

	}
	return Service{UID: uid,
		Name: name,
		VIP: vip,
		Type: string(type_),
		ExternalTrafficPolicy: string(externalTrafficPolicy),
		Ports: servicePorts,
	} , nil
}
func printService(entry Service) {
	log.Log.Info(fmt.Sprintf("uid: %s", entry.UID))
	log.Log.Info(fmt.Sprintf("name: %s", entry.Name))
	log.Log.Info(fmt.Sprintf("vip: %s", entry.VIP))
	log.Log.Info(fmt.Sprintf("ports: "))
	for _, y := range entry.Ports {
		log.Log.Info(fmt.Sprintf("--%s, %d, %s, %d", y.Name, y.Port, y.Proto, y.Nodeport))
		log.Log.Info(fmt.Sprintf("backends: "))
		for x := range y.Backends {
			log.Log.Info(fmt.Sprintf("--%s:%d", x.IP, x.Port))
		}
	}
}


// SetupWithManager sets up the controller with the Manager.
func (r *EndpointsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Endpoints{}).
		Complete(r)
}
