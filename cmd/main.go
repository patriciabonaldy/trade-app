package main

import (
	_ "embed"
	"log"

	"github.com/patriciabonaldy/zero/cmd/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
