/*
Copyright 2022 The Kubernetes Authors.

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

package core

import (
	ctrl "sigs.k8s.io/controller-runtime"

	config "sigs.k8s.io/kueue/apis/config/v1alpha2"
	"sigs.k8s.io/kueue/pkg/cache"
	"sigs.k8s.io/kueue/pkg/queue"
)

const updateChBuffer = 10

// SetupControllers sets up the core controllers. It returns the name of the
// controller that failed to create and an error, if any.
func SetupControllers(mgr ctrl.Manager, qManager *queue.Manager, cc *cache.Cache, performance *config.Performance) (string, error) {
	rfRec := NewResourceFlavorReconciler(mgr.GetClient(), qManager, cc, *performance.ResourceFlavorControllerWorkerCount)
	if err := rfRec.SetupWithManager(mgr); err != nil {
		return "ResourceFlavor", err
	}
	qRec := NewLocalQueueReconciler(mgr.GetClient(), qManager, cc, *performance.LocalQueueControllerWorkerCount)
	if err := qRec.SetupWithManager(mgr); err != nil {
		return "LocalQueue", err
	}
	cqRec := NewClusterQueueReconciler(mgr.GetClient(), qManager, cc, *performance.ClusterQueueControllerWorkerCount, rfRec)
	rfRec.AddUpdateWatcher(cqRec)
	if err := cqRec.SetupWithManager(mgr); err != nil {
		return "ClusterQueue", err
	}
	if err := NewWorkloadReconciler(mgr.GetClient(), qManager, cc,
		*performance.WorkloadControllerWorkerCount, qRec, cqRec).SetupWithManager(mgr); err != nil {
		return "Workload", err
	}
	return "", nil
}
