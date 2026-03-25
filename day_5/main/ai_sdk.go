package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
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

// 新增：Token估算函数（中文字符数换算）
func EstimateToken(text string) int {
	// 中文：1字 ≈ 1.3 Token，英文：1词 ≈ 1 Token
	chineseChars := 0
	englishWords := 0
	for _, c := range text {
		if c >= '\u4e00' && c <= '\u9fff' {
			chineseChars++
		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			englishWords++
		}
	}
	// 估算总Token
	token := int(math.Round(float64(chineseChars)*1.3)) + englishWords
	return token
}

// 新增：带Token控制的聊天函数
func (s *AISDKConfig) ChatWithTokenControl(ctx context.Context, prompt string, maxOutputToken int) (string, int, error) {
	// 1. 估算输入Token
	inputToken := EstimateToken(prompt)
	fmt.Printf("输入Token估算：%d\n", inputToken)

	// 2. 构造请求体（增加max_tokens参数）
	type Request struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Stream      bool      `json:"stream"`
		MaxTokens   int       `json:"max_tokens"`  // 新增：限制输出Token
		Temperature float64   `json:"temperature"` // 温度：0~1，越低越精准
	}

	reqBody := Request{
		Model:       s.ModelID,
		Messages:    []Message{{Role: "user", Content: prompt}},
		Stream:      false,
		MaxTokens:   maxOutputToken,
		Temperature: 0.1, // 企业场景用低温度，保证精准
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("序列化失败: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 3. 增加重试逻辑（最多3次）
	var resp *http.Response
	var errReq error
	for i := 0; i < 3; i++ {
		client := &http.Client{}
		resp, errReq = client.Do(req)
		if errReq == nil && resp.StatusCode == http.StatusOK {
			break
		}
		fmt.Printf("请求失败，第%d次重试...\n", i+1)
		time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
	}
	if errReq != nil {
		return "", 0, fmt.Errorf("请求失败（重试3次）: %v", errReq)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("状态码错误: %d", resp.StatusCode)
	}

	// 4. 解析响应
	type Response struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage struct { // 部分API会返回真实Token使用量
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	var respData Response
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", 0, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(respData.Choices) == 0 {
		return "", 0, fmt.Errorf("无返回内容")
	}

	// 5. 统计输出Token（优先用API返回的真实值，否则估算）
	outputToken := respData.Usage.CompletionTokens
	if outputToken == 0 {
		outputToken = EstimateToken(respData.Choices[0].Message.Content)
	}
	fmt.Printf("输出Token：%d\n", outputToken)

	return respData.Choices[0].Message.Content, outputToken, nil
}

// 新增：多模型适配（一键切换豆包/通义千问）
func (s *AISDKConfig) SwitchModel(modelType string) {
	switch modelType {
	case "doubao-pro":
		s.ModelID = "doubao-seed-2-0-code-preview-260215"
		s.APIUrl = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
	case "tongyi":
		s.ModelID = "sk-fcdb23bbc0c34222a4e30cfbebf3a10c"
		s.APIUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	case "kimi":
		s.ModelID = "sk-FsPJ9fqCyXBxnWJgcgddJMfyD8W5BUw2Da9cBCpKpIZDqJeA"
		s.APIUrl = ""
	}
}
