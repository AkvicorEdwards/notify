package tcp

const (
	// Unknown Error
	Err byte = 0

	// Message Request
	ReqMessage byte = 1
	// Message Response
	ResMessage byte = 2

	// Heartbeat Request
	ReqHeartbeat byte = 3
	// Heartbeat Response
	ResHeartbeat byte = 4

	// Registration request
	ReqRegister byte = 5
	// Registration Failure
	ResRegistrationFailure byte = 6
	// Registration Failure User Exist
	ResRegistrationFailureUserExist byte = 7
	// Registration Successful
	ResRegistrationSuccessful byte = 8

	// Notify Request
	ReqNotify byte = 20
	// Notify Response
	ResNotify byte = 21

	MaxMessageLength int = 4 * 200
)