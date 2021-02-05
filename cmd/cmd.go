package cmd

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/viknesh-nm/go-listener/utils"
	"github.com/viknesh-nm/go-listener/watcher"
)

var (
	path, name, commands string
	buildOnly            bool
	rootCmd              = &cobra.Command{
		Use: "go-listener",
		Run: func(cmd *cobra.Command, args []string) {
			go utils.GracefulShutdown()
			if strings.TrimSpace(path) == "" {
				path = utils.GetDefaultPath()
			}

			if name == "" {
				name = filepath.Base(path)
			}

			w := watcher.New(
				path,
				name,
				buildOnly,
				utils.ValidateCommands(commands),
			)

			if err := w.Watch(); err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "directory to be watch")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "name of the project")
	rootCmd.PersistentFlags().StringVarP(&commands, "command", "c", "", "custom commands that needs to be run")
	rootCmd.PersistentFlags().BoolVarP(&buildOnly, "build", "b", false, "build only mode that generates the binary file")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
