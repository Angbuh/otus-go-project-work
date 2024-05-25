package entities

type Note struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint64 `json:"user_id"`
}
