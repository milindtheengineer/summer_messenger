package witai

import "time"

type Response struct {
	Text     string `json:"_text"`
	Entities struct {
		Question []struct {
			Confidence float64 `json:"confidence"`
			Value      string  `json:"value"`
			Type       string  `json:"type"`
		} `json:"question"`
		Datetime []struct {
			Confidence float64 `json:"confidence"`
			Values     []struct {
				Value time.Time `json:"value"`
				Grain string    `json:"grain"`
				Type  string    `json:"type"`
			} `json:"values"`
			Value time.Time `json:"value"`
			Grain string    `json:"grain"`
			Type  string    `json:"type"`
		} `json:"datetime"`
	} `json:"entities"`
	MsgID string `json:"msg_id"`
}
