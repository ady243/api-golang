package services

import (
	"errors"

	"github.com/ady243/teamup/internal/models"
	"github.com/oklog/ulid/v2"
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
func (s *MatchService) DeleteMatch(matchID ulid.ULID) error {
	if err := s.DB.Where("id = ?", matchID).Delete(&models.Matches{}).Error; err != nil {
		return err
	}
	return nil
}
