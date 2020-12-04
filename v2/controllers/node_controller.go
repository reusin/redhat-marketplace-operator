// Copyright 2020 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	watchResourceTag   = "razee/watch-resource"
	watchResourceValue = "lite"
)

// ReconcileNode reconciles a Node object
type ReconcileNode struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func (r *ReconcileNode) SetupWithManager(mgr ctrl.Manager) error {
	// Create a new controller
	//
	labelPreds := []predicate.Predicate{
		predicate.Funcs{
			UpdateFunc: func(evt event.UpdateEvent) bool {
				watchResourceTag, ok := evt.MetaNew.GetLabels()[watchResourceTag]
				return !(ok && watchResourceTag == watchResourceValue)
			},
			CreateFunc: func(evt event.CreateEvent) bool {
				watchResourceTag, ok := evt.Meta.GetLabels()[watchResourceTag]
				return !(ok && watchResourceTag == watchResourceValue)
			},
			GenericFunc: func(evt event.GenericEvent) bool {
				watchResourceTag, ok := evt.Meta.GetLabels()[watchResourceTag]
				return !(ok && watchResourceTag == watchResourceValue)
			},
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		Watches(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestForObject{}, labelPreds...).
		Complete(r)
}

// Reconcile reads that state of the cluster for a Node object and makes changes based on the state read
// and what is in the Node.Spec
func (r *ReconcileNode) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)
	reqLogger.Info("Reconciling Node")

	// Fetch the Node instance
	instance := &corev1.Node{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Error(err, "node does not exist")
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get node")
		return reconcile.Result{}, err
	}

	labels := instance.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	nodeOriginalLabels := instance.DeepCopy().GetLabels()
	labels[watchResourceTag] = watchResourceValue
	if !reflect.DeepEqual(labels, nodeOriginalLabels) {
		instance.SetLabels(labels)
		if err := r.client.Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Failed to patch node with razee/watch-resource: lite label")
			return reconcile.Result{}, err
		}
		reqLogger.Info("Patched node with razee/watch-resource: lite label")
	} else {
		reqLogger.Info("No patch needed on node resource")
	}
	reqLogger.Info("reconcilation complete")
	return reconcile.Result{}, nil
}
