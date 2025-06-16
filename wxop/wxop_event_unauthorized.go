package wxop

import (
	"encoding/xml"
)

type xmlUnauthorizedEvent struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}

type unauthorizedEventProcessor struct {
	*eventImpl
}

func (impl wxopImpl) newUnauthorizedEventProcessor(e *eventImpl) EventProcessor {
	return &unauthorizedEventProcessor{
		eventImpl: e,
	}
}

func (h unauthorizedEventProcessor) Process(data []byte) (string, error) {
	var xmlEvent xmlUnauthorizedEvent
	if err := xml.Unmarshal(data, &xmlEvent); err != nil {
		return "", err
	}

	return xmlEvent.AuthorizerAppid, nil

}
