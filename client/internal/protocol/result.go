package protocol

type Result string

const (
	SUCCESS             Result = "SUCCESS"
	USER_ALREADY_EXISTS Result = "USER_ALREADY_EXISTS"
	NO_SUCH_USER        Result = "NO_SUCH_USER"
	NO_SUCH_ROOM        Result = "NO_SUCH_ROOM"
	ROOM_ALREADY_EXISTS Result = "ROOM_ALREADY_EXISTS"
	NOT_INVITED         Result = "NOT_INVITED"
	NOT_JOINED          Result = "NOT_JOINED"
	NOT_IDENTIFIED      Result = "NOT_IDENTIFIED"
	INVALID             Result = "INVALID"
)
