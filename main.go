package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultResponse = `Content-Type: text/plain

dumb-server default response`

var port = flag.Int("port", 7979, "The TCP port the server listens on")
var statusCode = flag.Int("sc", 200, "The HTTP Status Code returned with every request")
var responseFile = flag.String("resp", "", "A file containing the response to return with every request")

func main() {
	flag.Parse()

	var response io.Reader
	if *responseFile != "" {
		file, err := os.Open(*responseFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		response = io.Reader(file)
	} else {
		response = strings.NewReader(defaultResponse)
	}

	s := http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        NewDumbHandler(*statusCode, response),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

// DumbHandler is a Handler that doesn't care about the request parameters and always responds with the
// same response.
type DumbHandler struct {
	statusCode int
	headers    http.Header
	body       []byte
}

// NewDumbHandler creates a new DumbHandler instance by parsing the provided
// response file into headers and a response body.
func NewDumbHandler(statusCode int, response io.Reader) *DumbHandler {
	handler := &DumbHandler{
		statusCode: statusCode,
		headers:    http.Header{},
	}

	bodyBuf := &bytes.Buffer{}
	parseResponse(response, bodyBuf, &handler.headers)

	handler.body = bodyBuf.Bytes()

	return handler
}

func parseResponse(reader io.Reader, bodyBuf *bytes.Buffer, headers *http.Header) {
	inbody := false
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			inbody = true
			continue
		}

		if inbody {
			bodyBuf.WriteString(line)
			bodyBuf.WriteByte('\n')
		} else {
			parseHeader(headers, line)
		}
	}
}

func (h *DumbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for hn, hv := range h.headers {
		for _, v := range hv {
			w.Header().Add(hn, v)
		}
	}

	w.WriteHeader(h.statusCode)
	if 0 < len(h.body) {
		w.Write(h.body)
	}
}

func parseHeader(headers *http.Header, headerLine string) {
	// For a headerLine to be valid there needs to be a colon present
	// so i != -1, and the colon can't be in the first position, so
	// i != 0 either.
	if i := strings.Index(headerLine, ":"); i > 0 {
		headers.Add(headerLine[0:i], headerLine[i+1:])
	}
}
