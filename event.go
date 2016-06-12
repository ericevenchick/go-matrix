package matrix

import (
	"fmt"
)

type SendEventResponse struct {
	EventID string `json:"event_id"`
}

type MessageEvent struct {
	Body        string `json:"body"`
	MessageType string `json:"msgtype"`
}

func (me *MatrixClient) SendEvent(roomID string, eventType string, event interface{}) error {
	uri := (me.server + "/_matrix/client/r0/rooms/" + roomID + "/send/" +
		eventType + fmt.Sprintf("/m%d", me.transactionID) +
		"?access_token=" + me.accessToken)
	me.transactionID += 1

	var resp SendEventResponse
	err := me.makeMatrixRequest("PUT", uri, event, &resp)

	if err != nil {
		return err
	}

	return nil
}
