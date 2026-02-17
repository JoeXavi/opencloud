package command

import (
	"os"

	"github.com/opencloud-eu/opencloud/pkg/clihelper"
	"github.com/opencloud-eu/opencloud/services/filepostprocessing/pkg/config"
	"github.com/spf13/cobra"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) []*cobra.Command {
	return []*cobra.Command{
		Server(cfg),
	}
}

// Execute is the entry point for the filepostprocessing command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cobra.Command{
		Use:   "filepostprocessing",
		Short: "starts filepostprocessing service",
	})
	app.AddCommand(GetCommands(cfg)...)
	// We need to set the args to os.Args[1:] so that cobra doesn't see the "filepostprocessing" prefix
	// when running in supervised mode, the runtime calls the binary with the service name as the first arg.
	app.SetArgs(os.Args[1:])
	return app.ExecuteContext(cfg.Context)
}
