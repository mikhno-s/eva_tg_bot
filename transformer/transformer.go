package transformer

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/mikhno-s/eva_tg_bot/app/config"
	"github.com/mikhno-s/eva_tg_bot/app/scheme"
	"github.com/zelenin/go-tdlib/client"
)

type Car struct {
	Date         time.Time `json: date`
	Model        string    `json: model`
	LicensePlate string    `json: license_plate`
	VIN          string    `json: vin`
}

type EvacuationCarsLogProcessor struct {
	Config *config.Config
}

func (p *EvacuationCarsLogProcessor) Start() {
	p.Config = config.InitConfig()
	waiter := time.NewTicker(time.Second * 10)

	var lastModify time.Time

	cars := make([]*Car, 0)

	f, err := os.OpenFile(p.Config.MessageStoragePath, os.O_RDONLY, 0755)
	fOutput, _ := os.OpenFile(p.Config.OutputPath, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	defer fOutput.Close()
	defer waiter.Stop()

	connStr := "postgres://admin:admin@localhost/cars?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		select {
		case <-waiter.C:
			f.Sync()
			fStat, _ := f.Stat()
			if fStat.ModTime().After(lastModify) {

				// Writing to DB
				_, err := db.Query("SELECT 'evacuated_cars'::regclass")
				if err != nil {
					log.Fatalln(err.Error())
				}

				txn, err := db.Begin()
				stmt, err := txn.Prepare(pq.CopyIn("evacuated_cars", "date", "model", "license_plate", "vin"))

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

					newCars := GetCarsInfo(e.Content.Text.Text)

					for _, c := range newCars {
						c.Date = time.Unix(e.Date, 0)

						// Preparing statement
						_, err = stmt.Exec(c.Date, c.Model, c.LicensePlate, c.VIN)
						if err != nil {
							log.Fatal(err)
						}

						cars = append(cars, c)
						// fmt.Printf("%+v\n", c)
					}
				}

				fmt.Println(len(cars))

				err = stmt.Close()
				if err != nil {
					log.Fatal(err)
				}

				err = txn.Commit()
				if err != nil {
					log.Fatal(err)
				}

				fOutput.Truncate(0)
				fOutput.Seek(0, 0)

				fOutput.WriteString("[")
				for i, c := range cars {
					str, _ := json.Marshal(c)

					if err != nil {
						log.Fatal(err)
					}

					fOutput.Write(str)
					if i == len(cars)-1 {
						fOutput.WriteString("]")
					} else {
						fOutput.WriteString(",\n")
					}
				}

				lastModify = fStat.ModTime()
			}
		}
	}

}

func GetCarsInfo(row string) []*Car {
	cars := make([]*Car, 0)
	if strings.Count(row, "Евакуйовано") != strings.Count(row, "ДНЗ") && strings.Count(row, "Евакуйовано") != strings.Count(row, "VIN-код") {
		os.Exit(0)
	}

	b := strings.Split(row, "\n")
	for i, m := range b {
		if m != "" {
			// 09)
			// Евакуйовано: HYUNDAI
			// ДНЗ: AAAAAAAA
			// VIN-код: *01412
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
			// TODO probably, should use strings.Trim() here
			car.Model = strings.TrimSpace(t[len(t)-1])
		}
		if strings.Contains(str, "ДНЗ:") {
			t := strings.Split(str, "ДНЗ:")
			car.LicensePlate = strings.TrimSpace(t[len(t)-1])
		}
		if strings.Contains(str, "VIN-код:") {
			t := strings.Split(str, "VIN-код:")
			car.VIN = strings.TrimSpace(t[len(t)-1])
		}
	}
	return car
}
