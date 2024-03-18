package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "app"}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "about",
		Short: "About App",
		Run: func(cmd *cobra.Command, args []string) {
			PrintTextFunc()
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Database Seeding 'seed'",
		Run: func(cmd *cobra.Command, args []string) {
			SeedFunc()
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:migration",
		Short: "Make Migration File 'make:migration migration_file_name'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileName := args[0]
			MigrateMake(fileName)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate:up",
		Short: "Migrate Up 'migrate:up step(int)'",
		Run: func(cmd *cobra.Command, args []string) {
			MigrateUpFunc()
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate:down",
		Short: "Migrate Down 'migrate:down step(int)'",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				step := "1"
				MigrateDownFunc(step)
			} else {
				step := args[0]
				MigrateDownFunc(step)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate:fresh",
		Short: "Migrate Fresh 'migrate:fresh'",
		Run: func(cmd *cobra.Command, args []string) {
			MigrateFreshFunc()
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:controller",
		Short: "Create a new controller 'make:controller path/controller_name'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			MakeController(args[0])
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:middleware",
		Short: "Create a new middleware 'make:middleware path/middleware_name'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			MakeMiddleware(args[0])
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:model",
		Short: "Create a new model 'make:model path/model_name'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			MakeModel(args[0])
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:repository",
		Short: "Create a new repository 'make:repository path/repository_name'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			MakeRepository(args[0])
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "swag",
		Short: "Generate Swagger 'swag'",
		Run: func(cmd *cobra.Command, args []string) {
			GenerateSwag()
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Serve App 'serve'",
		Run: func(cmd *cobra.Command, args []string) {
			ServeFunc()
		},
	})
}
