package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mikhno-s/eva_tg_bot/app/scheme"
	"github.com/zelenin/go-tdlib/client"
)

// Start is start
func Start(done <-chan bool) {
	strAppID, _ := strconv.Atoi(os.Getenv("API_ID"))
	apiID := int32(strAppID)
	apiHash := os.Getenv("API_HASH")
	publicChannelUsername := os.Getenv("CHAN_NAME")
	dbConnString := os.Getenv("DB_CONN")

	tdlibClient, err := createTdlibClient(apiID, apiHash)

	checkError(err, "Creating telegram client")

	chat, err := tdlibClient.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: publicChannelUsername,
	})
	checkError(err, "Searching public channel")
	if chat == nil {
		checkError(fmt.Errorf("Cannot find %s", publicChannelUsername), "Searching public channel")
	}

	db := &carsDB{
		ConnString: dbConnString}

	err = db.Init()
	if err != nil {
		checkError(err, "DB initialization")
	}
	defer db.Close()

	// Creating tables
	db.Conn.MustExec(scheme.GetEvacuatedTableScheme())
	db.Conn.MustExec(scheme.GetTelegramLastChatMessagesTableScheme())

	cars, err := db.getAllCars()
	checkError(err, "Getting data from storage")

	// Fill storage if it's empty
	if len(cars) == 0 {
		messages := GetChanHistory(tdlibClient, chat.Id, chat.LastMessage.Id, 0)
		log.Println("Reading", len(messages), "messages")
		for _, m := range messages {
			e := scheme.MessageContentEntry{}
			// client.Message{} does not have Content.Text field, dirty hack
			mBytes, err := m.MarshalJSON()
			checkError(err, "Converting message to bytes")
			err = json.Unmarshal(mBytes, &e)
			checkError(err, "Unmarshaling message bytes message entry struct")

			// Fill
			for _, c := range getCarsInfoFromMessage(&e) {
				cars = append(cars, c)
			}
		}

		db.SetLastReadedMessageID(chat.Id, strconv.Itoa(int(messages[len(messages)-1].Id)))

		log.Println("Got", len(cars), "new cars")
		err = db.insertEvacuatedCars(cars)
		if err != nil {
			checkError(err, "Inserting data to storage")
		}

	}

	tickTenSeconds := time.NewTicker(time.Second * 10)

updatesLoop:
	for {
		select {
		case <-tickTenSeconds.C:
			// Check for a new message every 10 second
			lastMessageInChan, err := tdlibClient.GetChatMessageByDate(&client.GetChatMessageByDateRequest{
				ChatId: chat.Id,
				Date:   int32(time.Now().UTC().Unix()),
			})
			checkError(err, "Getting messages count")

			// Get last saved message
			lastSavedMessageID, err := db.getLastReadedMessageID(chat.Id)
			if err != nil {
				checkError(err, "Getting last saved messaged id")
			}

			if lastMessageInChan.Id > lastSavedMessageID {
				messages := GetChanHistory(tdlibClient, chat.Id, lastMessageInChan.Id, lastSavedMessageID)
				cars := make([]*scheme.Car, 0)
				log.Println("Reading", len(messages), "messages")
				for _, m := range messages {
					e := scheme.MessageContentEntry{}
					// client.Message{} does not have Content.Text field, dirty hack
					mBytes, err := m.MarshalJSON() //transform 'm' struct to bytes
					checkError(err, "Converting message to bytes")
					err = json.Unmarshal(mBytes, &e) // transforming bytes to struct with .Text field again
					checkError(err, "Unmarshaling message bytes message entry struct")

					// Fill
					for _, c := range getCarsInfoFromMessage(&e) {
						cars = append(cars, c)
					}
				}
				log.Println("Got", len(cars), "new cars")
				err = db.insertEvacuatedCars(cars)
				if err != nil {
					checkError(err, "Inserting data to storage")
				}

				db.SetLastReadedMessageID(chat.Id, strconv.Itoa(int(messages[len(messages)-1].Id)))
			}
		case <-done:
			break updatesLoop

		}
	}

}
