package hooks

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
)

func Register(pb core.App) {
	pb.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		log.Printf("[hook] record create request: collection=%s", e.Collection.Name)
		return e.Next()
	})
}
