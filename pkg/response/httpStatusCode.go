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
	ErrCodeInvalidData   = 4221  // Invalid request data
	ErrCodeUnauthorized  = 4010  // Unauthorized

	// task
	ErrCodeTaskExists           = 50101 // Task already exists
	ErrCodeTaskNotFound         = 4100  // Task not found
	ErrCodeTaskPermissionDenied = 40301 // You do not have permission to interact with this event

	// client
	ErrCodeClientNotFound = 4200  // Client not found
	ErrCodeClientExists   = 50201 // Client already exists

	// job
	ErrCodeJobNotFound = 4300  // Job not found
	ErrCodeJobExists   = 50301 // Job already exists

	// role
	ErrCodeRoleNotFound = 4400  // Role not found
	ErrCodeRoleExists   = 50401 // Role already exists

	// item
	ErrCodeItemNotFound = 4500  // Item not found
	ErrCodeItemExists   = 50501 // Item already exists
)

var (
	msg = map[int]string{
		//	common
		ErrCodeSuccess:       "Success",
		ErrInvalidToken:      "Token invalid",
		ErrCodeInvalidLogin:  "Invalid login credentials",
		ErrCodeAccessDenied:  "Access denied",
		ErrCodeInternalError: "Internal server error",
		ErrCodeInvalidData:   "Invalid request data",
		ErrCodeUnauthorized:  "Unauthorized",

		//	user
		ErrCodeInvalidParams: "Email invalid",
		ErrCodeUserHasExists: "User already exist",
		ErrCodeUserNotFound:  "User not found",

		//	task
		ErrCodeTaskNotFound:         "Task not found",
		ErrCodeTaskExists:           "Task already exists",
		ErrCodeTaskPermissionDenied: "You do not have permission to interact with this event",

		//	client
		ErrCodeClientNotFound: "Client not found",
		ErrCodeClientExists:   "Client already exists",

		//	job
		ErrCodeJobNotFound: "Job not found",
		ErrCodeJobExists:   "Job already exists",

		//	role
		ErrCodeRoleNotFound: "Role not found",
		ErrCodeRoleExists:   "Role already exists",

		//	item
		ErrCodeItemNotFound: "Item not found",
		ErrCodeItemExists:   "Item already exists",
	}
)

// GetMessage - Get message from error code
func GetMessage(errorCode int) string {
	if message, exists := msg[errorCode]; exists {
		return message
	}
	return "Unknown error"
}
