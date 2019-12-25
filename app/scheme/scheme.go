package scheme

// TelMessage is TDLib message struct
type MessageContentEntry struct {
	Content struct {
		Type string `json:"@type"`
		Text struct {
			Text string `json:"text"`
		} `json:"text"`
	} `json:"content"`
}
