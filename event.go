package matrix

import (
	"net/url"
	"path"
	"strconv"
)

type SendEventResponse struct {
	EventID string `json:"event_id"`
}

type MessageEvent struct {
	Body        string `json:"body"`
	MessageType string `json:"msgtype"`
}

type Sync struct {
	NextBatch string        `json:"next_batch"`
	Rooms     AllRoomEvents `json:"rooms"`
	Presence  Presence      `json:"presence"`
}
type AllRoomEvents struct {
	Invite map[string]Room `json:"invite"`
	Join   map[string]Room `json:"join"`
	Leave  map[string]Room `json:"leave"`
}
type Room struct {
	AccountData         AccountData              `json:"account_data"`
	State               State                    `json:"state"`
	Timeline            Timeline                 `json:"timeline"`
	Ephemeral           Ephemeral                `json:"ephemeral"`
	UnreadNotifications UnreadNotificationCounts `json:"unread_notifications"`
	InviteState         InviteState              `json:"invite_state"`
}
type Timeline struct {
	Limited       bool    `json:"limited"`
	Events        []Event `json:"events"`
	PreviousBatch string  `json:"prev_batch"`
}
type State struct {
	Events []Event `json:"events"`
}
type AccountData struct {
	Events []Event `json:"events"`
}
type Ephemeral struct {
	Events []Event `json:"events"`
}
type InviteState struct {
	Events []Event `json:"events"`
}
type Presence struct {
	Events []Event `json:"events"`
}
type UnreadNotificationCounts struct {
	HighlightCount    int `json:"highlight_count"`
	NotificationCount int `json:"notification_count"`
}
type Event struct {
	Type             string       `json:"type"`
	ID               string       `json:"event_id"`
	Sender           string       `json:"sender"`
	StateKey         string       `json:"state_key"`
	OriginServerTime int          `json:"origin_server_ts"`
	Content          EventContent `json:"content"`
	Unsigned         Unsigned     `json:"unsigned"`
}
type Unsigned struct {
	PreviousContent EventContent `json:"prev_content,omitempty"`
	Age             int          `json:"age"`
	TransactionID   string       `json:"transaction_id"`
}

type EventContent struct {
	MessageType string `json:"msg_type"`
	Body        string `json:"body"`
}

func (me *MatrixClient) SendEvent(roomID string, eventType string, event interface{}) error {
	uri := me.endpoints.room
	uri.Path += path.Join(roomID, "send", eventType, "m"+strconv.Itoa(me.transactionID))
	params := url.Values{}
	params.Add("access_token", me.accessToken)
	uri.RawQuery = params.Encode()

	me.transactionID += 1

	var response SendEventResponse
	err := me.makeMatrixRequest("PUT", uri.String(), event, &response)

	if err != nil {
		return err
	}

	return nil
}

func (me *MatrixClient) SyncOnce() (Sync, error) {
	uri := me.endpoints.sync
	params := url.Values{}
	params.Add("access_token", me.accessToken)
	params.Add("timeout", "10000")
	uri.RawQuery = params.Encode()

	if me.nextBatch != "" {
		uri.Query().Add("since", me.nextBatch)
	}

	var response Sync
	err := me.makeMatrixRequest("GET", uri.String(), nil, &response)

	// update nextBatch so we only sync new messages next time
	me.nextBatch = response.NextBatch

	if err != nil {
		return Sync{}, err
	}
	return response, nil
}

func (me *MatrixClient) StartSync() chan Sync {
	ch := make(chan Sync)

	go func() {
		for {
			sync, _ := me.SyncOnce()
			ch <- sync
		}
	}()

	return ch
}

func (me *Room) GetEvents() []Event {
	events := []Event{}
	for _, event := range me.Timeline.Events {
		events = append(events, event)
	}
	return events
}
