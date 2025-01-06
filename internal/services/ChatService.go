package services

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "time"

    "github.com/ady243/teamup/internal/models"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

type ChatService struct {
    DB          *gorm.DB
    RedisClient *redis.Client
}

func NewChatService(db *gorm.DB, redisClient *redis.Client) *ChatService {
    return &ChatService{
        DB:          db,
        RedisClient: redisClient,
    }
}

func (s *ChatService) AddMessage(matchID, userID, message string) error {
    var matchPlayer models.MatchPlayers
    if err := s.DB.Preload("Player").First(&matchPlayer, "match_id = ? AND player_id = ?", matchID, userID).Error; err != nil {
        return errors.New("player not found in match")
    }

    fullMessage := models.ChatMessage{
        PlayerID:   userID,
        Message:    message,
        Timestamp:  time.Now().Format(time.RFC3339),
        Username:   matchPlayer.Player.Username,
        ProfilePic: matchPlayer.Player.ProfilePhoto,
    }

    msgJSON, err := json.Marshal(fullMessage)
    if err != nil {
        return err
    }

    key := "chat:" + matchID
    ctx := context.Background()
    if err := s.RedisClient.RPush(ctx, key, msgJSON).Err(); err != nil {
        return err
    }

    return nil
}

func (s *ChatService) GetMessages(matchID string, limit int) ([]models.ChatMessage, error) {
    ctx := context.Background()
    key := "chat:" + matchID

    msgs, err := s.RedisClient.LRange(ctx, key, 0, int64(limit-1)).Result()
    if err != nil {
        return nil, err
    }

    messages := make([]models.ChatMessage, len(msgs))
    for i, msg := range msgs {
        var chatMessage models.ChatMessage
        if err := json.Unmarshal([]byte(msg), &chatMessage); err != nil {
            return nil, err
        }
        messages[i] = chatMessage
    }

    return messages, nil
}

func (s *ChatService) AddUserToChat(matchID, userID string) error {
    ctx := context.Background()
    chatKey := "chat:" + matchID + ":users"

    if err := s.RedisClient.SAdd(ctx, chatKey, userID).Err(); err != nil {
        return err
    }

    if err := s.RedisClient.Expire(ctx, chatKey, time.Hour*24*7).Err(); err != nil {
        log.Printf("Error setting expiration for chat key: %v", err)
        return err
    }

    return nil
}


//delete chat messages
func (s *ChatService) DeleteChatMessages(matchID string) error {
    ctx := context.Background()
    key := "chat:" + matchID

    if err := s.RedisClient.Del(ctx, key).Err(); err != nil {
        return err
    }

    return nil
}


func (s *ChatService) GetParticipants(matchID string) ([]models.Users, error) {
    var participants []models.Users
    if err := s.DB.Joins("JOIN match_players ON match_players.player_id = users.id").
        Where("match_players.match_id = ?", matchID).
        Find(&participants).Error; err != nil {
        return nil, err
    }
    return participants, nil
}