package app

import (
	"fmt"
	"sort"

	"github.com/zelenin/go-tdlib/client"
)

// GetChanHistory returns slice of messages
func GetChanHistory(tdlibClient *client.Client, chatID int64, fromMessageID int64, toMessageID int64) (messages []*client.Message) {
	var totalMessages int

	messagesSet := make(map[int]*client.Message)
	totalLimit := 99999999999

	// Read first message (newest) separetely, because messageReading does not return exactly message - fromMessageId
	if fromMessageID != 0 {
		lastMessage, err := tdlibClient.GetMessage(&client.GetMessageRequest{ChatId: chatID, MessageId: fromMessageID})
		checkErrorFatal(err, "Getting chan history")
		messagesSet[int(lastMessage.Id)] = lastMessage
	}
messageReading:
	for {
		fmt.Println("Retriving messages from ", fromMessageID, "..")
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
				break messageReading
			}
			// Read to needed MessageID
			if toMessageID == m.Id {
				break messageReading
			}
			totalMessages++

			// Read next set of messages
			fromMessageID = m.Id
			messagesSet[int(m.Id)] = m
		}
	}

	messagesIDsSorted := make([]int, 0, len(messagesSet))

	for k := range messagesSet {
		messagesIDsSorted = append(messagesIDsSorted, k)
	}
	sort.Ints(messagesIDsSorted)
	for _, i := range messagesIDsSorted {
		messages = append(messages, messagesSet[i])
	}

	return
}
