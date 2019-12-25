package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mikhno-s/eva_tg_bot/app/storage"
	"github.com/zelenin/go-tdlib/client"
)

func Start() {

	strAppID, _ := strconv.Atoi(os.Getenv("API_ID"))
	apiID := int32(strAppID)
	apiHash := os.Getenv("API_HASH")
	publicChannelUsername := os.Getenv("CHAN_NAME")
	storageFilePath := os.Getenv("STORAGE_FILE")

	tdlibClient, err := createTdlibClient(apiID, apiHash)

	checkErrorFatal(err, "Creating telegram client")

	chat, err := tdlibClient.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: publicChannelUsername,
	})
	checkErrorFatal(err, "Searching public channel")
	if chat == nil {
		checkErrorFatal(fmt.Errorf("Cannot find %s", publicChannelUsername), "Searching public channel")
	}

	storage, err := storage.InitStorage(&storage.DBFile{
		Path: storageFilePath,
	})
	checkErrorFatal(err, "File initialization")

	defer storage.Close()

	messages, err := storage.ReadMessages()

	var lastSavedMessageID int64

	// Fill storage if it's empty
	if len(messages) == 0 {
		messages = GetChanHistory(tdlibClient, chat.Id, chat.LastMessage.Id, 0)
		err = storage.WriteMessages(messages)
		checkErrorFatal(err, "Writing messages to file")
	}

	lastSavedMessageID = messages[len(messages)-1].Id

	tickTenSeconds := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-tickTenSeconds.C:
			// Check for a new message every 10 second
			lastMessageInChan, err := tdlibClient.GetChatMessageByDate(&client.GetChatMessageByDateRequest{
				ChatId: chat.Id,
				Date:   int32(time.Now().UTC().Unix()),
			})
			checkErrorFatal(err, "Getting messages count")
			if lastMessageInChan.Id != lastSavedMessageID {
				newMessages := GetChanHistory(tdlibClient, chat.Id, lastMessageInChan.Id, lastSavedMessageID)
				for _, m := range newMessages {
					messages = append(messages, m)
				}
				storage.WriteMessages(newMessages)
				lastSavedMessageID = messages[len(messages)-1].Id
			}

		}
	}
}
