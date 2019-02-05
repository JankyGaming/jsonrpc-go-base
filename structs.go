package main

//RPC Structs

var (
	//User Errors -32000 - -32010
	rpcErrUserNotFound           = -32000
	rpcErrEmailAlreadyRegistered = -32001
	rpcErrIncorrectPassword      = -32002
	rpcErrEmailInvalid           = -32003
	rpcErrUnknownLoginError      = -32010

	//Other Errors -32090 - -32099
	rpcErrUnhandledMongoError = -32090

	// Standard Errors
	rpcErrInvalidParams = -32602
	rpcErrInternalError = -32603
)

type exampleToken struct {
	FirstName string
	Email     string
}

type rpcRequest struct {
	Jsonrpc string                 `json:"jsonrpc"  bson:"jsonrpc"`
	Method  string                 `json:"method"  bson:"method"`
	Params  map[string]interface{} `json:"params"  bson:"params"`
	ID      int                    `json:"id"  bson:"id"`
	channel chan rpcResponse
}

type rpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"  bson:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"  bson:"result"`
	Error   *rpcError   `json:"error,omitempty"  bson:"error"`
	ID      int         `json:"id"  bson:"id"`
}

type rpcError struct {
	Code    int         `json:"code,omitempty"  bson:"code"`
	Message string      `json:"message,omitempty"  bson:"message"`
	Data    interface{} `json:"data,omitempty"  bson:"data"`
}
