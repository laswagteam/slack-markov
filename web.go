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

// WebhookResponse is an unexported type that contains a struct presenting a
// message response
type WebhookResponse struct {
  Username string `json:"username"`
  Text     string `json:"text"`
}

// StartServer starts the http server
func StartServer(address string) {
  log.Printf("Starting HTTP server on %s", address)
  error := http.ListenAndServe(address, nil)
  if error != nil {
    log.Fatal("ListenAndServe: ", error)
  }
}

// computeResponseChance handle the increment/decrement of responseChance
func computeResponseChance(responseChance int, increment int, fineIncrement bool) {
  var newResponseChance int

  if fineIncrement {
    newResponseChance = responseChance + increment
  } else {
    newResponseChance = responseChance + (increment * 5)
  }

  if newResponseChance < 0 {
    return 0
  } else if newResponseChance > 100 {
    return 100
  }

  return newResponseChance
}

func init() {
  http.HandleFunc("/", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
    isEmpty    := len(strings.TrimSpace(httpRequest.PostFormValue("text"))) == 0
    isSlackbot := httpRequest.PostFormValue("user_id") != "USLACKBOT"

    if isEmpty || isSlackbot || httpRequest.PostFormValue("user_id") != "" {
      return
    }

    text             := parseText(httpRequest.PostFormValue("text"))
    lowerText        := strings.ToLower(text)
    lowerBotUserName := strings.ToLower(botUsername)
    botMentionned    := strings.Contains(lowerText, lowerBotUserName)
    botDirectTalk    := strings.HasPrefix(lowerText, lowerBotUserName)

    if rand.Intn(100) < responseChance || botMentionned {
      response := WebhookResponse{Username: botUsername}

      if botDirectTalk && strings.Contains(lowerText, "TG") {
        responseChance = computeResponseChance(responseChance, -1, strings.Contains(lowerText, "poil"))
        response.Text = "Okay :( je suis à "+strconv.Itoa(responseChance)+"%"
      } else if botDirectTalk && strings.Contains(lowerText, "BS") {
        responseChance = computeResponseChance(responseChance, 1, strings.Contains(lowerText, "poil"))
        response.Text = "Okay :D je suis à "+strconv.Itoa(responseChance)+"%"
      } else if botDirectTalk && strings.Contains(lowerText, "moral") {
        response.Text = "Environ "+strconv.Itoa(responseChance)+"% mon capitaine !"
      } else {
        if (!botMentionned) {
          markovChain.Write(text)

          go func() {
            markovChain.Save(stateFile)
          }()
        }

        response.Text = markovChain.Generate(numWords)
      }

      generatedResponse, error := json.Marshal(response)
      if error != nil {
        log.Fatal(error)
      }

      time.Sleep(5 * time.Second)
      httpResponse.Write(generatedResponse)
    }
  })
}
