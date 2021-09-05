package jsonrpc2

import (
	"encoding/json"
	"github.com/google/uuid"
	"net"
)

type Client struct {
	conn net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
}

func Dial(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}, nil
}

func (c *Client) Call(method string, params ...interface{}) (*Response, error) {
	req := &Request{
		Version: JsonrpcVersion,
		Method:  method,
		Params:  params,
		Id:      uuid.NewString(),
	}

	err := c.encoder.Encode(req)
	if err != nil {
		return nil, err
	}

	res := new(Response)
	err = c.decoder.Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
