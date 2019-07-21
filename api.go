package ring

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	TokenExpiredError = "authentication token expired"
	PermissionError   = "no permissions"
	RateLimited       = "api rate limit exceeded"
)

var hardwareId = "A9118CAB-A774-40B7-9A83-EA16AE901B6F"

type Api struct {
	username string
	password string
}

func checkError(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201:
		break

	case 401:
		return errors.New(TokenExpiredError)

	case 403:
		return errors.New(PermissionError)

	case 429:
		return errors.New(RateLimited)

	default:
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *Api) getRaw(kind string, path string, vals *url.Values) ([]byte, error) {
	// Note that it is possible to use the same session.Profile.AuthenticationToken
	//   for multiple requests.. as long as they are all done within an apparent 5 second window
	//   ( token seems to only remain valid for about 5 seconds ) after which we will get a 429
	// So instead we just get a new token for each request
	session, err := a.authenticate()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	v := vals
	if v == nil {
		v = &url.Values{}
	}

	// Add API token and auth
	v.Set("api_version", API_VERSON)
	v.Set("auth_token", session.Profile.AuthenticationToken)

	req, err := http.NewRequest(kind, path+"?"+v.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkError(resp)
	if err != nil {
		return nil, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyText, nil
}


func (a *Api) get(kind string, url string, vals *url.Values, inter interface{}) error {
	body, err := a.getRaw(kind, url, vals)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, inter)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) authenticate() (*Session, error) {
	return a.getToken()
}

// Session returns Session information
func (a *Api) Session() (*Session, error) {
	session, err := a.authenticate()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Profile returns the account profile
func (a *Api) Profile() (*Profile, error) {
	session, err := a.authenticate()
	if err != nil {
		return nil, err
	}

	return &session.Profile, nil
}

// Devices returns a structure containing lists of various devices by type
func (a *Api) Devices() (*Devices, error) {
	devices := Devices{}

	err := a.get("GET", API_BASE_URL+API_PATH_DEVICES, nil, &devices)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}

// History return up to max number of history entries
func (a *Api) History(max int) ([]History, error) {
	history := make([]History, 0, 10)

	v := url.Values{}
	v.Set("limit", strconv.Itoa(max))

	err := a.get("GET", API_BASE_URL+API_PATH_HISTORY, &v, &history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// Recording downloads the recording identified by id and saves it to the saveFile
//       id is the id of the event as provided in History{}
func (a *Api) Recording(id uint64, saveFile string) error {
	session, err := a.authenticate()
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(saveFile)
	if err != nil {
		return err
	}
	defer out.Close()

	v := url.Values{}
	v.Set("api_version", API_VERSON)
	v.Set("auth_token", session.Profile.AuthenticationToken)

	link := fmt.Sprintf(API_BASE_URL+API_PATH_RECORDINGS+"?%s", strconv.FormatUint(id, 10), v.Encode())

	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Listen listens for a new Ding event every pollInterval duration and calls OnDing when an active event is processing
//         pollInterval should be greater than 17 seconds otherwise we are likely to hit an 429 throttle event from
//         cloudflare
func (a *Api) Listen(pollInterval time.Duration, OnDing func(ding *Ding)) error {
	for true {
		dings := make([]Ding, 0, 1)

		err := a.get("GET", API_BASE_URL+API_PATH_DINGS, nil, &dings)
		if err != nil {
			return err
		}

		for _, ding := range dings {
			OnDing(&ding)
		}

		time.Sleep(pollInterval)
	}
	return nil
}

// New makes a new Ring Api object
func New(username, password string) (*Api, error) {
	a := &Api{
		username: username,
		password: password,
	}
	return a, nil
}
