package protocol

import "encoding/json"

func Serialize(msg ClientMessage) (string, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Deserialize(data string) (ServerMessage, error) {
	var msg ServerMessage
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		return ServerMessage{}, err
	}
	return msg, nil
}
