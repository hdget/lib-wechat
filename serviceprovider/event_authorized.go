package serviceprovider

import (
	"encoding/xml"
)

type xmlAuthAuthorizedEvent struct {
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
}

type authorizedEventHandler struct {
	*xmlEvent
}

func (impl serviceProviderImpl) newAuthorizedEventProcessor(e *xmlEvent) EventProcessor {
	return &authorizedEventHandler{
		xmlEvent: e,
	}
}

func (h authorizedEventHandler) Process(data []byte) (string, error) {
	var xmlEvent xmlAuthAuthorizedEvent
	if err := xml.Unmarshal(data, &xmlEvent); err != nil {
		return "", err
	}

	return xmlEvent.AuthorizationCode, nil
}
