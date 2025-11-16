package repository

import (
	"github.com/ashurov-imomali/sbscribtion-service/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(subscription *models.Subscription) error
	GetByID(id uuid.UUID) (*models.Subscription, bool, error)
	GetByFilter(filter models.SubscriptionFilter) ([]models.Subscription, bool, error)
	Update(id uuid.UUID, updates map[string]interface{}) (*models.Subscription, bool, error)
	Delete(id uuid.UUID) (bool, error)
	GetTotal(from, to string, userID uuid.UUID, service string) (*models.Total, error)
}
