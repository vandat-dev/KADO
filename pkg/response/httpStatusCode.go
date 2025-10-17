package response

const (
	ErrCodeSuccess              = 2001  //Success
	ErrCodeInvalidParams        = 2002  //Email invalid
	ErrInvalidToken             = 3001  //Token invalid
	ErrCodeUserHasExists        = 50001 // User already exist
	ErrCodeUserNotFound         = 4000  // User not found
	ErrCodeInvalidLogin         = 4001  // Invalid login credentials
	ErrCodeAccessDenied         = 4003  // Access denied
	ErrCodeAccountLock          = 4004  // Your account has been locked
	ErrCodeUserPermissionDenied = 4005  // You do not have permission to interact with this user
	ErrCodeInternalError        = 5000  // Internal server error
	ErrCodeInvalidData          = 4221  // Invalid request data
	ErrCodeUnauthorized         = 4010  // Unauthorized

	// task
	ErrCodeTaskExists           = 50101 // Task already exists
	ErrCodeTaskNotFound         = 4100  // Task not found
	ErrCodeTaskPermissionDenied = 40301 // You do not have permission to interact with this event
	ErrCodeTaskExportFailed     = 40302 // Failed to export tasks

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
		ErrCodeSuccess:       "SUCCESS",
		ErrInvalidToken:      "TOKEN_INVALID",
		ErrCodeInvalidLogin:  "LOGIN_FAILED",
		ErrCodeAccessDenied:  "ACCESS_DENIED",
		ErrCodeInternalError: "INTERNAL_SERVER_ERROR",
		ErrCodeInvalidData:   "INVALID_DATA",
		ErrCodeUnauthorized:  "UNAUTHORIZED",

		//	user
		ErrCodeInvalidParams:        "EMAIL_INVALID",
		ErrCodeUserHasExists:        "USER_ALREADY_EXISTS",
		ErrCodeUserNotFound:         "USER_NOT_FOUND",
		ErrCodeAccountLock:          "USER_ACCOUNT_LOCKED",
		ErrCodeUserPermissionDenied: "YOU_DO_NOT_HAVE_PERMISSION_TO_INTERACT_WITH_THIS_USER",

		//	task
		ErrCodeTaskNotFound:         "TASK_NOT_FOUND",
		ErrCodeTaskExists:           "TASK_ALREADY_EXISTS",
		ErrCodeTaskPermissionDenied: "YOU_DO_NOT_HAVE_PERMISSION_TO_INTERACT_WITH_THIS_TASK",
		ErrCodeTaskExportFailed:     "FAILED_TO_EXPORT_TASKS",

		//	client
		ErrCodeClientNotFound: "CLIENT_NOT_FOUND",
		ErrCodeClientExists:   "CLIENT_ALREADY_EXISTS",

		//	job
		ErrCodeJobNotFound: "JOB_NOT_FOUND",
		ErrCodeJobExists:   "JOB_ALREADY_EXISTS",

		//	role
		ErrCodeRoleNotFound: "ROLE_NOT_FOUND",
		ErrCodeRoleExists:   "ROLE_ALREADY_EXISTS",

		//	item
		ErrCodeItemNotFound: "ITEM_NOT_FOUND",
		ErrCodeItemExists:   "ITEM_ALREADY_EXISTS",
	}
)

// GetMessage - Get message from error code
func GetMessage(errorCode int) string {
	if message, exists := msg[errorCode]; exists {
		return message
	}
	return "Unknown error"
}
