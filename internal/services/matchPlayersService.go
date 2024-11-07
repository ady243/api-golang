package services

import (
	"github.com/ady243/teamup/internal/models"

	"gorm.io/gorm"
)

type MatchPlayersService struct {
	DB *gorm.DB
}

func NewMatchPlayersService(db *gorm.DB) *MatchPlayersService {
	return &MatchPlayersService{
		DB: db,
	}
}

// CreateMatchPlayers crée un nouveau player dans la table match_players
func (s *MatchPlayersService) CreateMatchPlayers(matchPlayers *models.MatchPlayers) error {
	if err := s.DB.Create(matchPlayers).Error; err != nil {
		return err
	}
	return nil
}

func (s *MatchPlayersService) GetMatchPlayersByMatchID(matchID string) ([]models.MatchPlayers, error) {
	var matchPlayers []models.MatchPlayers
	if err := s.DB.Where("match_id = ?", matchID).Find(&matchPlayers).Error; err != nil {
		return nil, err
	}
	return matchPlayers, nil
}

func (s *MatchPlayersService) GetMatchPlayerByID(matchPlayerID string) (*models.MatchPlayers, error) {
	var matchPlayer models.MatchPlayers
	if err := s.DB.Where("id = ?", matchPlayerID).First(&matchPlayer).Error; err != nil {
		return nil, err
	}
	return &matchPlayer, nil
}

func (s *MatchPlayersService) UpdateMatchPlayer(matchPlayer *models.MatchPlayers) error {
	if err := s.DB.Save(matchPlayer).Error; err != nil {
		return err
	}
	return nil
}

func (s *MatchPlayersService) DeleteMatchPlayer(matchPlayerID string) error {
	if err := s.DB.Delete(&models.MatchPlayers{}, "id = ?", matchPlayerID).Error; err != nil {
		return err
	}
	return nil
}
