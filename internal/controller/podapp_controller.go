/*
Copyright 2024.

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

package controller

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "test/api/v1"
)

// Definitions to manage status conditions

// PodAppReconciler reconciles a PodApp object
type PodAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.example.com,resources=podapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.example.com,resources=podapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.example.com,resources=podapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=resources=pods,verbs=list;watch;create;
// +kubebuilder:rbac:groups=resources=services,verbs=list;watch;create;
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PodAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("enter reconcile", "req", req)
	podApp := &appv1.PodApp{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: "rye"}, podApp)
	l.Info("ENTEr reconcile", "Spec", podApp.Spec, "Status", podApp.Status)
	r.reconcilePod(ctx, podApp, l)
	l.Info("ENTEr reconcile", "Spec", podApp.Spec, "Status", podApp.Status)
	if podApp.Spec.Enable {
		servicename := "gin-Is-not-weak-after-all"
		l.Info("Reconcile complete", "flag", servicename)
	}
	return ctrl.Result{}, nil
}
func (r *PodAppReconciler) reconcilePod(ctx context.Context, PodApp *appv1.PodApp, l logr.Logger) error {

	pod := &corev1.Pod{}
	podname := "gin-pod"
	err := r.Get(ctx, types.NamespacedName{Name: podname, Namespace: "gin"}, pod)

	if err == nil {
		l.Info("Pod Found")
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	l.Info("Pod Not found")

	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podname,
			Namespace: "gin",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "gin-pod",
					Image: "nginx",
				},
			},
		},
	}

	l.Info("Creating pod")
	return r.Create(ctx, pod)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.PodApp{}).
		Complete(r)
}
