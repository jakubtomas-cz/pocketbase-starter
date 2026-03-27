package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func Register(pb core.App) {
	RegisterAPIs(pb)
	RegisterViews(pb)
}

func RegisterAPIs(pb core.App) {
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/hello/{name}", func(e *core.RequestEvent) error {
			name := e.Request.PathValue("name")
			return e.JSON(http.StatusOK, map[string]string{"message": "Hello, " + name})
		})

		return se.Next()
	})
}

func RegisterViews(pb core.App) {
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{$}", func(e *core.RequestEvent) error {
			registry := template.NewRegistry()
			html, err := registry.LoadFiles("views/index.html").Render(map[string]any{})
			if err != nil {
				return err
			}

			return e.HTML(http.StatusOK, html)
		})

		return se.Next()
	})
}
