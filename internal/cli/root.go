package cli

import (
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var banner = `
  __  __  ____    ____   _____  _   _
 |  \/  |/ ___|  / ___| | ____|| \ | |
 | |\/| |\___ \ | |  _  |  _|  |  \| |
 | |  | | ___) || |_| | | |___ | |\  |
 |_|  |_||____/  \____| |_____||_| \_|

 The Ultimate Microservice Scaffolder
`

var rootCmd = &cobra.Command{
	Use:   "go-micro-gen",
	Short: "Production-ready Go microservice generator",
	Long:  color.CyanString(banner) + "\n  Generate fully-wired, production-ready Go microservices in seconds.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags.
func Execute() error {
	return rootCmd.Execute()
}

var Version = "v1.3.7"

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return Version
	}
	return info.Main.Version
}

func init() {
	rootCmd.Version = getVersion()
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}
