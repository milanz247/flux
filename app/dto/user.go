package dto

// UserDTO is the public shape of a user sent to the frontend.
type UserDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	CreatedAt string `json:"createdAt"`
}
