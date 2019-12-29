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
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go app.Start(done)

	for {
		select {
		case <-sigs:
			fmt.Println("Done")
			close(done)
			os.Exit(0)
		}
	}
}
