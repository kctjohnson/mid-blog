package content

import (
	"context"
	"strings"

	"github.com/ayush6624/go-chatgpt"
)

// This package exists to hold the tools for creating bloggers, posts, titles, etc
// It just lets me refactor a bit of the longer pieces of code into its own package

// Kinda funny name ngl
type ContentCreator struct {
	gptClient *chatgpt.Client
}

func NewContentCreator(gptClient *chatgpt.Client) *ContentCreator {
	return &ContentCreator{
		gptClient: gptClient,
	}
}

func (c ContentCreator) GenerateFirstName(temperature float64) (string, error) {
	firstNameResp, err := c.gptClient.Send(
		context.Background(),
		&chatgpt.ChatCompletionRequest{
			Temperature: temperature,
			Model:       chatgpt.GPT35Turbo,
			Messages: []chatgpt.ChatMessage{
				{
					Role:    chatgpt.ChatGPTModelRoleSystem,
					Content: "I am a blogger generating system. I can generate bloggers, and their posts. Every blogger I generate will be entirely random.",
				},
				{
					Role:    chatgpt.ChatGPTModelRoleUser,
					Content: "In one word, give me a random first name. It can be any name, real or made up, and should only contain letters. Do not give me anything else, and do not give any punctuation.",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	firstName := firstNameResp.Choices[0].Message.Content

	return firstName, nil
}

func (c ContentCreator) GenerateLastName(temperature float64) (string, error) {
	lastNameResp, err := c.gptClient.Send(
		context.Background(),
		&chatgpt.ChatCompletionRequest{
			Temperature: temperature,
			Model:       chatgpt.GPT35Turbo,
			Messages: []chatgpt.ChatMessage{
				{
					Role:    chatgpt.ChatGPTModelRoleSystem,
					Content: "I am a blogger generating system. I can generate bloggers, and their posts. Every blogger I generate will be entirely random.",
				},
				{
					Role:    chatgpt.ChatGPTModelRoleUser,
					Content: "In one word, give me a random last name. It can be any name, real or made up, and should only contain letters. Do not give me anything else, and do not give any punctuation.",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	lastName := lastNameResp.Choices[0].Message.Content

	return lastName, nil
}

func (c ContentCreator) GenerateEmail(temperature float64) (string, error) {
	emailResp, err := c.gptClient.Send(
		context.Background(),
		&chatgpt.ChatCompletionRequest{
			Temperature: temperature,
			Model:       chatgpt.GPT35Turbo,
			Messages: []chatgpt.ChatMessage{
				{
					Role:    chatgpt.ChatGPTModelRoleSystem,
					Content: "I am a blogger generating system. I can generate bloggers, and their posts. Every blogger I generate will be entirely random.",
				},
				{
					Role:    chatgpt.ChatGPTModelRoleUser,
					Content: "Give me the name for a new email. It can be any name, real or made up. Do not give me anything else, and do not give any punctuation. Do not give the @ or anything that would follow it afterwards in an email, I just want the name part. Do not give anything else but the name part.",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	email := emailResp.Choices[0].Message.Content

	return email, nil
}

func (c ContentCreator) GenerateBio(
	temperature float64,
	firstName string,
	lastName string,
) (string, error) {
	bioResp, err := c.gptClient.Send(
		context.Background(),
		&chatgpt.ChatCompletionRequest{
			Temperature: temperature,
			Model:       chatgpt.GPT35Turbo,
			Messages: []chatgpt.ChatMessage{
				{
					Role:    chatgpt.ChatGPTModelRoleSystem,
					Content: "I am a blogger generating system. I can generate bloggers, and their posts. Every blogger I generate will be entirely random.",
				},
				{
					Role:    chatgpt.ChatGPTModelRoleUser,
					Content: "Write me a random bio that is AT MAX 255 characters long. I want the bio to have the blogger leaning towards one or two specific subjects of interest, or things that they like. Only include the blogger's bio in your response, do not include anything else, do not add any surrounding quotes around it. The blogger's name is " + firstName + " " + lastName + ".",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	bio := bioResp.Choices[0].Message.Content
	return bio, nil
}

func (c ContentCreator) GenerateTitle(temperature float64, bio string) (string, error) {
	titleResp, err := c.gptClient.Send(context.Background(), &chatgpt.ChatCompletionRequest{
		Temperature: temperature,
		Model:       chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "I am a blog generating system. I can generate blog post titles, and their content. Every title I generate will be entirely random.",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleUser,
				Content: "Please make me a blog post title, ranging from 5-8 words, and it can be on any topic. Do not include any extra words outside of the title, and it doesn't need to be wrapped in quotes. Base the topic for the title based on the bio of the blogger. The bio is \"" + bio + "\".",
			},
		},
	})
	if err != nil {
		return "", err
	}

	title := strings.ReplaceAll(titleResp.Choices[0].Message.Content, "\"", "")
	return title, nil
}

func (c ContentCreator) GeneratePost(temperature float64, title string) (string, error) {
	contentResp, err := c.gptClient.Send(context.Background(), &chatgpt.ChatCompletionRequest{
		Temperature: 1.0,
		Model:       chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "I am a blog generating system. I can generate blog post titles, and their content. Every title I generate will be entirely random.",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleUser,
				Content: "Based on this title, please write a blog post with 5 paragraphs, each with 5-8 sentences. The content should be relevant to the title, and should be entirely random. Do not include the title at the top of the content. Do not include any extra words outside of the content, and it doesn't need to be wrapped in quotes. The title is \"" + title + "\".",
			},
		},
	})
	if err != nil {
		return "", err
	}
	content := contentResp.Choices[0].Message.Content

	return content, nil
}
