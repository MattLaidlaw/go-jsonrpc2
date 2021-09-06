# Go-JsonRPC2
This project implements (most of) the JSON-RPC 2.0 specification for Go applications. Go-JsonRPC2 uses the [encoding/json](https://pkg.go.dev/encoding/json) package from the Go standard for marshalling and unmarshalling json data.

## Installation

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
  
  // (optionally) start a goroutine before calling the blocking `Listen` method
  go func() {
    // block and continuously accept new Go-JsonRPC2 client connections to port 6342 (unless error returned)
    // each connection handled on its own goroutine
    err := srv.Listen("6342")
    if err != nil {
      log.Fatalln(err)
    }    
  }()

}
```
### Client
The Client object attempts to connect to a Go-JsonRPC2 server and makes remote procedure calls. The following example works assuming the above server implementation is running.
```go
package main

import (
  "github.com/MattLaidlaw/go-jsonrpc2"  // import the jsonrpc2 package
  "log"
)

func main() {

  // attempt to connect to the RPC server hosted at localhost:6342
  client, err := jsonrpc2.Dial("localhost:6342")
  if err != nil {
    log.Fatalln(err)
  }
  
  // make the RPC to Class.Method with input argument 42
  res, err := client.Call("Class.Method", 42)
  if err != nil {
    log.Fatalln(err)
  }
  
  log.Println(res.Result)

}
```
We expect a result of 42, considering Class::Method(n int) returns exactly what was input. The protocol may also return an error in the case of user error or an internal server error. This can be queried from ```res.Error```.

## Limitations
Parameters and return values of methods registered by the RPC server must be one of the following types.
* float64
* string
* bool

More work has to be completed for reflecting on the parameters and return values of registered methods. When the RPC server decodes json data sent from the client, it cannot know the exact data type a number is, so it defaults to float64. This issue may not be present for all data types, but the ones shown above are the only data types that are tested. 
