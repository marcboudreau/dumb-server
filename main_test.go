package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHeader(t *testing.T) {
	headers := &http.Header{}

	parseHeader(headers, "")
	assert.Equal(t, 0, len(*headers))

	parseHeader(headers, ":bad")
	assert.Equal(t, 0, len(*headers))

	parseHeader(headers, "a:A")
	assert.Equal(t, 1, len(*headers))
	assert.Equal(t, "A", headers.Get("a"))

	parseHeader(headers, "b:")
	assert.Equal(t, 2, len(*headers))
	assert.Equal(t, "", headers.Get("b"))
}

func TestServeHTTP(t *testing.T) {
	handler := &DumbHandler{
		headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	testcases := []struct {
		body       []byte
		statusCode int
		method     string
		target     string
	}{
		{
			body:       []byte("{}"),
			statusCode: 200,
			method:     "GET",
			target:     "/",
		},
		{
			body:       []byte("{\"key\": \"value\"}"),
			statusCode: 200,
			method:     "GET",
			target:     "/path",
		},
		{
			body:       []byte("{}"),
			statusCode: 204,
			method:     "GET",
			target:     "/",
		},
		{
			body:       []byte("An error occurred"),
			statusCode: 500,
			method:     "POST",
			target:     "/error",
		},
	}

	for _, testcase := range testcases {
		writer := httptest.NewRecorder()

		handler.body = testcase.body
		handler.statusCode = testcase.statusCode

		request := httptest.NewRequest(testcase.method, testcase.target, strings.NewReader(""))

		handler.ServeHTTP(writer, request)

		assert.Equal(t, testcase.statusCode, writer.Code)
		assert.Equal(t, testcase.body, writer.Body.Bytes())
		assert.Equal(t, "application/json", writer.HeaderMap.Get("Content-Type"))
	}
}

func TestParseResponse(t *testing.T) {
	testcases := []struct {
		response string
		headers  http.Header
		body     string
	}{
		{
			response: "header:header value\n\nbody text",
			headers: http.Header{
				"Header": []string{"header value"},
			},
			body: "body text\n",
		},
		{
			response: "\nbody but no headers",
			headers:  http.Header{},
			body:     "body but no headers\n",
		},
		{
			response: "header:no body\n\n",
			headers: http.Header{
				"Header": []string{"no body"},
			},
			body: "",
		},
		{
			response: "header:value\n\nbody\n",
			headers: http.Header{
				"Header": []string{"value"},
			},
			body: "body\n",
		},
	}

	for _, testcase := range testcases {
		buffer := &bytes.Buffer{}
		headers := &http.Header{}
		parseResponse(strings.NewReader(testcase.response), buffer, headers)

		assert.Equal(t, testcase.body, buffer.String())
		assert.Equal(t, testcase.headers, *headers)
	}
}

func TestNewDumbHandler(t *testing.T) {
	handler := NewDumbHandler(200, strings.NewReader("Header:value\n\nBody Text\n"))
	assert.Equal(t, 200, handler.statusCode)
	assert.Equal(t, []byte("Body Text\n"), handler.body)
	assert.Equal(t, 1, len(handler.headers))
	assert.Equal(t, "value", handler.headers.Get("Header"))
}
