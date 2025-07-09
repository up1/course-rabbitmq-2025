package shared

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (m *Message) String() string {
	return "ID: " + m.ID + ", Content: " + m.Content
}
