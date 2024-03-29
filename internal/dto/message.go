package dto

type Message struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attachments []string `json:"attachments,omitempty"`
}
