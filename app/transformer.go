package app

import (
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mikhno-s/eva_tg_bot/app/scheme"
)

// getCarsInfoFromMessage open and transform a telegram message text into slice of car structs
func getCarsInfoFromMessage(msg *scheme.MessageContentEntry) (cars []*scheme.Car) {

	// Checking that messages text has equal count of these fields
	text := msg.Content.Text.Text
	if strings.Count(text, "Евакуйовано") != strings.Count(text, "ДНЗ") && strings.Count(text, "Евакуйовано") != strings.Count(text, "VIN-код") {
		os.Exit(0)
	}

	rows := strings.Split(text, "\n")
	for i, r := range rows {
		if r != "" {
			// Example of evacuation info:
			// 	01)
			// 	Евакуйовано: HYUNDAI i
			// 	ДНЗ: AAAAAAAA i + 1
			// 	VIN-код: *01234 i + 2
			if strings.Contains(r, "Евакуйовано:") {
				car := getCarInfo(rows[i : i+3])
				car.Date = time.Unix(msg.Date, 0)
				// Uses msg id, car license plate and string position in message, then it will be hashed by md5
				uniqCarID := strings.Join([]string{strconv.Itoa(int(msg.ID)), car.LicensePlate, strconv.Itoa(i)}, "-")
				hashedID := md5.Sum([]byte(uniqCarID))
				car.ID = fmt.Sprintf("%x", hashedID)
				cars = append(cars, car)
			}
		}

	}
	return cars
}

// getCarInfo fills car struct
// uses arguments a slice of string like ["Eвакуйовано: HYUNDAI i30", "ДНЗ: AAAAAAAA", "VIN-код: *01234"]
func getCarInfo(data []string) (car *scheme.Car) {
	car = &scheme.Car{}
	for _, str := range data {
		if strings.Contains(str, "Евакуйовано:") {
			t := strings.Split(str, "Евакуйовано:")
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
	return
}
