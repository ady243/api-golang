package services

import (
	"errors"
	"math"
	"time"

	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

// MatchService fournit les services pour gérer les matchs
type MatchService struct {
	DB *gorm.DB
}

// NewMatchService crée une nouvelle instance de MatchService
func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{
		DB: db,
	}
}

// CreateMatch crée un nouveau match dans la base de données
func (s *MatchService) CreateMatch(match *models.Matches) error {
	if err := s.DB.Create(match).Error; err != nil {
		return err
	}
	return nil
}

// Méthode pour récupérer tous les matchs
func (s *MatchService) GetAllMatches() ([]models.Matches, error) {
	var matches []models.Matches
	if err := s.DB.Find(&matches).Error; err != nil {
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

// DeleteMatch supprime un match par son ID (soft delete)
func (s *MatchService) DeleteMatch(matchID string) error {
	if err := s.DB.Model(&models.Matches{}).Where("id = ?", matchID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

// Fonction de Haversine pour calculer la distance en km entre deux points
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Rayon de la Terre en kilomètres
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
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
