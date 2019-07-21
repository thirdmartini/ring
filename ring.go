/**
 * Definitions of ring json structures reeturned in various calls
 */
package ring

import "fmt"

const (
	API_BASE_URL = "https://api.ring.com"
	API_VERSON   = "9"
	API_PATH_SESSION    = "/clients_api/session"
	API_PATH_DINGS      = "/clients_api/dings/active"
	API_PATH_DEVICES    = "/clients_api/ring_devices"
	API_PATH_HISTORY    = "/clients_api/doorbots/history"
	API_PATH_RECORDINGS = "/clients_api/dings/%s/recording"
)

// Account Profile information
//    The important bit here is the AuthenticationToken which is valid for about 5 seconds from the time
//    of authentication
type Profile struct {
	Id                  int    `json:"id"`
	EMail               string `json:"email"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PhoneNumber         string `json:"phone_number"`
	AuthenticationToken string `json:"authentication_token"`
}

// Session information
type Session struct {
	Profile    Profile `json:"profile"`
	HardwareId string  `json:"hardware_id"`
	UserFlow   string  `json:"user_flow"`
}

// RecordingRef is a a partial Recording record returned as part of History
type RecordingRef struct {
	Status string `json:"status"`
}

// DoorbotRef is an incomplete Doorbot record ( Returned inside other structures )
type DoorbotRef struct {
	Id          uint64 `json:"id"`
	Description string `json:"description"`
}

// History table entry as part of API_PATH_HISTORY
type History struct {
	Id        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
	Answered  bool   `json:"answered"`
	//Events []string  `json:"events"` // not sure the format of this yet
	Kind        string       `json:"kind"`
	Favorite    bool         `json:"favorite"`
	SnapshotUrl string       `json:"snapshot_url"`
	Recording   RecordingRef `json:"recording"`
	Doorbot     DoorbotRef   `json:"doorbot"`
}

func (h History) String() string {
	return fmt.Sprintf("%s: %s %s %t", h.CreatedAt, h.Doorbot.Description, h.Kind, h.Answered)
}

// Doorbot ring dootbell structure returned as part of API_PATH_DEVICES
type Doorbot struct {
	Id                 uint64  `json:"id"`
	Description        string  `json:"description"`
	DeviceId           string  `json:"device_id"`
	TimeZone           string  `json:"time_zone"`
	Subscibed          bool    `json:"subscribed"`
	SubscribedMotions  bool    `json:"subscribed_motions"`
	BatteryLife        string  `json:"battery_life"`
	ExternalConnection bool    `json:"external_connection"`
	FirmwareVersion    string  `json:"firmware_version"`
	Kind               string  `json:"kind"`
	Latitide           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Address            string  `json:"address"`
	// There is more stuff... but this is the only pertinent stuff here
	//Settings Settings
	//
}

// Devices structure returned as part of API_PATH_DEVICES
//          not complete as I don;t have all of the device types ring sells
type Devices struct {
	Doorbots []Doorbot `json:"doorbots"`
	// Don;t know these yet
	//AuthorizedDoorbots []interface{}  `json:"authorized_doorbots"`
	//Chimes      []interface{}  `json:"chimes"`
	//StickupCams []interface{}   `json:"stickup_cams"`
	//BaseStations []interface{}   `json:"base_stations"`
}

// Ding structure returned during a active ding event
type Ding struct {
	Id                 uint64 `json:"id"`
	IdStr              string `json:"id_str"`
	State              string `json:"state"`
	Protocol           string `json:"protocol"`
	DoorbotId          uint64 `json:"doorbot_id"`
	DoorbotDescription string `json:"doorbot_description"`
	DeviceKind         string `json:"device_kind"`
	Motion             bool   `json:"motion"`
	SipServerAddress   string `json:"sip_server_ip"`
	SipPort            int    `json:"sip_server_port"`
	SipServerTLS       bool   `json:"sip_server_tls"`
	SipSessionId       string `json:"sip_session_id"`
	SipFrom            string `json:"sip_from"`
	SipTo              string `json:"sip_to"`
	SipToken           string `json:"sip_token"`
	SipDingId          string `json:"sip_ding_id"`
}

