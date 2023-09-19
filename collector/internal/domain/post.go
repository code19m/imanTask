package domain

type Post struct {
	Id     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type CollectPostPayload struct {
	StartPage int32 `json:"start_page"`
}
