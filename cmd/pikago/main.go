package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Gnouc/pikago"
)

var (
	titleFlag   string
	syntaxFlag  string
	expiredFlag string
	client      *pikago.PikabinClient
	err         error
)

func init() {
	flag.StringVar(&titleFlag, "title", "", "Paste title")
	flag.StringVar(&syntaxFlag, "syntax", "plain", "Coloring syntax, see more: http://goo.gl/nLFqyB")
	flag.StringVar(&expiredFlag, "expired-at", "0", "Set expiration, in minute, -1 means no expiration. Default is 0, after reading")

	apiURL := os.Getenv("PIKABIN_URL")
	if apiURL == "" {
		client, _ = pikago.NewClient()
	} else {
		client, err = pikago.NewClient(pikago.APIUrl(apiURL))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.Parse()

	reader := os.Stdin
	if len(flag.Args()) > 1 {
		var err error
		reader, err = os.Open(flag.Args()[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	d := pikago.Document{
		Content:   string(content),
		Syntax:    syntaxFlag,
		ExpiredAt: expiredFlag,
		Title:     titleFlag,
	}

	resp, err := client.Paste(d)
	if err != nil {
		log.Fatal(err)
	}

	extractResponse(resp)
}

func extractResponse(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		defer func() {
			resp.Body.Close()
		}()
	}

	if err != nil {
		log.Fatal(err)
	}
	type response struct {
		Message []string
		URI     string
	}
	r := &response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatal(err)
	}

	if len(r.Message) > 0 {
		log.Fatal(strings.Join(r.Message, "\n"))
	}

	fmt.Println(r.URI)
}
