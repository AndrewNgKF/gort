package commands

import (
	"fmt"
	"os"

	"github.com/AndrewNgKF/gort/internal/generator"
	"github.com/spf13/cobra"
)

var apiMode bool

var newCmd = &cobra.Command{
	Use:   "new [app_name]",
	Short: "Create a new Gort application",
	Long: `Create a new Gort application with the complete directory structure,
configuration files, and boilerplate code.

Example:
  gort new myblog
  gort new my_awesome_api --api`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]

		// Check if directory already exists
		if _, err := os.Stat(appName); !os.IsNotExist(err) {
			fmt.Printf("Error: Directory '%s' already exists\n", appName)
			os.Exit(1)
		}

		if apiMode {
			fmt.Printf("Creating new Gort API application: %s\n", appName)
		} else {
			fmt.Printf("Creating new Gort application: %s\n", appName)
		}

		// Create the application
		if err := generator.NewApp(appName, apiMode); err != nil {
			fmt.Printf("Error creating application: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✨ Successfully created %s!\n\n", appName)
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", appName)
		fmt.Println("  gort server")
		fmt.Println("\nHappy coding! 🚀")
	},
}

func init() {
	newCmd.Flags().BoolVar(&apiMode, "api", false, "Create an API-only application (no views or assets)")
	rootCmd.AddCommand(newCmd)
}
