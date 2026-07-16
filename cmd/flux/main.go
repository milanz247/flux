// Command flux is the Flux framework CLI — the equivalent of Laravel's
// artisan. Run `go run ./cmd/flux` (or build it once: `go build -o flux.exe
// ./cmd/flux`) from the project root.
//
//	flux new <name>            scaffold a new project
//	flux serve                 run the Go server + Vite dev server
//	flux build                 build frontend + backend for production
//	flux migrate               run pending database migrations
//	flux seed                  run database seeders
//	flux make:controller User  generate app/controllers/user_controller.go
//	flux make:model User       generate app/models/user.go
//	flux make:service User     generate app/services/user_service.go
//	flux make:dto User         generate app/dto/user.go
//	flux make:middleware Auth  generate app/middleware/auth.go
//	flux make:migration create_users
//	flux make:page Users/Index generate resources/js/Pages/Users/Index.vue
package main

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"flux/config"
	"flux/database"
	_ "flux/database/migrations" // register migrations
	"flux/database/seeders"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "new":
		err = cmdNew(args)
	case "serve":
		err = cmdServe()
	case "build":
		err = cmdBuild()
	case "migrate":
		err = withDB(database.Migrate)
	case "seed":
		err = withDB(seeders.Run)
	case "make:controller":
		err = makeFromStub(args, "controller")
	case "make:model":
		err = makeFromStub(args, "model")
	case "make:service":
		err = makeFromStub(args, "service")
	case "make:dto":
		err = makeFromStub(args, "dto")
	case "make:middleware":
		err = makeFromStub(args, "middleware")
	case "make:migration":
		err = makeMigration(args)
	case "make:page":
		err = makePage(args)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "flux: unknown command %q\n\n", command)
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "flux: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`Flux — Laravel-style full-stack framework for Go

Usage:
  flux new <name>              Scaffold a new Flux project
  flux serve                   Run Go server + Vite dev server together
  flux build                   Production build (Vite bundle + Go binary)
  flux migrate                 Run pending database migrations
  flux seed                    Run database seeders

Generators:
  flux make:controller <Name>  New resource controller
  flux make:model <Name>       New GORM model
  flux make:service <Name>     New service
  flux make:dto <Name>         New DTO set
  flux make:middleware <Name>  New middleware
  flux make:migration <name>   New migration (e.g. create_posts)
  flux make:page <Path>        New Vue page (e.g. Users/Index)`)
}

// withDB loads config, makes sure the database exists, opens it and runs fn.
func withDB(fn func(db *gorm.DB) error) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := database.EnsureDatabase(cfg); err != nil {
		return err
	}
	db, err := database.Connect(cfg)
	if err != nil {
		return err
	}
	return fn(db)
}
