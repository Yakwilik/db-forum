package main

import (
	"github.com/db-forum.git/internal/app/forum"
	"log"
)

func main() {
	app, err := forum.NewForum()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.StartApp(); err != nil {
		log.Fatal(err)
	}

}
