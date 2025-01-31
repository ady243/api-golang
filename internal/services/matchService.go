package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// MatchService fournit les services pour gérer les matchs
type MatchService struct {
	DB          *gorm.DB
	ChatService *ChatService
	RedisClient *redis.Client
}

func NewMatchService(db *gorm.DB, chatService *ChatService, redisClient *redis.Client) *MatchService {
	return &MatchService{
		DB:          db,
		ChatService: chatService,
		RedisClient: redisClient,
	}
}

// CreateMatch crée un nouveau match dans la base de données et met à jour le rôle de l'utilisateur
func (s *MatchService) CreateMatch(match *models.Matches, userID string) error {
	// Générer un nouvel ID pour le match
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	newID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	match.ID = newID.String()
	match.OrganizerID = userID
	match.Status = models.Upcoming

	// Ajouter le rôle "organizer" à l'utilisateur
	var user models.Users
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// Créer le match dans la base de données
	if err := s.DB.Create(match).Error; err != nil {
		return err
	}

	// Ajouter l'utilisateur comme joueur dans la table match_players
	matchPlayer := &models.MatchPlayers{
		ID:       ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
		MatchID:  match.ID,
		PlayerID: userID,
	}
	if err := s.DB.Create(matchPlayer).Error; err != nil {
		return fmt.Errorf("failed to add organizer to match players: %w", err)
	}

	return nil
}

// CheckAndResetRole vérifie si le match a expiré et réinitialise le rôle de l'utilisateur
func (s *MatchService) VerifyAndResetRole(matchID string) error {
	var match models.Matches
	if err := s.DB.Where("id = ?", matchID).First(&match).Error; err != nil {
		return err
	}

	// Vérifier si le match a expiré
	if match.EndTime.Before(time.Now()) {
		var user models.Users
		if err := s.DB.Where("id = ?", match.OrganizerID).First(&user).Error; err != nil {
			return err
		}

		// Mettre à jour le statut du match à "expired"
		match.Status = models.Expired
		if err := s.DB.Save(&match).Error; err != nil {
			return err
		}
	}

	return nil
}

// Méthode pour récupérer tous les matchs avec préchargement des informations de l'organisateur
func (s *MatchService) GetAllMatches() ([]models.Matches, error) {
	var matches []models.Matches
	if err := s.DB.Preload("Organizer").Where("deleted_at IS NULL").Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// GetMatchByID récupère un match par son ID
func (s *MatchService) GetMatchByID(matchID string) (*models.Matches, error) {
	var match models.Matches
	if err := s.DB.Preload("Organizer").Where("id = ?", matchID).First(&match).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("match not found")
		}
		return nil, err
	}
	return &match, nil
}

// UpdateMatch met à jour un match existant dans la base de données
func (s *MatchService) UpdateMatch(match *models.Matches) error {
	if err := s.DB.Save(match).Error; err != nil {
		return err
	}
	return nil
}

// DeleteMatch supprime un match, si l'utilisateur est l'organisateur.
// L'ID du match et l'ID de l'utilisateur sont fournis en paramètres.
// La méthode renvoie une erreur si l'utilisateur n'est pas l'organisateur
// ou si le match n'est pas trouvé.
func (s *MatchService) DeleteMatch(matchID, userID string) error {
	// Récupère le match à partir de son ID
	match, err := s.GetMatchByID(matchID)
	if err != nil {
		return err
	}

	// Vérifie si l'utilisateur connecté est l'organisateur du match
	if match.OrganizerID != userID {
		return fmt.Errorf("you are not authorized to delete this match")
	}

	now := time.Now()
	match.DeletedAt = &now
	if err := s.DB.Save(&match).Error; err != nil {
		return err
	}

	// Supprime les joueurs du match
	if err := s.DB.Where("match_id = ?", matchID).Delete(&models.MatchPlayers{}).Error; err != nil {
		return err
	}

	// Supprime les messages de chat associés au match
	if err := s.ChatService.DeleteChatMessages(matchID); err != nil {
		return err
	}

	return nil
}

// Fonction de Haversine pour calculer la distance en km entre deux points
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Rayon de la Terre en kilomètres
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon1 - lon2) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// Trouver les matchs à proximité d'une position
func (s *MatchService) FindNearbyMatches(lat, lon, radius float64) ([]models.Matches, error) {
	var matches []models.Matches
	if err := s.DB.Find(&matches).Error; err != nil {
		return nil, err
	}

	var nearbyMatches []models.Matches
	for _, match := range matches {
		distance := calculateDistance(lat, lon, match.Latitude, match.Longitude)
		if distance <= radius {
			nearbyMatches = append(nearbyMatches, match)
		}
	}
	return nearbyMatches, nil
}

