package main

import (
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
)

func TestSuccessLogin(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	server, client := net.Pipe()
	greeter := NewGreeter(client, "/usr/bin/sh")
	go func() {
		defer wg.Done()
		greeter.CreateSession("username")
		greeter.HandleInput("password")
	}()
	for _, expected := range []struct {
		Request
		Response
	}{{NewCreateSessionRequest("username"), NewAuthMessageResponse(AUTH_MESSAGE_TYPE_SECRET, "password")},
		{NewPostAuthMessageResponseRequest("password"), NewSuccessResponse()},
		{NewStartSessionRequest([]string{"/usr/bin/sh"}, []string{}), NewSuccessResponse()}} {
		request := ReadFrom(server)
		require.Equal(t, expected.Request, request)
		WriteTo(server, expected.Response)

	}
	wg.Wait()
}

func TestFailLogin(t *testing.T) {
	server, client := net.Pipe()

	greeter := NewGreeter(client, "/usr/bin/sh")
	greeter.OnResetConnection = func() net.Conn {
		server.Close()
		server, client = net.Pipe()
		return client
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		greeter.CreateSession("username")
		greeter.HandleInput("password")
		// The reason is from the below comment
		wg.Done()
		greeter.CreateSession("username")
		greeter.HandleInput("password")
		wg.Done()
	}()

	for _, expected := range []struct {
		Request
		Response
	}{{NewCreateSessionRequest("username"), NewAuthMessageResponse(AUTH_MESSAGE_TYPE_SECRET, "password")},
		{NewPostAuthMessageResponseRequest("password"), NewErrorResponse(ERROR_TYPE_AUTH_ERROR, "")},
		{NewCancelSessionRequest(), nil}} {
		request := ReadFrom(server)
		require.Equal(t, expected.Request, request)
		if expected.Response != nil {
			WriteTo(server, expected.Response)
		}
	}
	// If we don't separate the 2 session, we will cause race condition when resetting the server/client pipe
	// causing read from closed yet to be reseted socket
	// So either we pause here wait for the threads, or we change from pipe to mock to a full server / client
	// implemetaion just for testing. Pipe for the time being
	wg.Wait()

	wg.Add(1)
	for _, expected := range []struct {
		Request
		Response
	}{
		{NewCreateSessionRequest("username"), NewAuthMessageResponse(AUTH_MESSAGE_TYPE_SECRET, "password")},
		{NewPostAuthMessageResponseRequest("password"), NewSuccessResponse()},
		{NewStartSessionRequest([]string{"/usr/bin/sh"}, []string{}), NewSuccessResponse()}} {
		request := ReadFrom(server)
		require.Equal(t, expected.Request, request)
		if expected.Response != nil {
			WriteTo(server, expected.Response)
		}
	}

	wg.Wait()
}
