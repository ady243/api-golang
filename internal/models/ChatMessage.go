package models

type ChatMessage struct {
    PlayerID   string `json:"playerId"`
    Message    string `json:"message"`
    Timestamp  string `json:"timestamp"`
    Username   string `json:"username"`
    ProfilePic string `json:"profile_pic"`
}