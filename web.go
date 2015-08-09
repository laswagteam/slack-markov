package main

// Slack outgoing webhooks are handled here. Requests come in and are run through
// the markov chain to generate a response, which is sent back to Slack.
//
// Create an outgoing webhook in your Slack here:
// https://my.slack.com/services/new/outgoing-webhook

import (
	"encoding/json"
	"log"
	"strconv"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type WebhookResponse struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		incomingText := r.PostFormValue("text")
		if incomingText != "" && r.PostFormValue("user_id") != "" && r.PostFormValue("user_id") != "USLACKBOT"{
			text := parseText(incomingText)
			//log.Printf("Handling incoming request: %s", text)

			if text != "" {
				markovChain.Write(text)
			}

			go func() {
				markovChain.Save(stateFile)
			}()

			if rand.Intn(100) < responseChance || strings.Contains(strings.ToLower(incomingText), strings.ToLower(botUsername)) {
				var response WebhookResponse
				response.Username = botUsername
				if strings.Contains(incomingText, "TG") && strings.HasPrefix(strings.ToLower(incomingText), strings.ToLower(botUsername)) {
					if strings.Contains(incomingText, "poil") {
						responseChance -= 1
					} else{
						responseChance -= 5
					}
					if responseChance < 0 {
						responseChance = 0
					}
					response.Text = "Okay :( je suis à "+strconv.Itoa(responseChance)+"%"
				} else if strings.Contains(incomingText, "BS") && strings.HasPrefix(strings.ToLower(incomingText), strings.ToLower(botUsername)) {
					if strings.Contains(incomingText, "poil") {
						responseChance += 1
					} else{
						responseChance += 5
					}
					if responseChance > 100 {
						responseChance = 100
					}
					response.Text = "Okay :D je suis à "+strconv.Itoa(responseChance)+"%"
				} else if strings.Contains(incomingText, "moral") && strings.HasPrefix(strings.ToLower(incomingText), strings.ToLower(botUsername)) {
					response.Text = "Environ "+strconv.Itoa(responseChance)+"% mon capitaine !"
				} else {
					response.Text = markovChain.Generate(numWords)
				}
				//log.Printf("Sending response: %s", response.Text)

				b, err := json.Marshal(response)
				if err != nil {
					log.Fatal(err)
				}

				time.Sleep(5 * time.Second)
				w.Write(b)
			}
		}
	})
}

func StartServer(addr string) {
	log.Printf("Starting HTTP server on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
