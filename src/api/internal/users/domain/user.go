package domain

// User represents a user in the system with personal and authentication details
type User struct {
	// Unique identifier for the user (UUID)
	ID string `json:"uuid" gorm:"type:nvarchar(36);primaryKey"`
	// User's name (required)
	Name string `json:"name" gorm:"type:nvarchar(255);not null"`
	// User's email (unique and required)
	Email string `json:"email" gorm:"type:nvarchar(255);unique;not null"`
	// User's hashed password
	Password string `json:"password" gorm:"type:nvarchar(64);not null"`
	// Profile picture (optional)
	Picture string `json:"picture" gorm:"type:nvarchar(max)"`
	// User's phone number (required)
	Phone string `json:"phone" gorm:"type:nvarchar(20)"`
	// User's role (admin or regular user)
	Role bool `json:"role" gorm:"type:bit;not null;default:0"`
}

// The functions above were made in case that the outputed result for role should be a string and not boolean
// Uncomment in case that's the purpose

// func (u *User) MarshalJSON() ([]byte, error) {
// 	type Alias User
// 	return json.Marshal(&struct {
// 		Role string `json:"role"`
// 		*Alias
// 	}{
// 		Role:  boolToString(u.Role),
// 		Alias: (*Alias)(u),
// 	})
// }

// func boolToString(b bool) string {
// 	if b {
// 		return "admin"
// 	}
// 	return "user"
// }
