package github

import "time"

// PR represents a GitHub Pull Request
type PR struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Body        string    `json:"body"`
	IsDraft     bool      `json:"isDraft"`
	HeadRefName string    `json:"headRefName"`
	Owner       string    `json:"owner"`
	Repo        string    `json:"repo"`
	TicketID    string    `json:"ticketId"`
}
