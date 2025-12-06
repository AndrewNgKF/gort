package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gort",
	Short: "Gort - A Rails-inspired web framework for Go",
	Long: `Gort is an opinionated web framework for Go that brings the developer
experience of Ruby on Rails to the Go ecosystem.

With Gort, you get:
  • Convention over configuration
  • Powerful CLI generators
  • RESTful routing DSL
  • Active Record-style ORM
  • Hot reload development
  • And much more!`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
}
