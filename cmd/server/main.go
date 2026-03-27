package main

import (
	"log"

	"pocketbase-starter/internal/app"

	"github.com/pocketbase/pocketbase"

	_ "pocketbase-starter/migrations"
)

func main() {
	pb := pocketbase.New()
	app.Register(pb)

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
