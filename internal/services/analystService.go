package services

import (
	"errors"

	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

type AnalystService struct {
	DB *gorm.DB
}

func NewAnalystService(db *gorm.DB) *AnalystService {
	return &AnalystService{
		DB: db,
	}
}

// CreateEvent crée un nouvel enregistrement dans la table Analyst
func (s *AnalystService) CreateEvent(event *models.Analyst) error {
	return s.DB.Create(event).Error
}

// GetEventsByMatchID renvoie tous les events pour un match donné
func (s *AnalystService) GetEventsByMatchID(matchID string) ([]models.Analyst, error) {
	var events []models.Analyst
	// Précharger la relation Player et Match si vous le souhaitez
	if err := s.DB.Preload("Player").Preload("Match").
		Where("match_id = ? AND deleted_at IS NULL", matchID).
		Find(&events).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no events found")
		}
		return nil, err
	}
	return events, nil
}

// GetEventByID renvoie un event en fonction de son ID
func (s *AnalystService) GetEventByID(eventID string) (*models.Analyst, error) {
	var event models.Analyst
	if err := s.DB.Preload("Player").Preload("Match").
		Where("id = ? AND deleted_at IS NULL", eventID).
		First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEventsByPlayerID renvoie tous les événements associés à un joueur donné
func (s *AnalystService) GetEventsByPlayerID(playerID string) ([]models.Analyst, error) {
	var events []models.Analyst
	// Précharger Player et Match si nécessaire
	if err := s.DB.Preload("Player").Preload("Match").
		Where("player_id = ? AND deleted_at IS NULL", playerID).
		Find(&events).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no events found for this player")
		}
		return nil, err
	}
	return events, nil
}

// UpdateEvent met à jour un event déjà existant
func (s *AnalystService) UpdateEvent(event *models.Analyst) error {
	return s.DB.Save(event).Error
}

// DeleteEvent supprime vraiment l’event (Hard Delete, si besoin)
func (s *AnalystService) DeleteEvent(eventID string) error {
	if err := s.DB.Where("id = ?", eventID).Delete(&models.Analyst{}).Error; err != nil {
		return err
	}
	return nil
}
