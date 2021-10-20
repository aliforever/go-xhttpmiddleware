package tests

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aliforever/go-xhttpmiddleware"
)

var message = []byte("Hello World!")

var server *http.Server

var methodChan = make(chan string, 1)

type HelloWorld struct{}

func (hw HelloWorld) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(message)
	methodChan <- r.Method
}

func TestServeHTTP(t *testing.T) {
	server = &http.Server{
		Addr:         ":8080",
		Handler:      xhttpmiddleware.NewXHTTPMethodOverrideHandler(HelloWorld{}, nil),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			t.Fatal("Test Failed.", err)
			return
		}
	}()

	var (
		expectedMethod string = "PATCH"
		err            error
	)

	err = testServer(expectedMethod)
	if err != nil {
		t.Fatal("Test Failed.", err)
		return
	}

	returnedMethod := <-methodChan

	if returnedMethod != "PATCH" {
		t.Fatalf("Test Failed.\nExpected Method: %s\nActual Value: %s", expectedMethod, returnedMethod)
		return
	}

	server.Close()
}

func testServer(expectedMethod string) (err error) {
	time.Sleep(time.Second * 2)

	var req *http.Request

	v := url.Values{}
	v.Set("phone_number", "+19021235263")

	req, err = http.NewRequest("POST", "http://localhost:8080/", strings.NewReader(v.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("X-HTTP-Method-Override", expectedMethod)

	_, err = http.DefaultClient.Do(req)

	return
}
