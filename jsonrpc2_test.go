package go_jsonrpc2

import (
	"github.com/MattLaidlaw/go-assert"
	"log"
	"os"
	"testing"
)

type Class struct {}
func (c *Class) Method() {}

func TestHandler_ValidRequest(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res , err := client.Call("Class.Method")
	assert.ExpectEq(err, nil, t)

	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Error, Error{}, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_UnregisteredClass(t *testing.T) {
	client, err := Dial("localhost:6343")
	assert.ExpectEq(err, nil, t)

	res , err := client.Call("Class.Method")
	assert.ExpectEq(err, nil, t)

	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Error, Error{MethodNotFound, "unregistered object: Class", nil}, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_UnknownMethod(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.UnknownMethod")
	assert.ExpectEq(err, nil, t)

	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Error, Error{MethodNotFound, "unknown method: UnknownMethod", nil}, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_InvalidParams(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.Method", "extra parameter")
	assert.ExpectEq(err, nil, t)

	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Error, Error{InvalidParams, "given parameters do not match desired method", nil}, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestMain(m *testing.M) {
	// This server should register classes needed for testing.
	RegisteringServer := NewServer()
	RegisteringServer.Register(Class{})
	go func() {
		err := RegisteringServer.Listen("6342")
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// This server does not register any classes. It is intended to test the error response of calling an unregistered
	// method.
	NonRegisteringServer := NewServer()
	go func() {
		err := NonRegisteringServer.Listen("6343")
		if err != nil {
			log.Fatalln(err)
		}
	}()

	os.Exit(m.Run())
}