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

// Package admission provides admission webhook configuration structures used for command execution
package admission

import (
	"context"
	"fmt"

	ionosapisinstall "github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/install"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos"
	"github.com/gardener/gardener/extensions/pkg/controller/cmd"
	"github.com/gardener/gardener/extensions/pkg/util"
	webhookcmd "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
	"github.com/gardener/gardener/pkg/apis/core/install"
	"github.com/spf13/cobra"
	"k8s.io/component-base/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var logger = log.Log.WithName("gardener-extension-admission-ionos")

// NewAdmissionCommand creates a new command for running an IONOS admission webhook.
func NewAdmissionCommand(ctx context.Context) *cobra.Command {
	var (
		restOpts = &cmd.RESTOptions{}

		mgrOpts  = &cmd.ManagerOptions{
			WebhookServerPort: 443,
		}

		webhookSwitches = webhookSwitchOptions()
		webhookOptions  = webhookcmd.NewAddToManagerSimpleOptions(webhookSwitches)

		aggOption = cmd.NewOptionAggregator(
			restOpts,
			mgrOpts,
			webhookOptions,
		)

		opts = manager.Options{}
	)

	cmdDefinition := &cobra.Command{
		Use: fmt.Sprintf("admission-%s", ionos.Type),

		RunE: func(cmdDefinition *cobra.Command, args []string) error {
			if err := aggOption.Complete(); err != nil {
				cmd.LogErrAndExit(err, "Error completing options")
			}

			util.ApplyClientConnectionConfigurationToRESTConfig(&config.ClientConnectionConfiguration{
				QPS:   100.0,
				Burst: 130,
			}, restOpts.Completed().Config)

			mgrOptions := mgrOpts.Completed().Options()
			mgrOptions.HealthProbeBindAddress = opts.HealthProbeBindAddress

			mgr, err := manager.New(restOpts.Completed().Config, mgrOptions)
			if err != nil {
				cmd.LogErrAndExit(err, "Could not instantiate manager")
			}

			install.Install(mgr.GetScheme())

			if err := ionosapisinstall.AddToScheme(mgr.GetScheme()); err != nil {
				cmd.LogErrAndExit(err, "Could not update manager scheme")
			}

			logger.Info("Setting up healthcheck endpoints")
			if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
				return err
			}

			logger.Info("Setting up webhook server")
			if err := webhookOptions.Completed().AddToManager(mgr); err != nil {
				return err
			}

			logger.Info("Setting up readycheck for webhook server")
			if err := mgr.AddReadyzCheck("webhook-server", mgr.GetWebhookServer().StartedChecker()); err != nil {
				return err
			}

			return mgr.Start(ctx)
		},
	}

	flags := cmdDefinition.Flags()
	aggOption.AddFlags(flags)

	flags.StringVar(&opts.HealthProbeBindAddress, "health-bind-address", ":8081", "bind address for the health server")

	return cmdDefinition
}
