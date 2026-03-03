package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//name:api-key-20260302200644, key:418bbb73-29bd-42dd-9908-c8ee99657fb6

const (
	apiKey = "418bbb73-29bd-42dd-9908-c8ee99657fb6"
	model  = "doubao-seed-2-0-code-preview-260215" // 你的模型ID
	apiUrl = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func main() {
	reqBody := Request{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: "你好，简单介绍一下自己",
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("返回结果：")
	fmt.Println(string(body))
}
