package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Display all registered routes",
	Long:  `Display all routes registered in the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we're in a Gort app directory
		if _, err := os.Stat("config/routes.go"); os.IsNotExist(err) {
			fmt.Println("Error: Not in a Gort application directory")
			fmt.Println("Make sure you're in the root directory of your Gort app.")
			os.Exit(1)
		}

		fmt.Println("\n📋 Application Routes:")
		fmt.Println("=====================\n")
		fmt.Println("To see routes, your application must be running.")
		fmt.Println("Add this to your config/routes.go after setting up routes:\n")
		fmt.Println("  r.PrintRoutes()")
		fmt.Println("\nOr run your application and check the router setup.")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(routesCmd)
}
