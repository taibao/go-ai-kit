package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

const (
	apiKey = "418bbb73-29bd-42dd-9908-c8ee99657fb6"
	model  = "doubao-seed-2-0-code-preview-260215" // 你的模型ID
	apiUrl = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
)

func main() {
	// 初始化SDK
	sdk := NewAISDK(apiKey, model)

	// 上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*1000*time.Millisecond)
	defer cancel()
	// 输入问题
	fmt.Print("请输入问题：")
	var question string
	fmt.Scanln(&question)

	// 调用流式聊天
	err := sdk.StreamChat(ctx, question)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
		os.Exit(1)
	}
}
