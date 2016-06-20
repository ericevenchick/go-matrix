package matrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type MatrixClient struct {
	server        url.URL
	endpoints     endpoints
	accessToken   string
	refreshToken  string
	transactionID int
	client        http.Client
	nextBatch     string
}

type endpoints struct {
	login     url.URL
	room      url.URL
	sync      url.URL
	sendEvent url.URL
}

type ErrorResponse struct {
	ErrorCode    string `json:"errcode"`
	ErrorMessage string `json:"error"`
}

func NewClient(server string) (*MatrixClient, error) {
	client := http.Client{}
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	return &MatrixClient{
		server: *u,
		endpoints: endpoints{
			login: url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/_matrix/client/r0/login"},
			room:  url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/_matrix/client/r0/rooms/"},
			sync:  url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/_matrix/client/r0/sync"},
		},
		transactionID: 0,
		client:        client,
		nextBatch:     "",
	}, nil
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
	uri := me.endpoints.room
	uri.Path += path.Join(roomID, "join")
	params := url.Values{}
	params.Add("access_token", me.accessToken)
	uri.RawQuery = params.Encode()

	err := me.makeMatrixRequest("POST", uri.String(), nil, nil)
	if err != nil {
		return err
	}

	return nil
}
