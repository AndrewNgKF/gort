package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	port string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Gort application server",
	Long: `Starts the Gort development server.

By default, the server runs on http://localhost:3000

Examples:
  gort server              # Start server on default port 3000
  gort server -p 8080      # Start server on custom port
  gort server --port 4000  # Alternative syntax`,
	Aliases: []string{"s", "serve"},
	Run: func(cmd *cobra.Command, args []string) {
		// Check if main.go exists
		if _, err := os.Stat("main.go"); os.IsNotExist(err) {
			fmt.Println("❌ Error: main.go not found in current directory")
			fmt.Println("Are you in a Gort application directory?")
			os.Exit(1)
		}

		// Set PORT environment variable if specified
		if port != "" {
			os.Setenv("PORT", port)
		} else {
			// Default port
			if os.Getenv("PORT") == "" {
				os.Setenv("PORT", "3000")
			}
		}

		fmt.Println("🚀 Starting Gort server...")
		fmt.Printf("=> Listening on http://localhost:%s\n", os.Getenv("PORT"))
		fmt.Println("=> Press Ctrl+C to stop")
		fmt.Println()

		// Get the current directory name for the module
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Run go run main.go
		runCmd := exec.Command("go", "run", "main.go")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin
		runCmd.Dir = currentDir

		if err := runCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Printf("❌ Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server on (default: 3000)")
}
