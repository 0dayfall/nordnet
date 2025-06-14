/*
Contains everything related to the public and private feeds
More information available on https://api.test.nordnet.se/next/2/api-docs/docs/feeds
*/
package feed

import (
	"crypto/tls"
	"encoding/json"
	"io"
)

// Represents the feed connection
type Feed struct {
	conn    io.ReadWriteCloser
	encoder *json.Encoder
	decoder *json.Decoder
}

// Returns a new Feed connected to the address specified
func newFeed(address string) (*Feed, error) {
	conn, err := tls.Dial("tcp", address, nil)
	if err != nil {
		return nil, err
	}

	return &Feed{conn, json.NewEncoder(conn), json.NewDecoder(conn)}, nil
}

// Feed implements the Writer interface
func (f *Feed) Write(any interface{}) error {
	return f.encoder.Encode(any)
}

// Feed implements the Closer interface
// closes the underlying conneciton
func (f *Feed) Close() error {
	return f.conn.Close()
}

// Send the login command with the specified session key
func (f *Feed) Login(session string, getState interface{}) error {
	return f.Write(&FeedCmd{Cmd: "login", Args: &LoginArgs{SessionKey: session, GetState: getState}})
}

// Used when sending feed commands
type FeedCmd struct {
	Cmd  string      `json:"cmd"`
	Args interface{} `json:"args"`
}

// Arguments for sending the login command
type LoginArgs struct {
	SessionKey string      `json:"session_key"`
	GetState   interface{} `json:"get_state,omitempty"`
}

// Arguments for getting orders and trades when logging in
type GetState struct {
	DeletedOrders bool  `json:"deleted_orders"`
	Days          int64 `json:"days,omitempty"`
}

// Used for receiving messages
type FeedMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Represents messages sent by the private feed
type PrivateMsg struct {
	Type string      `json:"t"`
	Data interface{} `json:"d"`
}

// Used in UnmarshalJSON overrides
type rawMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
