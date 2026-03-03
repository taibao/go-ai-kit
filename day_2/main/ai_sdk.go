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

// 新建SDK实例
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
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
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

		// 解析JSON
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

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
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

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBuffer(jsonData))
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
	return string(body), nil
}
