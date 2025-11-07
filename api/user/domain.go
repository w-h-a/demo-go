package user

// User represents the data model for a user in the database.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserDTO (Data Transfer Object) is used to capture
// the request body when creating a new user.
type CreateUserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
