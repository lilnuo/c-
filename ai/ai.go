package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

var ChatModel *ark.ChatModel

// InitAI 初始化 AI 客户端 
func InitAI() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 没有找到 .env 文件")
	}

	ctx := context.Background()

	apiKey := os.Getenv("ARK_API_KEY")
	modelID := os.Getenv("ARK_MODEL_ID") 

	fmt.Printf("正在初始化火山引擎, Model(Endpoint): %s\n", modelID)
	ChatModel, err = ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelID, 
	})

	if err != nil {
		return err
	}
	return nil
}


func Answer(question string) (string, error) {

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}。"),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我{task}。",
		},
	)

	variables := map[string]any{
		"role": "C语言编程专家。你的唯一任务是写代码。要求：1.只输出纯粹的C语言源代码。2.严禁输出Markdown标记。3.不要任何解释、前言或后缀。4.代码必须包含必要的头文件。5.尤其注意格式问题,比如空格等",
		"task": question,
	}

	messages, err := template.Format(context.Background(), variables)
	if err != nil {
		return "", err
	}

	result, err := ChatModel.Generate(context.Background(), messages)
	if err != nil {
		log.Printf("生成失败，err:%v", err)
		return "", err
	}
	if len(result.Content) == 0 {
		return "", err
	}

	rawContent := CleanCode(result.Content)
	return rawContent, nil

}

func CleanCode(input string) string {
	input = strings.ReplaceAll(input, "```c", "")
	input = strings.ReplaceAll(input, "```", "")
	return strings.TrimSpace(input)
}
