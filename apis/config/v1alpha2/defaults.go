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

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

const (
	DefaultNamespace              = "kueue-system"
	DefaultWebhookServiceName     = "kueue-webhook-service"
	DefaultWebhookSecretName      = "kueue-webhook-server-cert"
	DefaultWebhookPort            = 9443
	DefaultHealthProbeBindAddress = ":8081"
	DefaultMetricsBindAddress     = ":8080"
	DefaultLeaderElectionID       = "c1f6bfd2.kueue.x-k8s.io"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&Configuration{}, func(obj interface{}) {
		SetDefaults_Configuration(obj.(*Configuration))
	})
	return nil
}

// SetDefaults_Configuration sets default values for ComponentConfig.
func SetDefaults_Configuration(cfg *Configuration) {
	if cfg.Namespace == nil {
		cfg.Namespace = pointer.String(DefaultNamespace)
	}
	if cfg.Webhook.Port == nil {
		cfg.Webhook.Port = pointer.Int(DefaultWebhookPort)
	}
	if len(cfg.Metrics.BindAddress) == 0 {
		cfg.Metrics.BindAddress = DefaultMetricsBindAddress
	}
	if len(cfg.Health.HealthProbeBindAddress) == 0 {
		cfg.Health.HealthProbeBindAddress = DefaultHealthProbeBindAddress
	}
	if cfg.LeaderElection != nil && cfg.LeaderElection.LeaderElect != nil &&
		*cfg.LeaderElection.LeaderElect && len(cfg.LeaderElection.ResourceName) == 0 {
		cfg.LeaderElection.ResourceName = DefaultLeaderElectionID
	}
	if cfg.InternalCertManagement == nil {
		cfg.InternalCertManagement = &InternalCertManagement{}
	}
	if cfg.InternalCertManagement.Enable == nil {
		cfg.InternalCertManagement.Enable = pointer.Bool(true)
	}
	if *cfg.InternalCertManagement.Enable {
		if cfg.InternalCertManagement.WebhookServiceName == nil {
			cfg.InternalCertManagement.WebhookServiceName = pointer.String(DefaultWebhookServiceName)
		}
		if cfg.InternalCertManagement.WebhookSecretName == nil {
			cfg.InternalCertManagement.WebhookSecretName = pointer.String(DefaultWebhookSecretName)
		}
	}

	if cfg.Performance == nil {
		cfg.Performance = &Performance{}
	}
	SetDefaults_Performance(cfg.Performance)
}

// SetDefaults_Performance sets default values for Performance config.
func SetDefaults_Performance(performance *Performance) {
	if performance.JobControllerWorkerCount == nil {
		performance.JobControllerWorkerCount = pointer.Int(1)
	}
	if performance.ResourceFlavorControllerWorkerCount == nil {
		performance.ResourceFlavorControllerWorkerCount = pointer.Int(1)
	}
	if performance.ClusterQueueControllerWorkerCount == nil {
		performance.ClusterQueueControllerWorkerCount = pointer.Int(1)
	}
	if performance.LocalQueueControllerWorkerCount == nil {
		performance.LocalQueueControllerWorkerCount = pointer.Int(1)
	}
	if performance.WorkloadControllerWorkerCount == nil {
		performance.WorkloadControllerWorkerCount = pointer.Int(1)
	}
}
