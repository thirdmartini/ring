package ring

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}


type Metadata struct {
	ApiVersion string `json:"api_version"`
}

type Device struct {
	HardwareId string   `json:"hardware_id"`
	Metadata   Metadata `json:"metadata"`
	OS         string   `json:"os"`
}

type TokenRequest struct {
	Device Device `json:"device"`
}

func (a *Api) getToken() (*Session, error) {
	req := struct {
		ClientID  string `json:"client_id"`
		GrantType string `json:"grant_type"`
		Scope     string `json:"scope"`
		Password  string `json:"password"`
		Username  string `json:"username"`
	}{
		ClientID:  "ring_official_android",
		GrantType: "password",
		Scope:     "client",
		Username:  a.username,
		Password:  a.password,
	}

	token := Token{}

	err := a.do("POST", "https://oauth.ring.com/oauth/token", "", &req, &token)
	if err != nil {
		return nil, err
	}

	req2 := TokenRequest{
		Device: Device{
			HardwareId: hardwareId,
			OS:         "android",
			Metadata: Metadata{
				ApiVersion: "9",
			},
		},
	}

	session := &Session{}
	err = a.do("POST", "https://api.ring.com/clients_api/session?api_version=9", token.AccessToken, &req2, session)
	return session, err
}

func (a *Api) do(kind string, url string, token string, in interface{}, out interface{}) error {
	b := new(bytes.Buffer)
	if in != nil {
		json.NewEncoder(b).Encode(in)
	}

	req, err := http.NewRequest(kind, url, b)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		if out != nil {
			err = json.Unmarshal(body, out)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New(string(body))
}
