package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mikhno-s/eva_tg_bot/app"
	"github.com/mikhno-s/eva_tg_bot/app/config"
	"github.com/mikhno-s/eva_tg_bot/app/scheme"
	"github.com/zelenin/go-tdlib/client"
)

type Car struct {
	Date time.Time `json: date`
	Type string    `json: type`
	ID   string    `json: id`
	VIN  string    `json: vin`
}

type EvacuationCarsLogProcessor struct {
	Config *config.Config
}

func (p *EvacuationCarsLogProcessor) Start() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		cars := make([]*Car, 0)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println(event.Op & fsnotify.Write)
				// TODO FIX HERE. Find out why does not work here.
				if event.Op&fsnotify.Write == fsnotify.Write {
					f, err := os.OpenFile(p.Config.MessageStoragePath, os.O_RDONLY, 0755)
					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						e := scheme.MessageContentEntry{}
						m := client.Message{}
						err = json.Unmarshal(scanner.Bytes(), &m)
						if err != nil {
							return
						}
						mBytes, _ := m.MarshalJSON()
						json.Unmarshal(mBytes, &e)

						for _, c := range GetCarsInfo(e.Content.Text.Text) {
							c.Date = time.Unix(e.Date, 0)
							cars = append(cars, c)
							fmt.Printf("%+v\n", c)
						}
					}
					f.Close()

					fw, _ := os.OpenFile(p.Config.OutputPath, os.O_RDWR|os.O_CREATE, 0755)
					for _, c := range cars {
						str, _ := json.Marshal(c)
						fw.Write(str)
						fw.WriteString("\n")
					}

					fw.Close()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(p.Config.MessageStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {
	go app.Start()

	// Very VERY BAD CODE !!!
	// TODO CHANGE IT
	proc := &EvacuationCarsLogProcessor{}
	proc.Config = config.InitConfig()
	proc.Start()
}

func GetCarsInfo(row string) []*Car {
	cars := make([]*Car, 0)
	if strings.Count(row, "Евакуйовано") != strings.Count(row, "ДНЗ") && strings.Count(row, "Евакуйовано") != strings.Count(row, "VIN-код") {
		os.Exit(0)
	}

	b := strings.Split(row, "\n")
	for i, m := range b {
		if m != "" {
			if strings.Contains(m, "Евакуйовано:") {
				cars = append(cars, GetCarInfo(b[i:i+3]))
			}
		}

	}
	return cars
}

func GetCarInfo(data []string) *Car {
	car := new(Car)
	for _, str := range data {
		if strings.Contains(str, "Евакуйовано:") {
			t := strings.Split(str, "Евакуйовано:")
			car.Type = strings.TrimSpace(t[len(t)-1])
		}
		if strings.Contains(str, "ДНЗ:") {
			t := strings.Split(str, "ДНЗ:")
			car.ID = strings.TrimSpace(t[len(t)-1])
		}
		if strings.Contains(str, "VIN-код:") {
			t := strings.Split(str, "VIN-код:")
			car.VIN = strings.TrimSpace(t[len(t)-1])
		}
	}
	return car
}
