package main

import (
	"log"

	"github.com/project-ippl-dev/tanding-api/internal/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal("Running Failed : ", err.Error())
	}
}
