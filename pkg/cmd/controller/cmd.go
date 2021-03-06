/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved.

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

// Package controller provides Kubernetes controller configuration structures used for command execution
package controller

import (
	"context"
	"fmt"
	"os"

	ionoscontrolplane "github.com/23technologies/gardener-extension-provider-ionos/pkg/controller/controlplane"
	ionoshealthcheck "github.com/23technologies/gardener-extension-provider-ionos/pkg/controller/healthcheck"
	ionosinfrastructure "github.com/23technologies/gardener-extension-provider-ionos/pkg/controller/infrastructure"
	ionosworker "github.com/23technologies/gardener-extension-provider-ionos/pkg/controller/worker"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos"
	ionosapisinstall "github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/install"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/cmd"
    "github.com/gardener/gardener/extensions/pkg/controller/controlplane/genericactuator"
	"github.com/gardener/gardener/extensions/pkg/controller/worker"
	"github.com/gardener/gardener/extensions/pkg/util"
	webhookcmd "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"
	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"k8s.io/component-base/version/verflag"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// NewControllerManagerCommand creates a new command for running a IONOS provider controller.
//
// PARAMETERS
// ctx context.Context Execution context
func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	generalOpts := &cmd.GeneralOptions{}
	restOpts    := &cmd.RESTOptions{}

	mgrOpts  := &cmd.ManagerOptions{
		LeaderElection:          true,
		LeaderElectionID:        cmd.LeaderElectionNameID(ionos.Name),
		LeaderElectionNamespace: os.Getenv("LEADER_ELECTION_NAMESPACE"),
		WebhookServerPort:       443,
	}

	configFileOpts := &ConfigOptions{}

	// options for the infrastructure controller
	infraCtrlOpts := &cmd.ControllerOptions{
		MaxConcurrentReconciles: 5,
	}
	reconcileOpts := &cmd.ReconcilerOptions{}

	// options for the health care controller
	healthCareCtrlOpts := &cmd.ControllerOptions{
		MaxConcurrentReconciles: 5,
	}

	// options for the control plane controller
	controlPlaneCtrlOpts := &cmd.ControllerOptions{
		MaxConcurrentReconciles: 5,
	}

	// options for the worker controller
	workerCtrlOpts := &cmd.ControllerOptions{
		MaxConcurrentReconciles: 5,
	}
	workerReconcileOpts := &worker.Options{
		DeployCRDs: true,
	}
	workerCtrlOptsUnprefixed := cmd.NewOptionAggregator(workerCtrlOpts, workerReconcileOpts)

	// options for the webhook server
	webhookServerOptions := &webhookcmd.ServerOptions{
		Namespace: os.Getenv("WEBHOOK_CONFIG_NAMESPACE"),
	}

	controllerSwitches := controllerSwitchOptions()
	webhookSwitches    := webhookSwitchOptions()

    webhookOptions     := webhookcmd.NewAddToManagerOptions(ionos.Name,
		genericactuator.ShootWebhooksResourceName,
		genericactuator.ShootWebhookNamespaceSelector(ionos.Type),
		webhookServerOptions,
		webhookSwitches,
	)

	aggOption := cmd.NewOptionAggregator(
		generalOpts,
		restOpts,
		mgrOpts,
		cmd.PrefixOption("controlplane-", controlPlaneCtrlOpts),
		cmd.PrefixOption("infrastructure-", infraCtrlOpts),
		cmd.PrefixOption("worker-", &workerCtrlOptsUnprefixed),
		cmd.PrefixOption("healthcheck-", healthCareCtrlOpts),
		controllerSwitches,
		configFileOpts,
		reconcileOpts,
		webhookOptions,
	)

	cmdDefinition := &cobra.Command{
		Use: fmt.Sprintf("%s-controller-manager", ionos.Name),

		PreRun: func(cmdDefinition *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()
		},

		RunE: func(cmdDefinition *cobra.Command, args []string) error {
			if err := aggOption.Complete(); err != nil {
				return fmt.Errorf("Error completing options: %w", err)
			}

			util.ApplyClientConnectionConfigurationToRESTConfig(configFileOpts.Completed().Config.ClientConnection, restOpts.Completed().Config)

			if workerReconcileOpts.Completed().DeployCRDs {
				if err := worker.ApplyMachineResourcesForConfig(ctx, restOpts.Completed().Config); err != nil {
					return fmt.Errorf("Error ensuring the machine CRDs: %w", err)
				}
			}

			mgrOptions := mgrOpts.Completed().Options()

			configFileOpts.Completed().ApplyMetricsBindAddress(&mgrOptions.MetricsBindAddress)

			mgr, err := manager.New(restOpts.Completed().Config, mgrOptions)
			if err != nil {
				return fmt.Errorf("Could not instantiate manager: %w", err)
			}

			scheme := mgr.GetScheme()
			if err := controller.AddToScheme(scheme); err != nil {
				return fmt.Errorf("Could not update manager scheme: %w", err)
			}
			if err := ionosapisinstall.AddToScheme(scheme); err != nil {
				return fmt.Errorf("Could not update manager scheme: %w", err)
			}
			if err := druidv1alpha1.AddToScheme(scheme); err != nil {
				return fmt.Errorf("Could not update manager scheme: %w", err)
			}
			if err := machinev1alpha1.AddToScheme(scheme); err != nil {
				return fmt.Errorf("Could not update manager scheme: %w", err)
			}
			if err := autoscalingv1beta2.AddToScheme(scheme); err != nil {
				return fmt.Errorf("Could not update manager scheme: %w", err)
			}

			// add common meta types to schema for controller-runtime to use v1.ListOptions
			metav1.AddToGroupVersion(scheme, machinev1alpha1.SchemeGroupVersion)

			configFileOpts.Completed().ApplyGardenId(&ionoscontrolplane.DefaultAddOptions.GardenId)
			configFileOpts.Completed().ApplyGardenId(&ionosinfrastructure.DefaultAddOptions.GardenId)
			configFileOpts.Completed().ApplyHealthCheckConfig(&ionoshealthcheck.DefaultAddOptions.HealthCheckConfig)
			healthCareCtrlOpts.Completed().Apply(&ionoshealthcheck.DefaultAddOptions.Controller)
			controlPlaneCtrlOpts.Completed().Apply(&ionoscontrolplane.DefaultAddOptions.Controller)
			infraCtrlOpts.Completed().Apply(&ionosinfrastructure.DefaultAddOptions.Controller)
			reconcileOpts.Completed().Apply(&ionosinfrastructure.DefaultAddOptions.IgnoreOperationAnnotation)
			reconcileOpts.Completed().Apply(&ionoscontrolplane.DefaultAddOptions.IgnoreOperationAnnotation)
			reconcileOpts.Completed().Apply(&ionosworker.DefaultAddOptions.IgnoreOperationAnnotation)
			workerCtrlOpts.Completed().Apply(&ionosworker.DefaultAddOptions.Controller)

		    if _, err := webhookOptions.Completed().AddToManager(ctx, mgr); err != nil {
				return fmt.Errorf("Could not add webhooks to manager: %w", err)
			}

			if err := controllerSwitches.Completed().AddToManager(mgr); err != nil {
				return fmt.Errorf("Could not add controllers to manager: %w", err)
			}

			if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
				return err
			}

			if err := mgr.AddReadyzCheck("informer-sync", gardenerhealthz.NewCacheSyncHealthz(mgr.GetCache())); err != nil {
				return fmt.Errorf("could not add readycheck for informers: %w", err)
			}

			if err := mgr.AddReadyzCheck("webhook-server", mgr.GetWebhookServer().StartedChecker()); err != nil {
				return fmt.Errorf("could not add readycheck of webhook to manager: %w", err)
			}

			if err := mgr.Start(ctx); err != nil {
				return fmt.Errorf("Error running manager: %w", err)
			}

			return nil
		},
	}

	cmdFlags := cmdDefinition.Flags()
	aggOption.AddFlags(cmdFlags)
	verflag.AddFlags(cmdFlags)

	return cmdDefinition
}