// JoinMatch ajoute un joueur à un match
func (s *MatchService) JoinMatch(matchID, userID string) error {
	var count int64
	s.DB.Model(&models.MatchPlayers{}).Where("match_id = ? AND player_id = ?", matchID, userID).Count(&count)
	if count > 0 {
		return nil
	}

	player := models.MatchPlayers{
		MatchID:  matchID,
		PlayerID: userID,
	}

	return s.DB.Create(&player).Error
}

// IsUserInMatch vérifie si l'utilisateur est dans le match
func (s *MatchService) IsUserInMatch(matchID, userID string) error {
	var count int64
	if err := s.DB.Model(&models.MatchPlayers{}).Where("match_id = ? AND player_id = ?", matchID, userID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user is not in the match")
	}
	return nil
}

// GetMatchByOrganizerID récupère les matchs par l'ID de l'organisateur
func (s *MatchService) GetMatchByOrganizerID(organizerID string) ([]models.Matches, error) {
	var matches []models.Matches
	if err := s.DB.Where("organizer_id = ? AND deleted_at IS NULL", organizerID).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// GetMatchByOrganizerID récupère les matchs par l'ID de l'organisateur
func (s *MatchService) GetMatchByRefereeID(refereeID string) ([]models.Matches, error) {
	var matches []models.Matches
	if err := s.DB.Where("referee_id = ? AND deleted_at IS NULL", refereeID).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// PutRefereeID met à jour la table Matches pour y inscrire le RefereeID
func (s *MatchService) PutRefereeID(matcheID string, refereeID string) error {
	if err := s.DB.Model(&models.Matches{}).
		Where("id = ? AND deleted_at IS NULL", matcheID).
		Update("referee_id", refereeID).Error; err != nil {
		return err
	}
	return nil
}

func (s *MatchService) NotifyMatchStatusUpdate(matchID string, status string) error {
	message := map[string]string{
		"match_id": matchID,
		"status":   status,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.RedisClient.Publish(context.Background(), "match_status_updates", messageJSON).Err()
}

func (s *MatchService) UpdateMatchStatuses() error {
	var matches []models.Matches
	if err := s.DB.Find(&matches).Error; err != nil {
		return err
	}

	now := time.Now()
	for _, match := range matches {
		if match.Status != models.Completed && match.Status != models.Expired {
			matchStart := time.Date(match.MatchDate.Year(), match.MatchDate.Month(), match.MatchDate.Day(), match.MatchTime.Hour(), match.MatchTime.Minute(), match.MatchTime.Second(), 0, time.UTC)
			matchEnd := time.Date(match.MatchDate.Year(), match.MatchDate.Month(), match.MatchDate.Day(), match.EndTime.Hour(), match.EndTime.Minute(), match.EndTime.Second(), 0, time.UTC)

			if now.After(matchStart) && now.Before(matchEnd) {
				match.Status = models.Ongoing
			} else if now.After(matchEnd) {
				match.Status = models.Completed
			} else if now.Before(matchStart) {
				match.Status = models.Upcoming
			}

			if err := s.DB.Save(&match).Error; err != nil {
				return err
			}

			// Notifier les clients de la mise à jour du statut du match
			if err := s.NotifyMatchStatusUpdate(match.ID, string(match.Status)); err != nil {
				log.Println("Erreur lors de la notification de la mise à jour du statut du match:", err)
			}
		}
	}
	return nil
}

func (s *MatchService) AssignReferee(matchID, organizerID, refereeID string) error {
	// Vérifier si l'utilisateur est l'organisateur du match
	var match models.Matches
	if err := s.DB.Where("id = ? AND organizer_id = ?", matchID, organizerID).First(&match).Error; err != nil {
		return errors.New("unauthorized or match not found")
	}

	// Vérifier si l'utilisateur est un participant du match
	var count int64
	if err := s.DB.Model(&models.MatchPlayers{}).Where("match_id = ? AND player_id = ?", matchID, refereeID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user is not a participant in the match")
	}

	// Mettre à jour le referee_id du match
	if err := s.DB.Model(&models.Matches{}).Where("id = ?", matchID).Update("referee_id", refereeID).Error; err != nil {
		return err
	}

	return nil
}

func (s *MatchService) LeaveMatch(matchID, userID string) error {
	// Supprimer l'utilisateur de la liste des joueurs du match
	if err := s.DB.Where("match_id = ? AND player_id = ?", matchID, userID).Delete(&models.MatchPlayers{}).Error; err != nil {
		return err
	}
	return nil
}
