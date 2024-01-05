package main

import (
	"log"
	"net"
)

type Greeter struct {
	socket            net.Conn
	cmd               string
	done              bool
	nextRequest       Request
	OnNotice          func(notice string)
	OnRequestInput    func(label string, secret bool)
	OnDone            func()
	OnError           func()
	OnResetConnection func() net.Conn
}

func NewGreeter(socket net.Conn, cmd string) *Greeter {
	return &Greeter{
		socket:            socket,
		cmd:               cmd,
		OnNotice:          func(notice string) {},
		OnRequestInput:    func(label string, secret bool) {},
		OnDone:            func() {},
		OnError:           func() {},
		OnResetConnection: func() net.Conn { return socket },
	}
}

func (greeter *Greeter) HandleResponse(response Response) {
	switch response.(type) {
	case *AuthMessageResponse:
		authMessageResponse := response.(*AuthMessageResponse)
		switch authMessageResponse.AuthMessageType {
		case AUTH_MESSAGE_TYPE_VISIBLE:
			greeter.EmitRequestInput(authMessageResponse.AuthMessage, false)
		case AUTH_MESSAGE_TYPE_SECRET:
			greeter.EmitRequestInput(authMessageResponse.AuthMessage, true)
		case AUTH_MESSAGE_TYPE_INFO:
			greeter.nextRequest = NewPostAuthMessageResponseRequest("")
			greeter.EmitNotice("info:" + authMessageResponse.AuthMessage)
		case AUTH_MESSAGE_TYPE_ERROR:
			greeter.nextRequest = NewPostAuthMessageResponseRequest("")
			greeter.EmitNotice("error:" + authMessageResponse.AuthMessage)
		}

	case *SuccessResponse:
		if greeter.done {
			greeter.EmitNotice("Login Success")
			greeter.EmitDone()
			return
		}
		greeter.done = true
		greeter.nextRequest = NewStartSessionRequest([]string{greeter.cmd}, []string{})

	case *ErrorResponse:
		errorResponse := response.(*ErrorResponse)
		WriteTo(greeter.socket, NewCancelSessionRequest())
		if errorResponse.ErrorType == ERROR_TYPE_AUTH_ERROR {
			greeter.EmitNotice("Login Incorrect")
		} else {
			greeter.EmitNotice("Login error " + errorResponse.Description)
		}
		greeter.EmitError()
	}
}

func (greeter *Greeter) ExchangeRequest() {
	for greeter.nextRequest != nil {
		log.Printf("Request %+v\n", greeter.nextRequest)
		WriteTo(greeter.socket, greeter.nextRequest)
		greeter.nextRequest = nil
		response := ReadFrom(greeter.socket)
		log.Printf("Response %+v\n", response)
		greeter.HandleResponse(response)
	}
}

func (greeter *Greeter) HandleInput(input string) {
	greeter.nextRequest = NewPostAuthMessageResponseRequest(input)
	greeter.ExchangeRequest()
}

func (greeter *Greeter) CreateSession(input string) {
	greeter.nextRequest = NewCreateSessionRequest(input)
	greeter.ExchangeRequest()
}

func (greeter *Greeter) EmitRequestInput(label string, secret bool) {
	greeter.OnRequestInput(label, secret)
}

func (greeter *Greeter) EmitNotice(notice string) {
	greeter.OnNotice(notice)
}

func (greeter *Greeter) EmitError() {
	greeter.socket.Close()
	greeter.socket = greeter.OnResetConnection()
	greeter.OnError()
}

func (greeter *Greeter) EmitDone() {
	greeter.socket.Close()
	greeter.OnDone()
}
