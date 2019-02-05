package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/JankyGaming/easygo"
)

var publicFM = map[string]interface{}{
	"testFuncPublic": testFuncPublic,
}

var publicFMInfo = map[string]interface{}{
	"testFuncPrivate": map[string]interface{}{
		"desc":           "a test function so you can run this base and see it work, before putting in time",
		"requiredParams": map[string]interface{}{},
		"optionalParams": map[string]interface{}{
			"firtName": "string",
		},
		"result": map[string]interface{}{
			"respone": "Hello, user",
		},
	},
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//GET requests responds with map defined aboce publicFMInfo, you write that map to explain how your functions work, and whats available.
		easygo.Respond(w, r, 200, snsFMInfo, map[string]string{"Content-Type": "application/json"})
		return
	} else if r.Method != "POST" {
		easygo.Respond(w, r, 405, "only GET and POST are supported on this rpc endpoint", nil)
		return
	}

	//Decode request into array
	requestArr := []*rpcRequest{}
	rawBody := r.Body
	bytBody, err := ioutil.ReadAll(rawBody)
	if err != nil {
		easygo.Respond(w, r, 400, easygo.ResponseObject{Error: true, Message: "bad body"}, nil)
		return
	}
	if len(bytBody) == 0 {
		easygo.Respond(w, r, 400, easygo.ResponseObject{Error: true, Message: "empty body"}, nil)
		return
	}
	if bytBody[0] == []byte("[")[0] {
		err = json.Unmarshal(bytBody, &requestArr)
		if err != nil {
			easygo.Respond(w, r, 200, rpcResponse{Jsonrpc: "2.0", Error: &rpcError{Code: -32700, Message: "Parse error"}}, nil)
			return
		}
	} else {
		newSingleRequest := rpcRequest{}
		err = json.Unmarshal(bytBody, &newSingleRequest)
		if err != nil {
			easygo.Respond(w, r, 200, rpcResponse{Jsonrpc: "2.0", Error: &rpcError{Code: -32700, Message: "Parse error"}}, nil)
			return
		}
		requestArr = append(requestArr, &newSingleRequest)
	}

	//Loop array and run requests for each request
	for _, request := range requestArr {
		request.channel = make(chan rpcResponse, 1)
		if funcName, ok := publicFM[request.Method]; ok {
			go funcName.(func(*rpcRequest, *http.Request))(request, r)
		} else {
			newResponse := rpcResponse{
				Jsonrpc: request.Jsonrpc,
				Error: &rpcError{
					Code:    -32601,
					Message: "Method not found",
				},
				ID: request.ID,
			}
			request.channel <- newResponse
		}
	}
	responseArr := []rpcResponse{}

	for _, request := range requestArr {
		object := <-request.channel
		responseArr = append(responseArr, object)
	}

	if len(responseArr) == 1 {
		easygo.Respond(w, r, 200, responseArr[0], nil)
		return
	}

	easygo.Respond(w, r, 200, responseArr, nil)
	return
}

func testFuncPublic(rpc *rpcRequest, r *http.Request) {

	response := "Hello user"

	if name, ok := rpc.Params["firstName"].(string); ok {
		response = "Hello " + name
	}

	rpc.channel <- rpcResponse{
		Jsonrpc: rpc.Jsonrpc,
		Result:  response,
		ID:      rpc.ID,
	}
	return
}
