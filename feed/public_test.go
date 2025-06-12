package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Unmarshal Tests ---

var publicUnmarshalTests = []struct {
	json     string
	expected *PublicMsg
}{
	{
		`{"type":"heartbeat","data":{}}`,
		&PublicMsg{"heartbeat", struct{}{}},
	},
	{
		`{"type":"price","data":{"i":"test","m":123,"trade_timestamp":123,"tick_timestamp":123,"bid":1.1,"bid_volume":1.1,"ask":1.1,"ask_volume":1.1,"close":1.1,"high":1.1,"last":1.1,"last_volume":1.1,"low":1.1,"open":1.1,"turnover":1.1,"turnover_volume":1.1,"ep":1.1,"paired":1.1,"imbalance":1.1}}`,
		&PublicMsg{"price", PublicPrice{I: "test", M: 123, TradeTimestamp: 123, TickTimestamp: 123, Bid: 1.1, BidVolume: 1.1, Ask: 1.1, AskVolume: 1.1, Close: 1.1, High: 1.1, Last: 1.1, LastVolume: 1.1, Low: 1.1, Open: 1.1, Turnover: 1.1, TurnoverVolume: 1.1, EP: 1.1, Paired: 1.1, Imbalance: 1.1}},
	},
	// Add more cases like trade, depth, etc. (already tested individually below)
}

func TestPublicMsgUnmarshalJSON(t *testing.T) {
	for i, tt := range publicUnmarshalTests {
		msg := &PublicMsg{}
		err := json.Unmarshal([]byte(tt.json), msg)
		assert.NoError(t, err, "case %d: unmarshal failed", i)
		assert.Equal(t, tt.expected.Type, msg.Type, "case %d: type mismatch", i)
		assert.IsType(t, tt.expected.Data, msg.Data, "case %d: data type mismatch", i)
	}
}

// --- Dispatch Tests ---

var publicDispatchTests = []struct {
	json     string
	expected *PublicMsg
}{
	{`{"type":"heartbeat","data":{}}`, &PublicMsg{"heartbeat", struct{}{}}},
	{`{"type":"price","data":{}}`, &PublicMsg{"price", PublicPrice{}}},
	{`{"type":"trade","data":{}}`, &PublicMsg{"trade", PublicTrade{}}},
	{`{"type":"depth","data":{}}`, &PublicMsg{"depth", PublicDepth{}}},
	{`{"type":"indicator","data":{}}`, &PublicMsg{"indicator", PublicIndicator{}}},
	{`{"type":"news","data":{}}`, &PublicMsg{"news", PublicNews{}}},
	{`{"type":"trading_status","data":{}}`, &PublicMsg{"trading_status", PublicTradingStatus{}}},
}

func TestPublicFeedDispatch(t *testing.T) {
	var input bytes.Buffer
	for _, tt := range publicDispatchTests {
		input.WriteString(tt.json + "\n")
	}

	conn := &fakeConnection{Buffer: &input}
	feed := &PublicFeed{
		&Feed{
			conn,
			json.NewEncoder(&bytes.Buffer{}),
			json.NewDecoder(conn),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msgChan := make(chan *PublicMsg, len(publicDispatchTests))
	errChan := make(chan error, 1)

	feed.Dispatch(ctx, msgChan, errChan)

	for i, tt := range publicDispatchTests {
		select {
		case msg := <-msgChan:
			assert.Equal(t, tt.expected.Type, msg.Type, "case %d: type mismatch", i)
			assert.IsType(t, tt.expected.Data, msg.Data, "case %d: data type mismatch", i)
		case err := <-errChan:
			t.Fatalf("case %d: unexpected error from Dispatch: %v", i, err)
		case <-ctx.Done():
			t.Fatalf("case %d: test timed out", i)
		}
	}
}
