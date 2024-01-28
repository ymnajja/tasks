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

//+kubebuilder:rbac:groups=app.example.com,resources=podapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.example.com,resources=podapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.example.com,resources=podapps/finalizers,verbs=update

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
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, podApp)
	l.Info("ENTEr reconcile", "Spec", podApp.Spec, "Status", podApp.Status)
	r.reconcilePod(ctx, podApp, l)
	l.Info("ENTEr reconcile", "Spec", podApp.Spec, "Status", podApp.Status)
	r.reconcileservice(ctx, podApp, l)
	if podApp.Spec.Enable {
		servicename := podApp.Spec.PodName + "-Is-not-fun"
		l.Info("Reconcile complete", "flag", servicename)
	}
	return ctrl.Result{}, nil
}
func (r *PodAppReconciler) reconcilePod(ctx context.Context, PodApp *appv1.PodApp, l logr.Logger) error {

	pod := &corev1.Pod{}

	err := r.Get(ctx, types.NamespacedName{Name: PodApp.Spec.PodName, Namespace: PodApp.Spec.PodNamespace}, pod)

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
			Name:      PodApp.Spec.PodName,
			Namespace: PodApp.Spec.PodNamespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  PodApp.Spec.PodSpec.Containers[0].Name,
					Image: PodApp.Spec.PodSpec.Containers[0].Image,
				},
			},
		},
	}

	l.Info("Creating pod")
	return r.Create(ctx, pod)
}

func (r *PodAppReconciler) reconcileservice(ctx context.Context, PodApp *appv1.PodApp, l logr.Logger) error {
	Service := &corev1.Service{}

	err := r.Get(ctx, types.NamespacedName{Name: PodApp.Spec.PodName, Namespace: PodApp.Spec.PodNamespace}, Service)

	if err == nil {
		l.Info("service Found")
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	l.Info("service Not found")

	Service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PodApp.Spec.PodName,
			Namespace: PodApp.Spec.PodNamespace,
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{{Port: 80, NodePort: 30080, Protocol: corev1.ProtocolTCP}},
		},
	}

	l.Info("Creating service")
	return r.Create(ctx, Service)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.PodApp{}).
		Complete(r)
}
