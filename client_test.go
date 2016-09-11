package pikago

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	client *PikabinClient
	mux    *http.ServeMux
	server *httptest.Server
)

func setUp() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient()
	url := server.URL
	client.apiURL = url
}

func tearDown() {
	server.Close()
}

func assertEqual(t *testing.T, result interface{}, expect interface{}) {
	if result != expect {
		t.Fatalf("Expect (Value: %v) (Type: %T) - Got (Value: %v) (Type: %T)", expect, expect, result, result)
	}
}

func TestNewClient(t *testing.T) {
	c, _ := NewClient()

	assertEqual(t, c.UserAgent, ua)
	assertEqual(t, c.apiURL, defaultAPIURL)
}

func TestPaste(t *testing.T) {
	setUp()
	defer tearDown()

	content, title, syntax, expiredAt := "foo", "bar", "baz", "0"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.Method, "POST")

		payload := &payload{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &payload)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		assertEqual(t, payload.Payload.Content, content)
		assertEqual(t, payload.Payload.Title, title)
		assertEqual(t, payload.Payload.Syntax, syntax)
		assertEqual(t, payload.Payload.ExpiredAt, expiredAt)

		w.WriteHeader(http.StatusCreated)
	})

	resp, err := client.Paste(Document{
		Content:   content,
		Title:     title,
		Syntax:    syntax,
		ExpiredAt: expiredAt,
	})
	if err != nil {
		t.Fatalf("Paste(): %v", err)
	}

	assertEqual(t, http.StatusCreated, resp.StatusCode)
}

func TestAPIUrl(t *testing.T) {
	c, _ := NewClient()
	u := "http://cuonglm.xyz"
	if err := APIUrl(u)(c); err != nil {
		t.Fatalf("APIUrl() %+v", err)
	}

	assertEqual(t, c.apiURL, u)
}
