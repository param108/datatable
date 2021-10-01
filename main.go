package main

import (
	"log"

	"os"

	"github.com/param108/datatable/cmd"
	mylog "github.com/param108/datatable/log"
)

func main() {
	defer mylog.Close()

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
