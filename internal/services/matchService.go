package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

// MatchService fournit les services pour gérer les matchs
type MatchService struct {
	DB          *gorm.DB
	ChatService *ChatService
}

func NewMatchService(db *gorm.DB, chatService *ChatService) *MatchService {
	return &MatchService{
		DB:          db,
		ChatService: chatService,
	}
}

// CreateMatch crée un nouveau match dans la base de données
func (s *MatchService) CreateMatch(match *models.Matches) error {
	if err := s.DB.Create(match).Error; err != nil {
		return err
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
	if err := s.DB.Where("id = ?", matchID).First(&match).Error; err != nil {
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
