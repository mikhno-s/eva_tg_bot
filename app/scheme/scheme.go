package scheme

// TelMessage is TDLib message struct
type MessageContentEntry struct {
	Content struct {
		Text struct {
			Text string `json:"text"`
		} `json:"text"`
	} `json:"content"`
}
