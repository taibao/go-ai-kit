package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	apiKey = "418bbb73-29bd-42dd-9908-c8ee99657fb6"
	model  = "doubao-seed-2-0-code-preview-260215" // 你的模型ID
)

// 定义劳动合同审核的JSON结构（用于解析AI回复）
type ContractCheckResult struct {
	RiskPoints []struct {
		Level   string `json:"级别"`
		Clause  string `json:"条款"`
		Problem string `json:"问题"`
		Suggest string `json:"建议"`
	} `json:"风险点"`
	NoRiskClauses []string `json:"无风险条款"`
}

func main() {
	// 初始化SDK
	sdk := NewAISDK(apiKey, model)

	// 上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*1000*time.Millisecond)
	defer cancel()

	// 测试：强制JSON输出 + Token控制
	contractContent := `劳动合同
1. 合同期限：3年，试用期6个月；
2. 试用期工资：正式工资的50%；
3. 加班工资：按当地最低工资标准计算；
4. 社保：入职满1年缴纳。`

	// 生成Prompt
	prompt := GeneratePrompt("contract_check", map[string]string{
		"contract_content": contractContent,
	})

	// 调用带Token控制的函数（限制输出Token为1000）
	resp, token, err := sdk.ChatWithTokenControl(ctx, prompt, 1000)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
		return
	}
	fmt.Printf("=== 劳动合同审核结果（Token：%d）===\n%s\n", token, resp)

	// 解析JSON（验证是否可正常解析）
	var result ContractCheckResult
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		fmt.Printf("JSON解析失败：%v\n", err)
	} else {
		fmt.Println("\n=== 解析后的风险点 ===")
		for _, risk := range result.RiskPoints {
			fmt.Printf("级别：%s，问题：%s\n", risk.Level, risk.Problem)
		}
	}

	// 测试多模型切换
	//sdk.SwitchModel("doubao-pro")
	//fmt.Println("\n已切换到豆包模型，可直接调用")

	sdk.SwitchModel("tongyi")
	fmt.Println("\n已切换到通义千问模型，可直接调用")
}
