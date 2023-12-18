package main

import (
	"log"

	"github.com/ayush5588/shorturl/internal/router"
)

func main() {
	r := router.SetupRouter()
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("error in starting the router...", err)
	}
}
