package serviceprovider

import (
	"encoding/xml"
)

type xmlUnauthorizedEvent struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}

type unauthorizedEventProcessor struct {
	*xmlEvent
}

func (impl serviceProviderImpl) newUnauthorizedEventProcessor(e *xmlEvent) EventProcessor {
	return &unauthorizedEventProcessor{
		xmlEvent: e,
	}
}

func (h unauthorizedEventProcessor) Process(data []byte) (string, error) {
	var xmlEvent xmlUnauthorizedEvent
	if err := xml.Unmarshal(data, &xmlEvent); err != nil {
		return "", err
	}

	return xmlEvent.AuthorizerAppid, nil

}
