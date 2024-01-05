package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestRequestJson(t *testing.T) {
	testCases := []struct {
		name     string
		request  interface{}
		expected string
	}{
		{"create session", NewCreateSessionRequest("testuser"), `{"type": "create_session", "username": "testuser"}`},
		{"Authenticate", NewPostAuthMessageResponseRequest("password"), `{"type": "post_auth_message_response", "response": "password"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			j, err := json.Marshal(tc.request)
			if err != nil {
				t.Error(err)
			}
			require.JSONEq(t, tc.expected, string(j))
		})
	}
}

func TestResponse(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		expected interface{}
	}{
		{"success", `{"type": "success"}`, &SuccessResponse{}},
		{"error", `{"type": "error"}`, &ErrorResponse{}},
		{"auth", `{"type": "auth_message", "auth_message_type": "secret"}`, &AuthMessageResponse{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			res := parsePayload(data)
			require.IsType(t, tc.expected, res)
		})
	}
}

func TestComm(t *testing.T) {
	testCases := []struct {
		data interface{}
	}{
		{NewCreateSessionRequest("username")},
		{NewErrorResponse("test", "test")},
		{NewSuccessResponse()},
		{NewPostAuthMessageResponseRequest("password")},
	}

	server, client := net.Pipe()
	t.Cleanup(func() {
		server.Close()
		client.Close()
	})
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc.data), func(t *testing.T) {

			go WriteTo(client, tc.data)
			data := ReadFrom(server)
			require.Equal(t, data, tc.data)
		})
	}
}
