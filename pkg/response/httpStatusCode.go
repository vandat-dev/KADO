package response

const (
	ErrCodeSuccess       = 2001  //Success
	ErrCodeInvalidParams = 2002  //Email invalid
	ErrInvalidToken      = 3001  //Token invalid
	ErrCodeUserHasExists = 50001 // User already exist
	ErrCodeUserNotFound  = 4000  // User not found
	ErrCodeInvalidLogin  = 4001  // Invalid login credentials
	ErrCodeAccessDenied  = 4003  // Access denied
	ErrCodeInternalError = 5000  // Internal server error
)

var msg = map[int]string{
	ErrCodeSuccess:       "Success",
	ErrInvalidToken:      "Token invalid",
	ErrCodeInvalidParams: "Email invalid",
	ErrCodeUserHasExists: "User already exist",
	ErrCodeUserNotFound:  "User not found",
	ErrCodeInvalidLogin:  "Invalid login credentials",
	ErrCodeAccessDenied:  "Access denied",
	ErrCodeInternalError: "Internal server error",
}

// GetMessage - Get message from error code
func GetMessage(errorCode int) string {
	if message, exists := msg[errorCode]; exists {
		return message
	}
	return "Unknown error"
}
