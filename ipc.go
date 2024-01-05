package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

type Request interface{}

type Response interface{}

const (
	REQUEST_TYPE_CREATE_SESSION             = "create_session"
	REQUEST_TYPE_POST_AUTH_MESSAGE_RESPONSE = "post_auth_message_response"
	REQUEST_TYPE_START_SESSION              = "start_session"
	REQUEST_TYPE_CANCEL_SESSION             = "cancel_session"
	RESPONSE_TYPE_SUCCESS                   = "success"
	RESPONSE_TYPE_ERROR                     = "error"
	RESPONSE_TYPE_AUTH_MESSAGE              = "auth_message"
	AUTH_MESSAGE_TYPE_VISIBLE               = "visible"
	AUTH_MESSAGE_TYPE_SECRET                = "secret"
	AUTH_MESSAGE_TYPE_INFO                  = "info"
	AUTH_MESSAGE_TYPE_ERROR                 = "error"
	ERROR_TYPE_AUTH_ERROR                   = "auth_error"
	ERROR_TYPE_ERROR                        = "error"
)

type CreateSessionRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

type PostAuthMessageResponseRequest struct {
	Type     string `json:"type"`
	Response string `json:"response,omitempty"`
}

type StartSessionRequest struct {
	Type string   `json:"type"`
	Cmd  []string `json:"cmd"`
	Env  []string `json:"env"`
}

type CancelSessionRequest struct {
	Type string `json:"type"`
}

type SuccessResponse struct {
	Type string `json:"type"`
}

type ErrorResponse struct {
	Type        string `json:"type"`
	ErrorType   string `json:"error_type"`
	Description string `json:"description"`
}

type AuthMessageResponse struct {
	Type            string `json:"type"`
	AuthMessageType string `json:"auth_message_type"`
	AuthMessage     string `json:"auth_message"`
}

func NewCreateSessionRequest(username string) *CreateSessionRequest {
	return &CreateSessionRequest{Type: REQUEST_TYPE_CREATE_SESSION, Username: username}
}

func NewPostAuthMessageResponseRequest(response string) *PostAuthMessageResponseRequest {
	return &PostAuthMessageResponseRequest{Type: REQUEST_TYPE_POST_AUTH_MESSAGE_RESPONSE, Response: response}
}

func NewStartSessionRequest(cmd []string, env []string) *StartSessionRequest {
	return &StartSessionRequest{Type: REQUEST_TYPE_START_SESSION, Cmd: cmd, Env: env}
}

func NewCancelSessionRequest() *CancelSessionRequest {
	return &CancelSessionRequest{Type: REQUEST_TYPE_CANCEL_SESSION}
}

func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{Type: RESPONSE_TYPE_SUCCESS}
}

func NewAuthMessageResponse(authMessageType string, authMessage string) *AuthMessageResponse {
	return &AuthMessageResponse{
		Type:            RESPONSE_TYPE_AUTH_MESSAGE,
		AuthMessageType: authMessageType,
		AuthMessage:     authMessage,
	}
}

func NewErrorResponse(errorType string, description string) *ErrorResponse {
	return &ErrorResponse{
		Type:        RESPONSE_TYPE_ERROR,
		ErrorType:   errorType,
		Description: description,
	}
}

func WriteTo(conn net.Conn, payload interface{}) {
	v, err := json.Marshal(payload)
	if err != nil {
		log.Panicln("Can't convert to JSON", err)
	}
	length := make([]byte, 4)
	// original greetd implementation used hardware dependent endian order, so this might fuck something up later
	binary.LittleEndian.PutUint32(length, uint32(len(v)))

	_, err = conn.Write(length)

	if err != nil {
		log.Panicln("Error writing to socket", err)
	}

	_, err = conn.Write(v)
	if err != nil {
		log.Panicln("Error writing to socket", err)
	}
}

func ReadFrom(conn net.Conn) interface{} {
	length := make([]byte, 4)
	nbytes, err := conn.Read(length)
	if err != nil {
		log.Panicln("Error reading from socket", err)
	}
	if nbytes != 4 {
		log.Panicf("Byte read is %d but needed 4", nbytes)
	}

	// The case as in Write
	sz := binary.LittleEndian.Uint32(length)

	data := make([]byte, sz)

	nbytes, err = conn.Read(data)
	if err != nil {
		log.Panicln("Error reading from socket", nbytes, err)
	}

	return parsePayload(data)
}

func parsePayload(data []byte) interface{} {
	var res map[string]interface{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Panicln("Read data is not JSON", err)
	}

	var payload interface{}
	switch res["type"] {
	case REQUEST_TYPE_CREATE_SESSION:
		payload = &CreateSessionRequest{}
	case REQUEST_TYPE_POST_AUTH_MESSAGE_RESPONSE:
		payload = &PostAuthMessageResponseRequest{}
	case REQUEST_TYPE_START_SESSION:
		payload = &StartSessionRequest{}
	case REQUEST_TYPE_CANCEL_SESSION:
		payload = &CancelSessionRequest{}
	case RESPONSE_TYPE_SUCCESS:
		payload = &SuccessResponse{}
	case RESPONSE_TYPE_ERROR:
		payload = &ErrorResponse{}
	case RESPONSE_TYPE_AUTH_MESSAGE:
		payload = &AuthMessageResponse{}
	default:
		log.Panicln("Unrecognized type", res["type"])
	}

	json.Unmarshal(data, &payload)
	return payload
}
