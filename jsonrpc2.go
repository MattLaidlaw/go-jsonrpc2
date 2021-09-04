// Package jsonrpc2 implements the JSON-RPC 2.0 specification
package go_jsonrpc2

const JsonrpcVersion string = "2.0"

type ErrCode int
const (
	ParseError     ErrCode = -32700
	InvalidRequest ErrCode = -32600
	MethodNotFound ErrCode = -32601
	InvalidParams  ErrCode = -32602
	InternalError  ErrCode = -32603
)

type Request struct {
	Version	string			`json:"jsonrpc"`
	Method 	string 			`json:"method"`
	Params 	[]interface{}	`json:"params,omitempty"`
	Id 		string			`json:"id,omitempty"`
}

type Response struct {
	Version string		`json:"jsonrpc"`
	Result interface{} `json:"result,omitempty"`
	Error  Error       `json:"error,omitempty"`
	Id     string      `json:"id,omitempty"`
}

type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
	Data	interface{}	`json:"data,omitempty"`
}
