package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Callback struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender    User    `json:"sender,omitempty"`
	Recipient User    `json:"recipient,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Message struct {
	MID        string `json:"mid,omitempty"`
	Text       string `json:"text,omitempty"`
	QuickReply *struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
	Attachment  *Attachment   `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Response struct {
	Recipient User    `json:"recipient,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type Payload struct {
	URL string `json:"url,omitempty"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeEndpoint)
	r.HandleFunc("/webhook", VerificationEndPoint).Methods("GET")
	r.HandleFunc("/webhook", MessagesEndPoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

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

	body := fmt.Sprintf(`{
		"recipient": {
		  "id": "%s"
		},
		"message": {
		  "text": "hello, world!"
		}
	  }`, event.Sender.ID)

	url := fmt.Sprintf("https://graph.facebook.com/v2.6/me/messages?access_token=%s", os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	fmt.Println("The response is", resp.Body)
	fmt.Println("The response code is", resp.StatusCode)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
