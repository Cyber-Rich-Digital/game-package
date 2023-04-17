package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewRecommentRepository(db *gorm.DB) RecommentRepository {
	return &repo{db}
}

type RecommentRepository interface {
	GetRecommentList(query model.RecommentQuery) ([]model.RecommentList, int64, error)
	CreateRecomment(recomment model.CreateRecomment) error
	UpdateRecomment(id int64, body model.CreateRecomment) error
}

func (r repo) GetRecommentList(query model.RecommentQuery) ([]model.RecommentList, int64, error) {

	var recomments []model.RecommentList

	db := r.db.Table("Recommend_channels").Select("id, title, status, url, created_at")

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Filter != "" {
		db = db.Where("title LIKE ?", "%"+query.Filter+"%")
	}

	if err := db.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&recomments).
		Order("id ASC").
		Error; err != nil {
		return nil, 0, err
	}

	var total int64

	if err := r.db.Table("Recommend_channels").
		Select("id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	return recomments, total, nil
}

func (r repo) CreateRecomment(recomment model.CreateRecomment) error {

	if err := r.db.Table("Recommend_channels").
		Create(&recomment).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateRecomment(id int64, body model.CreateRecomment) error {

	if err := r.db.Table("Recommend_channels").
		Where("id = ?", id).
		Updates(&body).
		Error; err != nil {
		return err
	}

	return nil
}
