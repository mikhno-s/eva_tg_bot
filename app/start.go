package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mikhno-s/eva_tg_bot/app/scheme"
	"github.com/zelenin/go-tdlib/client"
)

type App struct {
	config struct {
		APIID    int32
		APIHash  string
		ChanName string
		DbURL    string
	}
	stopChan    chan bool
	db          carsDB
	tdlibClient *client.Client
	chat        *client.Chat
	wg          sync.WaitGroup
}

func (app *App) Stop() {
	close(app.stopChan)
	app.wg.Wait()
	app.db.Close()
}

func (app *App) Init() {
	app.wg.Add(1)
	var err error

	app.stopChan = make(chan bool)

	strAppID, _ := strconv.Atoi(os.Getenv("API_ID"))
	app.config.APIID = int32(strAppID)
	app.config.APIHash = os.Getenv("API_HASH")
	app.config.ChanName = os.Getenv("CHAN_NAME")
	app.config.DbURL = os.Getenv("DB_CONN")

	app.tdlibClient, err = createTdlibClient(app.config.APIID, app.config.APIHash)
	checkError(err, "Creating telegram client")

	app.chat, err = app.tdlibClient.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: app.config.ChanName,
	})
	checkError(err, "Searching public channel")
	if app.chat == nil {
		checkError(fmt.Errorf("Cannot find %s", app.config.ChanName), "Searching public channel")
	}
	app.wg.Done()

	app.db = carsDB{
		ConnString: app.config.DbURL}

	err = app.db.Init()
	if err != nil {
		checkError(err, "DB initialization")
	}

}

// Start is start
func (app *App) Start() {
	app.wg = sync.WaitGroup{}
	app.wg.Add(1)
	app.Init()
	log.Println("App started")

	// Creating tables
	app.db.Conn.MustExec(scheme.GetEvacuatedTableScheme())
	app.db.Conn.MustExec(scheme.GetTelegramLastChatMessagesTableScheme())

	cars, err := app.db.getAllCars()
	checkError(err, "Getting data from storage")

	// Fill storage if it's empty
	if len(cars) == 0 {
		messages := GetChanHistory(app.tdlibClient, app.chat.Id, app.chat.LastMessage.Id, 0)
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

		app.db.SetLastReadedMessageID(app.chat.Id, strconv.Itoa(int(messages[len(messages)-1].Id)))

		log.Println("Got", len(cars), "new cars")
		err = app.db.insertEvacuatedCars(cars)
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
			lastMessageInChan, err := app.tdlibClient.GetChatMessageByDate(&client.GetChatMessageByDateRequest{
				ChatId: app.chat.Id,
				Date:   int32(time.Now().UTC().Unix()),
			})
			checkError(err, "Getting messages count")

			// Get last saved message
			lastSavedMessageID, err := app.db.getLastReadedMessageID(app.chat.Id)
			if err != nil {
				checkError(err, "Getting last saved messaged id")
			}

			if lastMessageInChan.Id > lastSavedMessageID {
				messages := GetChanHistory(app.tdlibClient, app.chat.Id, lastMessageInChan.Id, lastSavedMessageID)
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
				err = app.db.insertEvacuatedCars(cars)
				if err != nil {
					checkError(err, "Inserting data to storage")
				}

				app.db.SetLastReadedMessageID(app.chat.Id, strconv.Itoa(int(messages[len(messages)-1].Id)))
			}
		case <-app.stopChan:
			app.wg.Done()
			break updatesLoop
		}
	}

}
