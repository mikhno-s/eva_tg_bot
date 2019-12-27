package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikhno-s/eva_tg_bot/app"
	"github.com/mikhno-s/eva_tg_bot/transformer"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go app.Start()
	proc := &transformer.EvacuationCarsLogProcessor{}
	go proc.Start()

	for {
		select {
		case <-sigs:
			fmt.Println("Done")
			os.Exit(0)
		}
	}
}
