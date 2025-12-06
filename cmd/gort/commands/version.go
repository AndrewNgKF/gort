package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/AndrewNgKF/gort/pkg/gort"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the Gort framework version",
	Long:  "Display the current version of the Gort framework installed on your system.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Gort v%s\n", gort.Version)
		fmt.Println("A Rails-inspired web framework for Go")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
