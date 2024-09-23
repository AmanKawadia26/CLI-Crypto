package models

// User represents a user in the system
type User struct {
	UserID                 int    `json:"userId"`
	Username               string `json:"username"`
	Password               string `json:"password"`
	Email                  string `json:"email"`
	Mobile                 int    `json:"mobile"`
	NotificationPreference string `json:"notificationPreference"`
	IsAdmin                bool   `json:"isAdmin"`
	Role                   string `json:"Role"`
}
