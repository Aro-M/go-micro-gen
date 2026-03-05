package cli

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/Aro-M/go-micro-gen/internal/ui"
	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch the interactive web dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		port := "8080"
		url := fmt.Sprintf("http://localhost:%s", port)

		fmt.Printf("🚀 Starting Interactive Web Server on %s\n", url)

		// Wait a small amount of time to ensure the HTTP server yields properly, then open browser
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser(url)
		}()

		// Start the embedded HTTP listener
		return ui.StartServer(port)
	},
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Printf("ℹ️  Automatic browser launch failed. Please open %s in your browser\n", url)
	}
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
