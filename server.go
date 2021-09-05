package jsonrpc2

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
)

type Server struct {
	registered map[string]reflect.Type
}

func NewServer() *Server {
	return &Server{
		registered: make(map[string]reflect.Type),
	}
}

func (s *Server) Listen(port string) error {
	listener, err := net.Listen("tcp", "localhost:" + port)
	if err != nil {
		return err
	}
	log.Println("== listening on port", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("== accepted connection from", conn.RemoteAddr().Network())
		h := NewHandler(conn, s.registered)
		go h.Handle()
	}
}

// The Register method tells the RPC Server to accept method signatures matching those implemented by the registered
// class.
func (s *Server) Register(class interface{}) {
	t := reflect.TypeOf(class)
	s.registered[t.Name()] = t
}

// The Handler object drives a single JSON-RPC 2.0 connection.
type Handler struct {
	conn       net.Conn
	registered map[string]reflect.Type
}

// The NewHandler method constructs a Handler on top of a given connection with a list of registered methods.
func NewHandler(conn net.Conn, registered map[string]reflect.Type) *Handler {
	return &Handler{
		conn:       conn,
		registered: registered,
	}
}

// The Handle method decodes incoming requests, executes them, and sends an encoded response. This method blocks until
// EOF received from client.
func (h *Handler) Handle() {
	decoder := json.NewDecoder(h.conn)
	encoder := json.NewEncoder(h.conn)

	for {
		// decode incoming JSON request
		var req Request
		err := decoder.Decode(&req)
		if err == io.EOF {
			break
		} else if err != nil {
			resErr := Error{
				Code:    ParseError,
				Message: "unable to parse JSON request body",
			}
			resData := Response{
				Version: JsonrpcVersion,
				Error:   resErr,
			}
			_ = encoder.Encode(resData)
			if err != nil {
				log.Println(err)
			}
		}

		res := h.execute(&req)

		// encode outgoing JSON response
		err = encoder.Encode(res)
		if err != nil {
			log.Println(err)
		}
	}
}

// The execute method attempts the given function call with the supplied arguments. The call string should be of
// the form 'class.method', where the class has been registered, and the method is a valid member of that class. The
// arguments should match those of the registered method. If any of these statements are not true, the function call is
// not made and an error is returned.
func (h *Handler) execute(req *Request) Response {
	// get, and validate, call string of form class.method
	fields := strings.Split(req.Method, ".")
	if len(fields) < 2 {
		resErr := Error{
			Code:    MethodNotFound,
			Message: "call string must be of form class.method: " + req.Method,
		}
		return Response{
			Version: JsonrpcVersion,
			Error:   resErr,
			Id:      req.Id,
		}
	}

	// check if the supplied object is registered
	object, ok := h.registered[fields[0]]
	if !ok {
		resErr := Error{
			Code:    MethodNotFound,
			Message: "unregistered object: " + fields[0],
		}
		return Response{
			Version: JsonrpcVersion,
			Error:   resErr,
			Id:      req.Id,
		}
	}

	// check if the supplied method is a valid method of the given class
	obj := reflect.New(object).Interface()
	method := reflect.ValueOf(obj).MethodByName(fields[1])
	if !method.IsValid() {
		resErr := Error{
			Code:    MethodNotFound,
			Message: "unknown method: " + fields[1],
		}
		return Response{
			Version: JsonrpcVersion,
			Error:   resErr,
			Id:      req.Id,
		}
	}

	// check if the supplied arguments correlate with those of the matched method
	if method.Type().NumIn() != len(req.Params) {
		resErr := Error{
			Code:    InvalidParams,
			Message: "given parameters do not match desired method",
		}
		return Response{
			Version: JsonrpcVersion,
			Error:   resErr,
			Id:      req.Id,
		}
	}

	// convert variadic args into an array of values
	params := make([]reflect.Value, len(req.Params))
	for i := range params {
		params[i] = reflect.ValueOf(req.Params[i])
	}

	resultArr := method.Call(params)
	if len(resultArr) != 0 {
		return Response {
			Version: JsonrpcVersion,
			Result:  resultArr[0].Interface(),
			Id:      req.Id,
		}
	} else {
		return Response {
			Version: JsonrpcVersion,
			Id:      req.Id,
		}
	}
}
