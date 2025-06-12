package feed

import (
	"context"
	"encoding/json"

	"github.com/0dayfall/nordnet/util/models"
)

type PrivateFeed struct {
	*Feed
}

func NewPrivateFeed(address string) (*PrivateFeed, error) {
	f, err := newFeed(address)
	if err != nil {
		return nil, err
	}
	return &PrivateFeed{f}, nil
}

// Starts reading from the connection and sends data through given channels
func (pf *PrivateFeed) Dispatch(ctx context.Context, msgChan chan *PrivateMsg, errChan chan error) {
	go func() {
		defer close(msgChan)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				pMsg := new(PrivateMsg)
				if err := pf.decoder.Decode(pMsg); err != nil {
					errChan <- err
					return
				}
				msgChan <- pMsg
			}
		}
	}()
}

// Order data section in the private message
type PrivateOrder models.Order

// Trade data section in the private message
type PrivateTrade models.Trade

// Implements the Unmarshaler interface
// decodes the json into proper data types depending on the type field
func (pm *PrivateMsg) UnmarshalJSON(b []byte) (err error) {
	rawMsg := rawMsg{} // to avoid endless recursion below
	if err = json.Unmarshal(b, &rawMsg); err != nil {
		return
	}

	*pm = PrivateMsg{Type: rawMsg.Type}

	switch rawMsg.Type {
	case heartbeatType:
		pm.Data = struct{}{}
	case orderType:
		order := PrivateOrder{}
		if err = json.Unmarshal(rawMsg.Data, &order); err != nil {
			return
		}
		pm.Data = order
	case tradeType:
		trade := PrivateTrade{}
		if err = json.Unmarshal(rawMsg.Data, &trade); err != nil {
			return
		}
		pm.Data = trade
	}

	return
}
