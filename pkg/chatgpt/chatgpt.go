package chatgpt

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DLzer/gojira/models"
	"github.com/ayush6624/go-chatgpt"
)

func MakeGPTRequest(clientKey string, message string) (string, error) {
	c, err := chatgpt.NewClient(clientKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	ctx := context.Background()

	res, err := c.SimpleSend(ctx, message)
	if err != nil {
		log.Println(err)
		return "", err
	}

	a, _ := json.MarshalIndent(res, "", " ")

	var response models.ChatGPTResponse
	if err := json.Unmarshal(a, &response); err != nil {
		log.Println(err)
		return "", err
	}

	return response.Choices[0].Message.Content, err
}
