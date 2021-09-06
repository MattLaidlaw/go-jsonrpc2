# Go-JsonRPC2
This project implements (most of) the JSON-RPC 2.0 specification for Go applications. The exact format of request and response objects adheres to the specification defined [here](https://www.jsonrpc.org/specification). Go-JsonRPC2 uses the [encoding/json](https://pkg.go.dev/encoding/json) package from the Go standard for marshalling and unmarshalling json data.

## Installation
```
go get github.com/MattLaidlaw/go-jsonrpc2@60e3cb0bde2df6880bb1748caf67c744b5116c3d
```

## Usage
### Server
The Server object registers methods for use with the RPC protocol and supports concurrent client connections.
```go
package main

import (
  "github.com/MattLaidlaw/go-jsonrpc2"  // import the jsonrpc2 package
  "log"
)

// define a class that implements some set of methods
type Class struct {}
func (c *Class) Method(n float64) float64 {
  return n
}

func main() {

  // obtain a handle to a new server instance
  srv := jsonrpc2.NewServer()
  
  // tell the srv instance to accept methods implemented by the Class object
  srv.Register(Class{})
  
  // block and continuously accept new Go-JsonRPC2 client connections to port 6342 (unless error returned)
  // each connection handled on its own goroutine
  err := srv.Listen("6342")
  if err != nil {
    log.Fatalln(err)
  }

}
```
### Client
The Client object connects to a Go-JsonRPC2 server and makes remote procedure calls. The following example works assuming the above server implementation is running on the same host.
```go
package main

import (
  "fmt"
  "github.com/MattLaidlaw/go-jsonrpc2"  // import the jsonrpc2 package
  "log"
)

func main() {

  // attempt to connect to the RPC server hosted at localhost:6342
  client, err := jsonrpc2.Dial("localhost:6342")
  if err != nil {
    log.Fatalln(err)
  }
  
  // make a remote procedure call to Class.Method with input argument 42
  res, err := client.Call("Class.Method", 42)
  if err != nil {
    log.Fatalln(err)
  }
  
  fmt.Println(res.Result)

}
```
We expect a result of 42, considering Class::Method(n int) returns the input parameter as-is. The protocol may also return an error which can be queried from ```res.Error```.

## Limitations
Parameters and return values of methods registered by the RPC server must be one of the following types.
* float64
* string
* bool

In the current state of this project, all numbers encoded into json are decoded into float64 because the exact type information is lost in the encoding/decoding process. This issue may not be present for all kinds of data, but the data types shown above are the only ones tested for correctness.
