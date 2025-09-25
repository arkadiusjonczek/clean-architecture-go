package common

const (
	DEFAULT_USER_ID = "1337"
)

// GetUserID returns a static user id
func GetUserID() string {
	return DEFAULT_USER_ID
}
