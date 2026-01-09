package event

import (
	"encoding/xml"
)

type xmlUnauthorizedEvent struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}

type unauthorizedEventHandler struct {
}

func newUnauthorizedEventHandler() Handler {
	return &unauthorizedEventHandler{}
}

func (h unauthorizedEventHandler) Handle(data []byte) (string, error) {
	var e xmlUnauthorizedEvent
	if err := xml.Unmarshal(data, &e); err != nil {
		return "", err
	}

	return e.AuthorizerAppid, nil

}
