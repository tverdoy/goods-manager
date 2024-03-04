package entity

type Good struct {
	Id          int    `json:"id"`
	ProjectId   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"created_at"`
}
