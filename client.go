package matrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type MatrixClient struct {
	server        string
	accessToken   string
	refreshToken  string
	transactionID int
	client        http.Client
}

type ErrorResponse struct {
	ErrorCode    string `json:"errcode"`
	ErrorMessage string `json:"error"`
}

func NewClient(server string) MatrixClient {
	client := http.Client{}
	return MatrixClient{
		server:        server,
		transactionID: 0,
		client:        client,
	}
}

func (me *MatrixClient) makeMatrixRequest(method string, uri string, reqIf interface{}, respIf interface{}) error {
	// create the JSON request
	reqBody, err := json.Marshal(reqIf)
	if err != nil {
		return err
	}
	reqBuf := bytes.NewBuffer(reqBody)

	// send the request
	req, err := http.NewRequest(method, uri, reqBuf)
	req.Header.Add("Content-Type", "application/json")
	resp, err := me.client.Do(req)

	if err != nil {
		return err
	}

	// retrieve the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// decode JSON as error response to see if an error occurred
	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if errResp.ErrorCode != "" {
		return errors.New(errResp.ErrorMessage)
	}

	err = json.Unmarshal(body, &respIf)
	if err != nil {
		return err
	}

	return nil
}

func (me *MatrixClient) JoinRoom(roomID string) error {
	uri := me.server + "/_matrix/client/r0/rooms/" + roomID + "/join?access_token=" + me.accessToken
	err := me.makeMatrixRequest("POST", uri, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
