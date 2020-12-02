package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	consulv1alpha1 "github.com/ashwin-venkatesh/consul-operator/api/v1alpha1"
)

// OperatorReconciler reconciles a Operator object
type OperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=consul.hashicorp.com,resources=operators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=consul.hashicorp.com,resources=operators/status,verbs=get;update;patch

func (r *OperatorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	operator := &consulv1alpha1.Operator{}
	ctx := context.Background()
	log := r.Log.WithValues("operator", req.NamespacedName)

	err := r.Get(ctx, req.NamespacedName, operator)
	if err != nil {
		log.Error(err, "failed getting operator")
		return ctrl.Result{}, err
	}

	if operator.Spec.Global.Enabled || operator.Spec.Server.Enabled {
		serverServiceAccount := &corev1.ServiceAccount{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: operator.Namespace,
			Name:      fmt.Sprintf("%s-server", operator.Name),
		}, serverServiceAccount)

		if err != nil && errors.IsNotFound(err) {
			serverServiceAccount = operator.ServerServiceAccount()
			if err := r.Create(ctx, serverServiceAccount); err != nil {
				log.Error(err, "failed creating server serviceaccount")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "failed getting server serviceaccount")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}

		serverRole := &rbacv1.Role{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: operator.Namespace,
			Name:      fmt.Sprintf("%s-server", operator.Name),
		}, serverRole)

		if err != nil && errors.IsNotFound(err) {
			serverRole = operator.ServerRole()
			if err := r.Create(ctx, serverRole); err != nil {
				log.Error(err, "failed creating server role")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "failed getting server role")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}

		serverRoleBinding := &rbacv1.RoleBinding{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: operator.Namespace,
			Name:      fmt.Sprintf("%s-server", operator.Name),
		}, serverRoleBinding)

		if err != nil && errors.IsNotFound(err) {
			serverRoleBinding = operator.ServerRoleBinding()
			if err := r.Create(ctx, serverRoleBinding); err != nil {
				log.Error(err, "failed creating server rolebinding")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "failed getting server rolebinding")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}

		serverConfigMap := &corev1.ConfigMap{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: operator.Namespace,
			Name:      fmt.Sprintf("%s-server", operator.Name),
		}, serverConfigMap)

		if err != nil && errors.IsNotFound(err) {
			serverConfigMap = operator.ServerConfigMap()
			if err := r.Create(ctx, serverConfigMap); err != nil {
				log.Error(err, "failed creating server configmap")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "failed getting server configmap")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}

		serverService := &corev1.Service{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: operator.Namespace,
			Name:      fmt.Sprintf("%s-server", operator.Name),
		}, serverService)

		if err != nil && errors.IsNotFound(err) {
			serverService = operator.ServerService()
			if err := r.Create(ctx, serverService); err != nil {
				log.Error(err, "failed creating server service")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "failed getting server service")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}

		serverStatefulSet := &appsv1.StatefulSet{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: operator.Namespace,
			Name:      fmt.Sprintf("%s-server", operator.Name),
		}, serverStatefulSet)

		if err != nil && errors.IsNotFound(err) {
			serverStatefulSet = operator.ServerStatefulSet()
			if err := r.Create(ctx, serverStatefulSet); err != nil {
				log.Error(err, "failed creating server statefulset")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "failed getting server statefulset")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}
	}

	uiService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{
		Namespace: operator.Namespace,
		Name:      fmt.Sprintf("%s-ui", operator.Name),
	}, uiService)

	if err != nil && errors.IsNotFound(err) {
		uiService = operator.UIService()
		if err := r.Create(ctx, uiService); err != nil {
			log.Error(err, "failed creating ui service")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "failed getting ui service")
		return ctrl.Result{}, err
	} else {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *OperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&consulv1alpha1.Operator{}).
		Complete(r)
}
