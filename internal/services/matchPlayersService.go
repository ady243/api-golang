package services

import (
	"errors"

	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

type MatchPlayersService struct {
	DB *gorm.DB
}

// NewMatchPlayersService crée une nouvelle instance de MatchPlayersService
func NewMatchPlayersService(db *gorm.DB) *MatchPlayersService {
	return &MatchPlayersService{
		DB: db,
	}
}

// CreateMatchPlayers crée un nouveau joueur dans la table match_players
func (s *MatchPlayersService) CreateMatchPlayers(matchPlayers *models.MatchPlayers) error {
	if err := s.DB.Create(matchPlayers).Error; err != nil {
		return err
	}
	return nil
}

// GetMatchPlayersByMatchID retrieves all match players associated with the given match ID
func (service *MatchPlayersService) GetMatchPlayersByMatchID(matchID string) ([]models.MatchPlayers, error) {
	var matchPlayers []models.MatchPlayers
	if err := service.DB.Preload("Player").Where("match_id = ?", matchID).Find(&matchPlayers).Error; err != nil {
		return nil, err
	}
	return matchPlayers, nil
}

// CheckMatchExists vérifie si un match existe en fonction de son ID
func (s *MatchPlayersService) CheckMatchExists(matchID string) error {
	var match models.Matches
	if err := s.DB.First(&match, "id = ?", matchID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("match not found")
		}
		return err
	}
	return nil
}

// DeleteMatchPlayer supprime un joueur d'un match
func (s *MatchPlayersService) DeleteMatchPlayer(playerID, matchID string) error {
	if err := s.DB.Where("player_id = ? AND match_id = ?", playerID, matchID).Delete(&models.MatchPlayers{}).Error; err != nil {
		return err
	}
	return nil
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

// GetMatchesByPlayerID récupère les matchPlayers joués par playerID
func (s *MatchPlayersService) GetMatchesByPlayerID(playerID string) ([]models.MatchPlayers, error) {
	var matchPlayers []models.MatchPlayers
	if err := s.DB.Where("player_id = ?", playerID).Find(&matchPlayers).Error; err != nil {
		return nil, err
	}
	return matchPlayers, nil
}
