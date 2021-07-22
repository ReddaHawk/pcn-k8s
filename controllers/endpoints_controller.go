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
	"polycube.com/utils"
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

	msg := fmt.Sprintf("received reconcile request for %q (namespace: %q)", endpoint.GetName(), endpoint.GetNamespace())
	log.Log.Info(msg)
	// is object marked for deletion?
	if !endpoint.DeletionTimestamp.IsZero() {
		log.Log.Info("Pod marked for deletion")
		removeEndpoint(req.NamespacedName.String())
		return ctrl.Result{}, nil
	}

	addEndpoint(req.NamespacedName.String(),endpoint)
	return ctrl.Result{}, nil
}

func removeEndpoint(namespacedName string) {
	
}

func addEndpoint(namespacedName string, endpoint *corev1.Endpoints) {
	if _, ok := utils.Services[namespacedName]; ok {

	}
}

func updateEndpoint() {

}

// SetupWithManager sets up the controller with the Manager.
func (r *EndpointsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Endpoints{}).
		Complete(r)
}
