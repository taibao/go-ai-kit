package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// AI SDK 配置
type AISDKConfig struct {
	APIKey  string
	ModelID string
	APIUrl  string
}

// 消息结构体（如果 main.go 中没有定义）
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 流式请求结构体
type StreamRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// 流式响应相关结构体
type Delta struct {
	Content string `json:"content"`
}

type Choice struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type StreamResponse struct {
	Choices []Choice `json:"choices"`
}

// 普通请求结构体
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// 新建 SDK 实例
func NewAISDK(apiKey, modelID string) *AISDKConfig {
	return &AISDKConfig{
		APIKey:  apiKey,
		ModelID: modelID,
		APIUrl:  "https://ark.cn-beijing.volces.com/api/v3/chat/completions",
	}
}

// 流式聊天
func (s *AISDKConfig) StreamChat(ctx context.Context, question string) error {
	// 构造请求体
	reqBody := StreamRequest{
		Model: s.ModelID,
		Messages: []Message{
			{Role: "user", Content: question},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败：%v", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败：%v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}

	// 逐行读取流式响应
	scanner := bufio.NewScanner(resp.Body)
	fmt.Print("AI 回复：")
	for scanner.Scan() {
		line := scanner.Text()
		// 过滤掉空行和分隔符
		if line == "" || line == "data: [DONE]" {
			continue
		}
		// 去掉前缀 "data: "
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		// 解析 JSON
		var streamResp StreamResponse
		err := json.Unmarshal([]byte(line), &streamResp)
		if err != nil {
			continue // 忽略解析失败的行
		}

		// 输出内容
		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			fmt.Print(content)
			// 如果是结束标志，退出
			if streamResp.Choices[0].FinishReason == "stop" {
				break
			}
		}
	}

	fmt.Println("\n\n对话结束")
	return scanner.Err()
}

// 普通聊天（非流式）
func (s *AISDKConfig) Chat(ctx context.Context, question string) (string, error) {
	reqBody := Request{
		Model: s.ModelID,
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化失败：%v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败：%v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败：%v", err)
	}

	fmt.Println("返回结果：")
	fmt.Println(string(body))
	return string(body), nil
}

// chatwithPrompt
func (s *AISDKConfig) ChatWithPrompt(ctx context.Context, prompt string) (string, error) {
	//构造请求体，（非流式）
	type Request struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
	}
	reqBody := Request{
		Model: s.ModelID,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化失败：%v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败：%v", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}

	// 解析完整响应
	type Response struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}

	var respData Response
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(respData.Choices) == 0 {
		return "", fmt.Errorf("无返回内容")
	}

	return respData.Choices[0].Message.Content, nil
}
