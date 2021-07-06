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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"

	corev1 "k8s.io/api/core/v1"

	hotdrinkv1alpha1 "github.com/kube-go/drink-operator/api/v1alpha1"
)

var log = ctrl.Log.WithName("controllers").WithName("Pod")

// TeaReconciler reconciles a Tea object
type TeaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hotdrink.slashpai.github.io,resources=teas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hotdrink.slashpai.github.io,resources=teas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hotdrink.slashpai.github.io,resources=teas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tea object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *TeaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// your logic here
	reqLogger := log.WithValues("namespace", req.Namespace, "Tea", req.Name)
	reqLogger.Info("=== Reconciling At")

	instance := &hotdrinkv1alpha1.Tea{}

	err := r.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after
			// reconcile request—return and don't requeue:
			return reconcile.Result{}, nil
		}
		// Error reading the object—requeue the request:
		return reconcile.Result{}, err
	}

	// If no phase set, default to pending (the initial phase):
	if instance.Status.Phase == "" {
		reqLogger.Info("Phase empty in Status")
		instance.Status.Phase = hotdrinkv1alpha1.PhasePending
	}

	// Now let's make the main case distinction: implementing
	// the state diagram PENDING -> RUNNING -> DONE

	switch instance.Status.Phase {
	case hotdrinkv1alpha1.PhasePending:
		reqLogger.Info("Phase: PENDING")
		reqLogger.Info("Checking recipe", "Target", instance.Spec.Recipe)
		instance.Status.Phase = hotdrinkv1alpha1.PhaseRunning
	case hotdrinkv1alpha1.PhaseRunning:
		reqLogger.Info("Phase: Running")
		pod := newPodForCR(instance)
		// Set Tea instance as the owner and controller
		err := controllerutil.SetControllerReference(instance, pod, r.Scheme)
		if err != nil {
			// requeue with error
			return reconcile.Result{}, err
		}
		found := &corev1.Pod{}
		nsName := types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}
		err = r.Get(context.TODO(), nsName, found)
		// Try to see if the pod already exists and if not
		// (which we expect) then create a one-shot pod as per spec:
		if err != nil && errors.IsNotFound(err) {
			err = r.Create(context.TODO(), pod)
			if err != nil {
				// requeue with error
				return reconcile.Result{}, err
			}
			reqLogger.Info("Pod launched", "name", pod.Name)
		} else if err != nil {
			// requeue with error
			return reconcile.Result{}, err
		} else if found.Status.Phase == corev1.PodFailed ||
			found.Status.Phase == corev1.PodSucceeded {
			reqLogger.Info("Container terminated", "reason",
				found.Status.Reason, "message", found.Status.Message)
			instance.Status.Phase = hotdrinkv1alpha1.PhaseDone
		} else {
			// Don't requeue because it will happen automatically when the
			// pod status changes.
			return reconcile.Result{}, nil
		}
	case hotdrinkv1alpha1.PhaseDone:
		reqLogger.Info("Phase: DONE")
		return reconcile.Result{}, nil
	default:
		reqLogger.Info("NOP")
		return reconcile.Result{}, nil
	}

	// Update the tea instance, setting the status to the respective phase:
	err = r.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Don't requeue. We should be reconcile because either the pod
	// or the CR changes.
	return reconcile.Result{}, nil

}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *hotdrinkv1alpha1.Tea) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: strings.Split(cr.Spec.Recipe, " "),
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hotdrinkv1alpha1.Tea{}).
		Complete(r)
}
