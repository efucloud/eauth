package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/efucloud/eauth/cmd/server"
)

func main() {
	rand.NewSource(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	command := server.NewRunnerServerCommand()
	command.Flags().SortFlags = true
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

}
