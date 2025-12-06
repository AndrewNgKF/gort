package commands

import (
	"fmt"
	"os"

	"github.com/AndrewNgKF/gort/internal/generator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate code (aliases: g, gen)",
	Long:    `Generate models, controllers, migrations, and more.`,
	Aliases: []string{"g", "gen"},
}

var generateControllerCmd = &cobra.Command{
	Use:   "controller [name] [actions...]",
	Short: "Generate a new controller",
	Long: `Generate a new controller with specified actions.

Example:
  gort generate controller Posts index show create update destroy
  gort g controller Users index show
  gort gen controller Home index about contact`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		actions := []string{}
		if len(args) > 1 {
			actions = args[1:]
		}

		fmt.Printf("Creating controller: %s\n", name)

		if err := generator.GenerateController(name, actions); err != nil {
			fmt.Printf("Error generating controller: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✨ Controller created successfully!\n\n")
		fmt.Println("Next steps:")
		fmt.Printf("  1. Add routes in config/routes.go\n")
		fmt.Printf("  2. Create views in app/views/%s/\n", generator.ToSnakeCase(name))

		if len(actions) > 0 {
			fmt.Printf("\nGenerated actions:\n")
			for _, action := range actions {
				fmt.Printf("  - %s\n", action)
			}
		}
	},
}

var generateModelCmd = &cobra.Command{
	Use:   "model [name] [field:type...]",
	Short: "Generate a new model",
	Long: `Generate a new model with specified fields and a migration.

Example:
  gort generate model Post title:string body:text published:boolean
  gort g model User name:string email:string:unique age:integer
  gort gen model Product name:string price:decimal stock:integer`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fields := []string{}
		if len(args) > 1 {
			fields = args[1:]
		}

		fmt.Printf("Creating model: %s\n", name)

		if err := generator.GenerateModel(name, fields); err != nil {
			fmt.Printf("Error generating model: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✨ Model created successfully!\n\n")
		fmt.Println("Next steps:")
		fmt.Println("  1. Run migrations: gort db:migrate")
		fmt.Println("  2. Use the model in your controllers")

		if len(fields) > 0 {
			fmt.Printf("\nGenerated fields:\n")
			for _, field := range fields {
				fmt.Printf("  - %s\n", field)
			}
		}
	},
}

var generateScaffoldCmd = &cobra.Command{
	Use:   "scaffold [name] [field:type...]",
	Short: "Generate a complete CRUD resource",
	Long: `Generate a model, migration, controller, views, and routes for a complete CRUD resource.

Example:
  gort generate scaffold Post title:string body:text published:boolean
  gort g scaffold Product name:string price:decimal stock:integer
  gort gen scaffold Article title:string content:text author:string`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fields := []string{}
		if len(args) > 1 {
			fields = args[1:]
		}

		if err := generator.GenerateScaffold(name, fields); err != nil {
			fmt.Printf("Error generating scaffold: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateControllerCmd)
	generateCmd.AddCommand(generateModelCmd)
	generateCmd.AddCommand(generateScaffoldCmd)
}
