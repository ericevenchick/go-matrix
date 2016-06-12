package matrix

type LoginRequest struct {
	Password string `json:"password"`
	Medium   string `json:"medium,omitempty"`
	Type     string `json:"type"`
	User     string `json:"user,omitempty"`
	Address  string `json:"address,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	HomeServer   string `json:"home_server"`
	UserID       string `json:"user_id"`
}

func (me *MatrixClient) PasswordLogin(user string, pass string) error {
	uri := me.server + "/_matrix/client/r0/login"
	req := LoginRequest{
		Password: pass,
		Medium:   "email",
		Type:     "m.login.password",
		User:     user,
	}

	var resp LoginResponse
	err := me.makeMatrixRequest("POST", uri, req, &resp)
	if err != nil {
		return err
	}

	me.accessToken = resp.AccessToken
	me.refreshToken = resp.RefreshToken

	return nil
}
