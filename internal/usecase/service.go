package usecase

import (
	"github.com/ashurov-imomali/sbscribtion-service/internal/models"
	"github.com/ashurov-imomali/sbscribtion-service/internal/repository"
	"github.com/ashurov-imomali/sbscribtion-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Service struct {
	repo repository.Repository
	l    logger.Logger
}

func New(repo repository.Repository, l logger.Logger) *Service {
	return &Service{repo: repo, l: l}
}

func (s *Service) CreateSubscription(sub *models.Subscription) (int, error) {
	startDate, err := time.Parse("2006-01-02", sub.StartDate)
	if errMessage, ok := s.validateSubscribe(map[string]bool{
		"service_name is required": sub.ServiceName != "",
		"price must be >= 0":       sub.Price > 0,
		"user_id is required":      sub.UserID != uuid.Nil,
		"start_date is required":   err == nil,
	}); !ok {
		return 422, errors.New(errMessage)
	}

	if sub.EndDate != nil {
		if parse, err := time.Parse("2006-01-02", *sub.EndDate); err == nil && !startDate.Before(parse) {
			return 422, errors.New("End Date have to been > Start Date")
		}
	}

	if err := s.repo.Create(sub); err != nil {
		s.l.Errorf("Error while create subscription in bd. ERROR: %v", err)
		return 500, errors.New("INTERNAL_SERVER_ERROR")
	}
	return 200, nil
}

func (s *Service) validateSubscribe(mp map[string]bool) (string, bool) {
	for k, v := range mp {
		if !v {
			return k, false
		}
	}

	return "All is ok", true
}

func (s *Service) GetSubscribe(id uuid.UUID) (*models.Subscription, int, error) {
	sub, notFound, err := s.repo.GetByID(id)
	if notFound {
		return nil, 404, errors.New("NOT_FOND")
	}
	if err != nil {
		s.l.Errorf("Error while get subscribtion from DB. ERROR: %v", err)
		return nil, 500, errors.New("INTERNAL_SERVER_ERROR")
	}
	return sub, 200, nil
}

func (s *Service) GetSubscriptions(filters models.SubscriptionFilter) ([]models.Subscription, int, error) {
	subscriptions, notFound, err := s.repo.GetByFilter(filters)
	if notFound {
		return nil, 404, errors.New("NOT_FOUND")
	}
	if err != nil {
		s.l.Errorf("Error while get subscription with filters. Error: %v", err)
		return nil, 500, errors.New("INTERNAL_SERVER_ERROR")
	}

	return subscriptions, 200, nil
}

func (s *Service) UpdateSubscription(sub models.Subscription) (*models.Subscription, int, error) {
	updates := make(map[string]interface{})

	if sub.ID == uuid.Nil {
		return nil, 422, errors.New("ID required")
	}

	if len(strings.TrimSpace(sub.ServiceName)) > 0 {
		updates["service_name"] = sub.ServiceName
	}

	if sub.Price > 0 {
		if sub.Price < 0 {
			return nil, 422, errors.New("Price must be > 0")
		}
		updates["price"] = sub.Price
	}

	if sub.UserID != uuid.Nil {
		updates["user_id"] = sub.UserID
	}

	if _, err := time.Parse("2006-01-02", sub.StartDate); err == nil {
		updates["start_date"] = sub.StartDate
	}

	if sub.EndDate != nil {
		updates["end_date"] = sub.EndDate
	}

	if len(updates) == 0 {
		return nil, 422, errors.New("no fields to update")
	}

	update, notFound, err := s.repo.Update(sub.ID, updates)
	if notFound {
		return nil, 404, errors.New("NOT_FOUND")
	}
	if err != nil {
		s.l.Errorf("Error while update subscription. Error: %v", err)
		return nil, 500, errors.New("INTERNAL_SERVER_ERROR")
	}

	return update, 200, nil

}

func (s *Service) DeleteSubscription(id uuid.UUID) (int, error) {
	notFound, err := s.repo.Delete(id)
	if err != nil {
		s.l.Errorf("Error while deleting subscription. Error: %v", err)
		return 500, errors.New("INTERNAL_SERVER_ERROR")
	}
	if notFound {
		return 404, errors.New("NOT_FOUND")
	}

	return 200, nil
}

func (s *Service) GetTotalCost(from, to time.Time, userID uuid.UUID, serviceName string) (*models.Total, int, error) {
	if to.Before(from) {
		return nil, 422, errors.New("to date cannot be earlier than from date")
	}

	total, err := s.repo.GetTotal(from.Format("2006-01-02"), to.Format("2006-01-02"), userID, serviceName)
	if err != nil {
		s.l.Errorf("Error while get total cost. Err: %v", err)
		return nil, 500, errors.New("INTERNAL_SERVER_ERROR")
	}

	return total, 200, nil
}

func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
func minDate(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func diffMonths(from, to time.Time) int {
	return (to.Year()*12 + int(to.Month())) - (from.Year()*12 + int(from.Month())) + 1
}
