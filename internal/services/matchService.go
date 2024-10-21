package services

import (
	"errors"
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
