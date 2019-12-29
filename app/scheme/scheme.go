package scheme

import "time"

// MessageContentEntry is TDLib message struct with text field
type MessageContentEntry struct {
	ID      int64 `json:"id"`
	Date    int64 `json:"date"`
	Content struct {
		Type string `json:"@type"`
		Text struct {
			Text string `json:"text"`
		} `json:"text"`
	} `json:"content"`
}

// Car evacuated entity
type Car struct {
	Date          time.Time `db:"date"           json:"date"`
	Model         string    `db:"model"          json:"model"`
	LicensePlate  string    `db:"license_plate"  json:"license_plate"`
	VIN           string    `db:"vin"            json:"vin"`
	ID            string    `db:"entity_id"      json:"id"`
	RetriavedFrom int64
}

// GetEvacuatedTableScheme returns evacuated_cars create query
func GetEvacuatedTableScheme() string {
	scheme := `
	CREATE TABLE IF NOT EXISTS evacuated_cars (
		date timestamp(0),
		model text, 
		license_plate text,
		vin text,
		entity_id text UNIQUE,
		id serial,
		PRIMARY KEY (id, license_plate)
	  );
	`
	return scheme
}

// GetTelegramLastChatMessagesTableScheme returns telegram_last_chat_messages create query
func GetTelegramLastChatMessagesTableScheme() string {
	scheme := `
	CREATE TABLE IF NOT EXISTS telegram_last_chat_messages (
		chat_id text UNIQUE,
		last_message_id text,
		PRIMARY KEY (chat_id)
	  );
	`
	return scheme
}
