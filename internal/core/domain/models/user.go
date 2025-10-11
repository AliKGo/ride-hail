package models

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Password  string `json:"password"`
	Attrs     struct{}
}
