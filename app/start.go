package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/mikhno-s/eva_tg_bot/app/scheme"
	"github.com/zelenin/go-tdlib/client"
)

func Start() {

	strAppID, _ := strconv.Atoi(os.Getenv("API_ID"))
	apiID := int32(strAppID)
	apiHash := os.Getenv("API_HASH")
	publicChannelUsername := os.Getenv("CHAN_NAME")
	storageFile := os.Getenv("STORAGE_FILE")

	tdlibClient, err := createTdlibClient(apiID, apiHash)

	checkErrorFatal(err, "Creating telegram client")

	chat, err := tdlibClient.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: publicChannelUsername,
	})

	checkErrorFatal(err, "Searching public channel")
	if chat == nil {
		checkErrorFatal(fmt.Errorf("Cannot find %s", publicChannelUsername), "Searching public channel")
	}

	// Read storage file
	f, err := os.OpenFile(storageFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	checkErrorFatal(err, "Opening file")
	defer f.Close()
	scanner := bufio.NewScanner(f)

	fStat, err := f.Stat()
	checkErrorFatal(err, "Reading file size")

	messages := make([]*client.Message, 0)

	if fStat.Size() > 0 {
		for scanner.Scan() {
			m := client.Message{}
			err := json.Unmarshal(scanner.Bytes(), &m)
			checkErrorFatal(err, "Reading file")
			messages = append(messages, &m)
		}
	}

	var lastSavedMessageID int64

	if len(messages) != 0 {
		lastSavedMessageID = messages[len(messages)-1].Id
	}

	// If chat has newer messages or we don't have saved messages at all
	if chat.LastMessage.Id != lastSavedMessageID {
		fmt.Println(chat.LastMessage.Id, lastSavedMessageID)
		for _, m := range GetChanHistory(tdlibClient, chat.Id, chat.LastMessage.Id, lastSavedMessageID) {

			// Append to saved state
			marhMessage, err := m.MarshalJSON()
			checkErrorFatal(err, "Printing messages")
			f.WriteString(string(marhMessage) + "\n")

			// Append to memory state
			messages = append(messages, m)
		}
	}

	for _, m := range messages {
		e := scheme.MessageContentEntry{}
		mBytes, err := m.MarshalJSON()
		checkErrorFatal(err, "Json marshalling err")
		err = json.Unmarshal(mBytes, &e)
		checkErrorFatal(err, "Json unmarshaling")
		fmt.Println(e.Content.Text.Text)
	}

}
