package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

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

	// Create messages slice that can be used as in mem storage

	// Read storage file
	f, err := os.OpenFile(storageFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	checkErrorFatal(err, "Opening file")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messages := make([]*client.Message, 0)
	messagesStack := make([]*client.Message, 0)

	for scanner.Scan() {
		m := client.Message{}
		err := json.Unmarshal(scanner.Bytes(), &m)
		checkErrorFatal(err, "Reading file")
		messages = append(messages, &m)
	}

	var lastSavedMessageID int64
	if len(messages) != 0 {
		lastSavedMessageID = messages[len(messages)-1].Id
	}

	if lastSavedMessageID != chat.LastMessage.Id {
		// Should be append, because we want to save previous data. Append should perform in revert(?) orders (Fifo)
		messagesStack = GetChanHistory(tdlibClient, chat.Id, lastSavedMessageID, 0)
	}

	for _, m := range messages {
		fmt.Printf("%v ", m.Id)
	}
	fmt.Println()

	for i := range messagesStack {
		messages = append(messages, messagesStack[len(messagesStack)-i-1])
	}

	for _, m := range messages {
		fmt.Printf("%v ", m.Id)
	}
	fmt.Println()
	os.Exit(0)
	// Create flushing data to file
	// File must save data in log format (last messages - last in file)
	for _, m := range messages {
		marhMessage, err := m.MarshalJSON()
		checkErrorFatal(err, "Printing messages")
		fmt.Println(string(marhMessage))
		f.WriteString(string(marhMessage) + "\n")
	}
	// 696254464

}

// GetChanHistory returns slise of messages
// TODO Add offset message read and limit as param
func GetChanHistory(tdlibClient *client.Client, chatID int64, fromMessageID int64, toMessageID int64) (messages []*client.Message) {
	var totalMessages int

	totalLimit := 10

	for {
		chanHistory, err := tdlibClient.GetChatHistory(&client.GetChatHistoryRequest{
			ChatId:        chatID,
			Limit:         100,
			OnlyLocal:     false,
			FromMessageId: fromMessageID,
		})
		checkErrorFatal(err, "Getting chan history")
		if chanHistory.TotalCount == 0 {
			break
		}
		for _, m := range chanHistory.Messages {
			if totalLimit > 0 && totalMessages >= totalLimit {
				return
			}
			// Read til needed message
			if toMessageID == m.Id {
				return
			}
			totalMessages++
			messages = append(messages, m)
		}
		fromMessageID = messages[totalMessages-1].Id
		if totalLimit > 0 && totalMessages >= totalLimit {
			break
		}
	}

	return
}
