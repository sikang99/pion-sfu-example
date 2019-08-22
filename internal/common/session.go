package common

// SessionInfo desscrbes session information of user login
type SessionInfo struct {
	UserID    string
	LoginTime string
	Gender    string
	Age       int
}

// GetSession gets user session information
func GetSession(sessionid string) (SessionInfo, error) {
	sess := SessionInfo{UserID: "Stoney"}
	return sess, nil
}

// RoomInfo means information record of video room
type RoomInfo struct {
	RoomID     string
	CreateTime string
	UseTime    string
}
