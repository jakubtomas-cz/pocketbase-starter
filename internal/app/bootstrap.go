package app

import (
	"net/http"
	"os"

	"pocketbase-starter/internal/hooks"
	"pocketbase-starter/internal/routes"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/ghupdate"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/hook"
)

// Register wires flags, plugins, hooks, and routes onto the PocketBase app.
// Mirrors the setup of examples/base/main.go in the official repo.
func Register(pb *pocketbase.PocketBase) {
	// Default PB setup
	var hooksDir string
	pb.RootCmd.PersistentFlags().StringVar(
		&hooksDir, "hooksDir", "", "the directory with the JS app hooks",
	)

	var hooksWatch bool
	pb.RootCmd.PersistentFlags().BoolVar(
		&hooksWatch, "hooksWatch", true, "auto restart the app on pb_hooks file change; it has no effect on Windows",
	)

	var hooksPool int
	pb.RootCmd.PersistentFlags().IntVar(
		&hooksPool, "hooksPool", 15, "the total prewarm goja.Runtime instances for the JS app hooks execution",
	)

	var migrationsDir string
	pb.RootCmd.PersistentFlags().StringVar(
		&migrationsDir, "migrationsDir", "", "path to user-defined migrations",
	)

	var automigrate bool
	pb.RootCmd.PersistentFlags().BoolVar(
		&automigrate, "automigrate", true, "enable/disable auto migrations",
	)

	var publicDir string
	pb.RootCmd.PersistentFlags().StringVar(
		&publicDir, "publicDir", "./pb_public", "the directory to serve static files",
	)

	var indexFallback bool
	pb.RootCmd.PersistentFlags().BoolVar(
		&indexFallback, "indexFallback", true, "fallback the request to index.html on missing static path, e.g. when pretty urls are used with SPA",
	)

	pb.RootCmd.ParseFlags(os.Args[1:])

	// ---------------------------------------------------------------
	// Plugins and hooks:
	// ---------------------------------------------------------------

	// load jsvm (pb_hooks and pb_migrations)
	jsvm.MustRegister(pb, jsvm.Config{
		MigrationsDir: migrationsDir,
		HooksDir:      hooksDir,
		HooksWatch:    hooksWatch,
		HooksPoolSize: hooksPool,
	})

	migratecmd.MustRegister(pb, pb.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangGo,
		Automigrate:  automigrate,
		Dir:          migrationsDir,
	})

	ghupdate.MustRegister(pb, pb.RootCmd, ghupdate.Config{})

	// static route to serves files from the provided public dir
	// (if publicDir exists and the route path is not already defined)
	pb.App.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Func: func(e *core.ServeEvent) error {
			if !e.Router.HasRoute(http.MethodGet, "/{path...}") {
				e.Router.GET("/{path...}", apis.Static(os.DirFS(publicDir), indexFallback))
			}

			return e.Next()
		},
		Priority: 999, // execute as latest as possible to allow users to provide their own route
	})

	hooks.Register(pb)
	routes.Register(pb)
}
