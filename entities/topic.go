package entities

type Topic struct {
	Id    *int    `json:"id"`
	Title string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
