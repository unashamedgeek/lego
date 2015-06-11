package acme

import "time"

type registrationMessage struct {
	Contact []string `json:"contact"`
}

// Registration is returned by the ACME server after the registration
// The client implementation should save this registration somewhere.
type Registration struct {
	ID  int `json:"id"`
	Key struct {
		Kty string `json:"kty"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"key"`
	Recoverytoken string   `json:"recoveryToken"`
	Contact       []string `json:"contact"`
	Agreement     string   `json:"agreement,omitempty"`
}

// RegistrationResource represents all important informations about a registration
// of which the client needs to keep track itself.
type RegistrationResource struct {
	Body        Registration
	URI         string
	NewAuthzURL string
	TosURL      string
}

type authorizationResource struct {
	Body       authorization
	Domain     string
	NewCertURL string
}

type authorization struct {
	Identifier   identifier  `json:"identifier"`
	Status       string      `json:"status,omitempty"`
	Expires      time.Time   `json:"expires,omitempty"`
	Challenges   []challenge `json:"challenges,omitempty"`
	Combinations [][]int     `json:"combinations,omitempty"`
}

type identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type challenge struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	URI    string `json:"uri"`
	Token  string `json:"token"`
}