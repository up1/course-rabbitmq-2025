package shared

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Catalog string `json:"catalog,omitempty"` // Optional routing key
}

func (m *Message) String() string {
	return "ID: " + m.ID + ", Content: " + m.Content + ", Catalog: " + m.Catalog
}
