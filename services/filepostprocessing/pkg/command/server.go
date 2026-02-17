package command

import (
	"github.com/opencloud-eu/opencloud/pkg/config/configlog"
	"github.com/opencloud-eu/opencloud/pkg/generators"
	"github.com/opencloud-eu/opencloud/pkg/log"
	"github.com/opencloud-eu/opencloud/services/filepostprocessing/pkg/config"
	"github.com/opencloud-eu/opencloud/services/filepostprocessing/pkg/config/parser"
	"github.com/opencloud-eu/opencloud/services/filepostprocessing/pkg/service"
	"github.com/opencloud-eu/reva/v2/pkg/events/stream"
	"github.com/spf13/cobra"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start filepostprocessing service",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.Configure(cfg.Service.Name, cfg.Commons, cfg.LogLevel)

			logger.Info().
				Str("endpoint", cfg.Events.Endpoint).
				Str("cluster", cfg.Events.Cluster).
				Msg("Initializing NATS stream")

			connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
			bus, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				logger.Error().Err(err).Msg("Failed to initialize NATS stream")
				return err
			}

			svc, err := service.NewFilePostProcessingService(logger, cfg, bus)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create filepostprocessing service")
				return err
			}

			return svc.Run()
		},
	}
}
