package main

import (
	"context"
	"fmt"
	"time"
)

const (
	apiKey = "418bbb73-29bd-42dd-9908-c8ee99657fb6"
	model  = "doubao-seed-2-0-code-preview-260215" // 你的模型ID
)

func main() {
	// 初始化SDK
	sdk := NewAISDK(apiKey, model)

	// 上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*1000*time.Millisecond)
	defer cancel()
	// 输入问题
	// 测试1：生成请假制度
	//prompt1 := GeneratePrompt("leave_policy", map[string]string{})
	//fmt.Println("=== 测试1：企业请假制度 ===")
	//resp1, err := sdk.ChatWithPrompt(ctx, prompt1)
	//if err != nil {
	//	fmt.Printf("错误：%v\n", err)
	//	return
	//}
	//fmt.Println(resp1)

	// 测试2：生成Go面试问题
	//prompt2 := GeneratePrompt("interview_questions", map[string]string{})
	//fmt.Println("\n=== 测试2：Go面试问题 ===")
	//resp2, err := sdk.ChatWithPrompt(ctx, prompt2)
	//if err != nil {
	//	fmt.Printf("错误：%v\n", err)
	//	return
	//}
	//fmt.Println(resp2)

	//// 测试3：客户投诉回复（自定义投诉内容）
	//prompt3 := GeneratePrompt("complaint_reply", map[string]string{
	//	"complaint_content": "我买的衣服收到后尺寸不对，联系客服3天没人回复，要求退货+赔偿！",
	//})
	//fmt.Println("\n=== 测试3：客户投诉回复 ===")
	//resp3, err := sdk.ChatWithPrompt(ctx, prompt3)
	//if err != nil {
	//	fmt.Printf("错误：%v\n", err)
	//	return
	//}
	//	fmt.Println(resp3)

	//测试4：
	prompt4 := GeneratePrompt("private_operation", map[string]string{})
	fmt.Println("\n=== 测试4：私域运营SOP ===")
	resp4, err := sdk.ChatWithPrompt(ctx, prompt4)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
		return
	}
	fmt.Println(resp4)
}
