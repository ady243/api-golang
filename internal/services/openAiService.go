package services

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/ady243/teamup/internal/models"
    openai "github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
    client *openai.Client
}

func NewOpenAIService() *OpenAIService {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable is not set")
    }
    client := openai.NewClient(apiKey)
    return &OpenAIService{
        client: client,
    }
}

func (s *OpenAIService) SuggestFormations(players []models.Users) ([]string, error) {
    var formations []string

    // Construire le prompt pour l'API OpenAI
    prompt := "Voici les statistiques des joueurs :\n"
    for _, player := range players {
        prompt += fmt.Sprintf("Joueur: %s, Pac: %d, Sho: %d, Pas: %d, Dri: %d, Def: %d, Phy: %d\n",
            player.Username, player.Pac, player.Sho, player.Pas, player.Dri, player.Def, player.Phy)
    }
    prompt += "Suggérez une formation pour ces joueurs."
    // Appeler l'API OpenAI avec CreateChatCompletion
    resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    "system",
                Content: "You are a football coach assistant.",
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
        MaxTokens: 100,
    })
    if err != nil {
        log.Printf("Erreur lors de l'appel à l'API OpenAI: %v", err)
        return nil, err
    }

    // Extraire les formations de la réponse
    if len(resp.Choices) > 0 {
        formations = append(formations, resp.Choices[0].Message.Content)
    } else {
        return nil, fmt.Errorf("Aucune formation suggérée par OpenAI")
    }

    return formations, nil
}