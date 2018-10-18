package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"summer/witai"

	"github.com/joho/godotenv"
)

func HomeEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Milind :)")
}

func VerificationEndPoint(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	challenge := r.URL.Query().Get("hub.challenge")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	if mode != "" && token == os.Getenv("PAGE_VERIFY_TOKEN") {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}

func MessagesEndPoint(w http.ResponseWriter, r *http.Request) {
	var callback Callback
	json.NewDecoder(r.Body).Decode(&callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				ProcessMessage(event)
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}
}

func ProcessMessage(event Messaging) {
	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := &http.Client{}

	fmt.Println(event.Sender.ID)
	text := event.Message.Text
	message := witai.ExtractMessage(text)
	fmt.Println("The message is ", message)
	body := fmt.Sprintf(`{
		"recipient": {
		  "id": "%s"
		},
		"message": {
		  "text": "%s"
		}
	  }`, event.Sender.ID, message)

	url := fmt.Sprintf("https://graph.facebook.com/v2.6/me/messages?access_token=%s", os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
