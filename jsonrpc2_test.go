package jsonrpc2

import (
	"github.com/MattLaidlaw/go-assert"
	"log"
	"os"
	"testing"
)

type Class struct {}

func (c *Class) Method() {}

func (c *Class) ArgMethod(n float64) float64 {
	x := 42 + n
	return x
}

func (c *Class) ReturnFloat64() float64 {
	return float64(42)
}

func (c *Class) ReturnString() string {
	return "Hello, World!"
}

func (c *Class) ReturnBool() bool {
	return true
}

func TestHandler_ValidRequest(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res , err := client.Call("Class.Method")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{}, t)

	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_Args(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.ArgMethod", 42)
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{}, t)

	assert.ExpectEq(res.Result, float64(84), t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_ReturnFloat64(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.ReturnFloat64")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{}, t)

	assert.ExpectEq(res.Result, float64(42), t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_ReturnString(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.ReturnString")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{}, t)

	assert.ExpectEq(res.Result, "Hello, World!", t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_ReturnBool(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.ReturnBool")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{}, t)

	assert.ExpectEq(res.Result, true, t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_UnregisteredClass(t *testing.T) {
	client, err := Dial("localhost:6343")
	assert.ExpectEq(err, nil, t)

	res , err := client.Call("Class.Method")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{MethodNotFound, "unregistered object: Class", nil}, t)

	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_UnknownMethod(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.UnknownMethod")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{MethodNotFound, "unknown method: UnknownMethod", nil}, t)

	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
	assert.ExpectEq(res.Id != "", true, t)
}

func TestHandler_InvalidParams(t *testing.T) {
	client, err := Dial("localhost:6342")
	assert.ExpectEq(err, nil, t)

	res, err := client.Call("Class.Method", "extra parameter")
	assert.ExpectEq(err, nil, t)
	assert.ExpectEq(res.Error, Error{InvalidParams, "given parameters do not match desired method", nil}, t)

	assert.ExpectEq(res.Result, nil, t)
	assert.ExpectEq(res.Version, JsonrpcVersion, t)
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