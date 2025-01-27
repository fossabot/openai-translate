/*
 * @Author: Vincent Young
 * @Date: 2023-11-13 11:16:26
 * @LastEditors: Vincent Young
 * @LastEditTime: 2023-11-13 12:08:36
 * @FilePath: /openai-translate/main.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 * 
 * Copyright © 2023 by Vincent, All Rights Reserved. 
 */
package main

import (
    "context"
    "net/http"
    openai "github.com/sashabaranov/go-openai"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"flag"
	"os"
	"fmt"
)

type ResData struct {
	TransText  string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

func translator(apiKey string, targetLang string, transText string)(string, error) {
    c := openai.NewClient(apiKey)
    resp, err := c.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleSystem,
                    Content: "You're a translator. Translate to " + targetLang + ".",
                },
				{
                    Role:    openai.ChatMessageRoleUser,
                    Content: transText,
                },
            },
        },
    )

    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}

func main(){
	// Define a command line flag
	apiKeyFlag := flag.String("apiKey", "", "API key for OpenAI")
	flag.Parse()

	// First try to get the API key from the command line flag
	apiKey := *apiKeyFlag

	// If it's not provided, try to get it from the environment variable
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_KEY")
	}

	// If the API key is still empty, return an error and exit
	if apiKey == "" {
		fmt.Println("Error: No API key provided. Set the apiKey flag or the OPENAI_KEY environment variable.")
		os.Exit(1)
	}

	fmt.Println("Starting server on port 23333. Made by Vincent.")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Hello",
		})
	})

	r.POST("/translate", func(c *gin.Context) {
		req := ResData{}
		c.BindJSON(&req)
		sourceLang := req.SourceLang
		targetLang := req.TargetLang
		translateText := req.TransText
		targetText, _ := translator(apiKey, targetLang, translateText)
		
		if targetText == "" {
			c.JSON(http.StatusTooManyRequests, gin.H{ // 429 Too Many Requests
				"code":    http.StatusTooManyRequests,
				"message": "Translation limit exceeded or service unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":         http.StatusOK,
			"data":         targetText,
			"source_lang":  sourceLang,
			"target_lang":  targetLang,
		})

	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Path not found",
		})
	})

	r.Run(":23333")
}