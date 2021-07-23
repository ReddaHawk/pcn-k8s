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

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	pod := &corev1.Pod{}
	if err := r.Client.Get(ctx, req.NamespacedName, pod); err != nil {
		// add some debug information if it's not a NotFound error

		log.Log.Info("unable to fetch VmGroup")

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	msg := fmt.Sprintf("received reconcile request for Pod %q (namespace: %q)", pod.GetName(), pod.GetNamespace())
	log.Log.Info(msg)
	// is object marked for deletion?
	if !pod.DeletionTimestamp.IsZero() {
		log.Log.Info("Pod marked for deletion" + pod.Status.PodIP)
		removePod(req.NamespacedName.String())
		return ctrl.Result{}, nil
	}
	// If pod is not in running return
	if pod.Status.Phase != corev1.PodRunning {
		log.Log.Info("Pod is not in running: " + req.NamespacedName.String())

		return ctrl.Result{},nil
	}

	addPod(req.NamespacedName.String(),pod)
	return ctrl.Result{}, nil
}

func addPod(namespacedName string, pod *corev1.Pod) {
	if _, ok := utils.Pods[namespacedName]; ok {
		log.Log.Info("Pod already in the list: " + namespacedName)
		return
	}
	// todo aggiungi rotte esistenti al lbrp

	utils.Pods[namespacedName] = parsePod(pod)
	log.Log.Info("Added pod in the list: " + namespacedName + "lenght "  + fmt.Sprintf("%d",len(utils.Pods)))
}


func parsePod(p *corev1.Pod) utils.Pod {
	// UID is unique within the whole system
	return utils.Pod{p.UID, p.Name,p.Status.PodIP }
}

// todo complete
func removePod(namespacedName string)  {

	if _, ok := utils.Pods[namespacedName]; !ok {
		log.Log.Info("Removing Pod: "+namespacedName)
		delete(utils.Pods,namespacedName)
		return
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
