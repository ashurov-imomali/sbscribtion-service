package repository

import (
	"fmt"
	"github.com/ashurov-imomali/sbscribtion-service/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *repo) GetByID(id uuid.UUID) (*models.Subscription, bool, error) {
	var result models.Subscription
	if err := r.db.First(&result, id).Error; err != nil {
		return nil, errors.Is(err, gorm.ErrRecordNotFound), err
	}
	return &result, false, nil
}

func (r *repo) GetByFilter(filter models.SubscriptionFilter) ([]models.Subscription, bool, error) {
	var subs []models.Subscription
	db := r.db.Model(&models.Subscription{})
	if filter.UserID != nil {
		db = db.Where("user_id = ?", *filter.UserID)
	}
	if filter.ServiceName != nil {
		db = db.Where("service_name ilike ?", fmt.Sprintf("%%%s%%", *filter.ServiceName))
	}
	if filter.StartDate != nil {
		db = db.Where("start_date::date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		db = db.Where("end_date::date <= ?", *filter.EndDate)
	}

	tx := db.Find(&subs)
	if tx.Error != nil {
		return nil, false, tx.Error
	}

	return subs, tx.RowsAffected == 0, nil
}

func (r *repo) Update(id uuid.UUID, updates map[string]interface{}) (*models.Subscription, bool, error) {
	var result models.Subscription
	tx := r.db.Model(&result).
		Clauses(clause.Returning{}).
		Where("id=?", id).Updates(updates)
	if tx.Error != nil {
		return nil, false, tx.Error
	}
	return &result, tx.RowsAffected == 0, nil
}

func (r *repo) Delete(id uuid.UUID) (bool, error) {
	tx := r.db.Delete(models.Subscription{}, id)
	return tx.RowsAffected == 0, tx.Error
}

func (r *repo) GetTotal(from, to string, userID uuid.UUID, service string) (*models.Total, error) {
	var result models.Total

	tx := r.db.Table("subscriptions").Select(`
		sum(
            price * (
                date_part('year', age(least(?::date, end_date::date), greatest(?::date, start_date::date))) * 12 +
                date_part('month', age(least(?::date, end_date::date), greatest(?::date, start_date::date)))
            )
        ) as total_cost`, to, from, to, from).
		Where("start_date::date < ? and end_date::date >= ?", to, from)

	if userID != uuid.Nil {
		tx = tx.Where("user_id = ?", userID)
	}

	if service != "" {
		tx = tx.Where("service_name = ?", service)
	}

	return &result, tx.Scan(&result.TotalCost).Error
}

//select id, start_date, end_date,
//	date_part('year',age(least('2024-05-01', end_date), greatest('2024-02-01', start_date))) * 12+
//		date_part('month',age(least('2024-05-01', end_date), greatest('2024-02-01', start_date))) as m,
//	price * (date_part('year',age(least('2024-05-01', end_date), greatest('2024-02-01', start_date))) * 12+
//		date_part('month',age(least('2024-05-01', end_date), greatest('2024-02-01', start_date)))) total
//	from subscriptions
//	where start_date < '2024-05-01' and end_date >= '2024-02-01';
