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

	//Ici on fait appel à l'API OpenAI pour suggérer une formation
	prompt := "Voici les statistiques des joueurs de football :\n"
	for _, player := range players {
		prompt += fmt.Sprintf("Joueur: %s, Pac: %d, Sho: %d, Pas: %d, Dri: %d, Def: %d, Phy: %d\n",
			player.Username, player.Pac, player.Sho, player.Pas, player.Dri, player.Def, player.Phy)
	}
	prompt += "Donne moi la formation la plus adaptée suivant leurs profil par exemple Formation : 4-4-2. Pour chacun des joueurs de la liste, donne moi sa positions idéale sur le terrain, par exemple : Joueur : Username, Position : Defenseur gauche."
	// Ici on fait la gestion des erreurs et on initialise la version du modèle à utiliser
	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
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

	// Ici on extrait la formation suggérée par OpenAI
	if len(resp.Choices) > 0 {
		formations = append(formations, resp.Choices[0].Message.Content)
	} else {
		return nil, fmt.Errorf("aucune formation suggérée par OpenAI")
	}

	return formations, nil
}
