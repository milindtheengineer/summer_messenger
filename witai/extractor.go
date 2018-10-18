package witai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"summer/football"

	"github.com/joho/godotenv"
)

func make_get_call(message string) []byte {
	err := godotenv.Load("../secrets.env")
	client := &http.Client{}
	url := fmt.Sprintf("https://api.wit.ai/message")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("q", message)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	req.Header.Add("Authorization", os.Getenv("WIT_AI_ACCESS_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		panic(err)
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody
}

func ExtractMessage(message string) string {
	byt := make_get_call(message)
	responseMessage := Response{}
	if err := json.Unmarshal(byt, &responseMessage); err != nil {
		panic(err)
	}
	fmt.Println(responseMessage)
	if responseMessage.Entities.Question[0].Value == "football matches" && responseMessage.Entities.Question != nil {
		return football.SendMatches(responseMessage.Entities.Datetime[0].Value.Format("2006-01-02"))
	} else {
		return "No action for this message yet"
	}
}
