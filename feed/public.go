package feed

import (
	"context"
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

func (pf *PublicFeed) Dispatch(ctx context.Context, msgChan chan *PublicMsg, errChan chan error) {
	go func(d *json.Decoder, mc chan<- *PublicMsg, ec chan<- error) {
		defer close(mc)
		defer close(ec)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				pMsg := new(PublicMsg)
				if err := d.Decode(pMsg); err != nil {
					ec <- err
					return
				}
				mc <- pMsg
			}
		}
	}(pf.decoder, msgChan, errChan)
}
