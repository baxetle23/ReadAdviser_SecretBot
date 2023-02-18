package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	tgBotHost   = "api.clients.org"
	storagePath = "storage"
	batchSize   = 100
)

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to clients bot",
	)

	flag.Parse()

	return *token
}

type UpdateResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type payload struct {
	Offset  int    `json:"offset"`
	Limit   string `json:"limit"`
	Timeout string `json:"timeout"`
}

type messageResponse struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type configInfo struct {
	URL string `json:"URL"`
}

const (
	getUpdatesAPI  = "getUpdates"
	sendMessageAPI = "sendMessage"
)

var (
	UrlTokenBot string
)

func main() {
	fmt.Println("telegram bot Lonna start to work")
	//
	//tgClient := clients.New(tgBotHost, mustToken())
	//
	//eventsProcessor := telegram.New(tgClient, files.New(storagePath))
	//
	//consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	//
	//if err := consumer.Start(); err != nil {
	//	log.Fatal(err)
	//}

	UrlTokenBot := getURLFromConfig("./config.json")

	payload_ := payload{}

	for {
		jsonPayload, err := json.Marshal(payload_)
		if err != nil {
			log.Fatal(err)
		}

		payload := strings.NewReader(string(jsonPayload))
		req, _ := http.NewRequest(http.MethodPost, UrlTokenBot+getUpdatesAPI, payload)

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		UpdateResponse_ := &UpdateResponse{}
		jsonErr := json.Unmarshal(body, &UpdateResponse_)

		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		if UpdateResponse_.Ok == true && len(UpdateResponse_.Result) > 0 {
			fmt.Println(string(body))
			fmt.Println()
			payload_.Offset = UpdateResponse_.Result[0].UpdateID + 1
			if err := sendMessage(UpdateResponse_, UrlTokenBot); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(1 * time.Second)
		defer res.Body.Close()
	}

}

func getURLFromConfig(pathFile string) string {
	plan, err := ioutil.ReadFile(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	configInfo_ := configInfo{}
	err = json.Unmarshal(plan, &configInfo_)
	if err != nil {
		log.Fatal(err)
	}
	return configInfo_.URL
}

func sendMessage(UpdateResponse *UpdateResponse, Token string) error {
	mirrorResponse := messageResponse{
		strconv.Itoa(UpdateResponse.Result[0].Message.Chat.ID),
		UpdateResponse.Result[0].Message.Text,
	}

	response, err := json.Marshal(mirrorResponse)
	if err != nil {
		log.Fatal(err)
	}
	payload := strings.NewReader(string(response))
	req, _ := http.NewRequest(http.MethodPost, Token+sendMessageAPI, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	http.DefaultClient.Do(req)

	return err
}
