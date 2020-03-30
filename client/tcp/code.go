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
	ResRegister byte = 7
	// Registration Failure
	ResRegistrationFailure byte = 8
	// Registration Failure User Exist
	ResRegistrationFailureUserExist byte = 9
	// Registration Successful
	ResRegistrationSuccessful byte = 10

	ReqKey byte = 11
	ResKey byte = 12

	// Notify Request
	ReqNotify byte = 20
	// Notify Response
	ResNotify byte = 21

	MaxMessageLength int = 4 * 200
)