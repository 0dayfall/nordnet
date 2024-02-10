package feed

import (
	"encoding/json"
)

type PublicFeed struct {
	*Feed
}

func NewPublicFeed(address string) (*PublicFeed, error) {
	f, err := newFeed(address)
	if err != nil {
		return nil, err
	}

	return &PublicFeed{f}, err
}

// Starts reading from the connection and sends data through given channels
func (pf *PublicFeed) Dispatch(msgChan chan *PublicMsg, errChan chan error) {
	go func(d *json.Decoder, mc chan<- *PublicMsg, ec chan<- error) {
		var (
			pMsg *PublicMsg
			err  error
		)

		for {
			pMsg = new(PublicMsg)
			if err = d.Decode(pMsg); err != nil {
				ec <- err
			}
			msgChan <- pMsg
		}
	}(pf.decoder, msgChan, errChan)

	return
}
