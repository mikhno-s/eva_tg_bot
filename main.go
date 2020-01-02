package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikhno-s/eva_tg_bot/app"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := app.App{}
	go app.Start()

	for {
		select {
		case <-sigs:
			fmt.Println("Done")
			app.Stop()
			os.Exit(0)
		}
	}
}
