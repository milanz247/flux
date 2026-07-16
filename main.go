// Flux — a Laravel + Inertia inspired full-stack framework for Go.
//
// This is the application entry point: load config, connect the database,
// run pending migrations, build the app container, register routes, serve.
package main

import (
	"fmt"
	"os"

	"flux/app/routes"
	"flux/config"
	"flux/database"
	_ "flux/database/migrations" // register migrations
	"flux/framework"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fail("config", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		fail("database", err)
	}

	if err := database.Migrate(db); err != nil {
		fail("migrate", err)
	}

	app, err := framework.New(cfg, framework.WithDB(db))
	if err != nil {
		fail("bootstrap", err)
	}

	routes.Register(app)

	if err := app.Serve(); err != nil {
		fail("serve", err)
	}
}

func fail(stage string, err error) {
	fmt.Fprintf(os.Stderr, "flux: %s: %v\n", stage, err)
	os.Exit(1)
}
