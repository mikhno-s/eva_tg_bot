package storage

import (
	"github.com/zelenin/go-tdlib/client"
)

type Row struct {
	Row []byte
}

type MessageRow interface {
	GetRows() []*client.Message
	GetRow(n int) *client.Message
}

type DBFile struct {
	Path    string
	Content []Row
}

// read

// write

// TelMessage is TDLib message struct
type TelMessage struct {
	Type                    string `json:"@type"`
	Extra                   string `json:"@extra"`
	ID                      int    `json:"id"`
	SenderUserID            int    `json:"sender_user_id"`
	ChatID                  int64  `json:"chat_id"`
	IsOutgoing              bool   `json:"is_outgoing"`
	CanBeEdited             bool   `json:"can_be_edited"`
	CanBeForwarded          bool   `json:"can_be_forwarded"`
	CanBeDeletedOnlyForSelf bool   `json:"can_be_deleted_only_for_self"`
	CanBeDeletedForAllUsers bool   `json:"can_be_deleted_for_all_users"`
	IsChannelPost           bool   `json:"is_channel_post"`
	ContainsUnreadMention   bool   `json:"contains_unread_mention"`
	Date                    int    `json:"date"`
	EditDate                int    `json:"edit_date"`
	ReplyToMessageID        int    `json:"reply_to_message_id"`
	TTL                     int    `json:"ttl"`
	TTLExpiresIn            int    `json:"ttl_expires_in"`
	ViaBotUserID            int    `json:"via_bot_user_id"`
	AuthorSignature         string `json:"author_signature"`
	Views                   int    `json:"views"`
	MediaAlbumID            string `json:"media_album_id"`
	Content                 struct {
		Type  string `json:"@type"`
		Extra string `json:"@extra"`
		Text  struct {
			Type     string        `json:"@type"`
			Extra    string        `json:"@extra"`
			Text     string        `json:"text"`
			Entities []interface{} `json:"entities"`
		} `json:"text"`
		WebPage interface{} `json:"web_page"`
	} `json:"content"`
	ReplyMarkup interface{} `json:"reply_markup"`
}
