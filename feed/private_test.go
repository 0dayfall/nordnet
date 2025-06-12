package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/0dayfall/nordnet/util/models"
	"github.com/stretchr/testify/assert"
)

var privateUnmarshalTests = []struct {
	json     string
	expected *PrivateMsg
}{
	{
		`{
			"type":"heartbeat",
			"data":{}
		}`,
		&PrivateMsg{"heartbeat", struct{}{}},
	},
	{
		`{
			"type":"order",
			"data":{
				"accno":123,
				"order_id":123,
				"price":{
					"value":1.1,
					"currency":"test"
				},
				"volume":1.1,
				"tradable":{
					"identifier":"test",
					"market_id":123
				},
				"open_volume":1.1,
				"traded_volume":1.1,
				"side":"test",
				"modified":123,
				"reference":"test",
				"activation_condition":{
					"type":"test",
					"trailing_value":1.1,
					"trigger_value":1.1,
					"trigger_condition":"test"
				},
				"price_condition":"test",
				"volume_condition":"test",
				"validity":{
					"type":"test",
					"valid_until":123
				},
				"action_state":"test",
				"order_state":"test"
			}
		}`,
		&PrivateMsg{"order", PrivateOrder{
			Accno:               123,
			OrderId:             123,
			Price:               models.Amount{1.1, "test"},
			Volume:              1.1,
			Tradable:            models.TradableId{"test", 123},
			OpenVolume:          1.1,
			TradedVolume:        1.1,
			Side:                "test",
			Modified:            123,
			Reference:           "test",
			ActivationCondition: models.ActivationCondition{"test", 1.1, 1.1, "test"},
			PriceCondition:      "test",
			VolumeCondition:     "test",
			Validity:            models.Validity{"test", 123},
			ActionState:         "test",
			OrderState:          "test",
		}},
	},
	{
		`{
			"type": "trade",
			"data": {
				"accno": 123,
				"order_id": 123,
				"trade_id": "test",
				"tradable": {
					"identifier": "test",
					"market_id": 123
				},
				"price": {
					"value": 1.1,
					"currency": "test"
				},
				"volume": 1.1,
				"side": "test",
				"counterparty": "test",
				"tradetime": 123
			}
		}`,
		&PrivateMsg{"trade", PrivateTrade{
			Accno:        123,
			OrderId:      123,
			TradeId:      "test",
			Tradable:     models.TradableId{"test", 123},
			Price:        models.Amount{1.1, "test"},
			Volume:       1.1,
			Side:         "test",
			Counterparty: "test",
			Tradetime:    123,
		}},
	},
}

func TestPrivateMsgUnmarshalJSON(t *testing.T) {
	for i, tt := range privateUnmarshalTests {
		msg := &PrivateMsg{}
		err := json.Unmarshal([]byte(tt.json), msg)
		assert.NoError(t, err, "case %d: failed to unmarshal", i)
		assert.Equal(t, tt.expected.Type, msg.Type, "case %d: unexpected type", i)
		assert.Equal(t, tt.expected.Data, msg.Data, "case %d: unexpected data", i)
	}
}

var privateDispatchTests = []struct {
	json     string
	expected *PrivateMsg
}{
	{
		`{"type":"heartbeat","data":{}}`,
		&PrivateMsg{"heartbeat", struct{}{}},
	},
	{
		`{"type":"trade","data":{}}`,
		&PrivateMsg{"trade", PrivateTrade{}},
	},
	{
		`{"type":"order","data":{}}`,
		&PrivateMsg{"order", PrivateOrder{}},
	},
}

func TestPrivateFeedDispatch(t *testing.T) {
	var input bytes.Buffer
	for _, tt := range privateDispatchTests {
		input.WriteString(tt.json + "\n")
	}

	conn := &fakeConnection{Buffer: &input}
	feed := &PrivateFeed{
		&Feed{
			conn,
			json.NewEncoder(&bytes.Buffer{}), // unused encoder
			json.NewDecoder(conn),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msgChan := make(chan *PrivateMsg, len(privateDispatchTests))
	errChan := make(chan error, 1)

	go feed.Dispatch(ctx, msgChan, errChan)

	for i, tt := range privateDispatchTests {
		select {
		case msg := <-msgChan:
			assert.Equal(t, tt.expected.Type, msg.Type, "case %d: wrong type", i)
			assert.IsType(t, tt.expected.Data, msg.Data, "case %d: wrong data type", i)
		case err := <-errChan:
			t.Fatalf("case %d: error from Dispatch: %v", i, err)
		}
	}
}
