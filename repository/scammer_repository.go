package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewScammerRepository(db *gorm.DB) ScammerRepository {
	return &repo{db}
}

type ScammerRepository interface {
	GetScammerList(query model.ScammerQuery) ([]model.ScammertList, error)
	CreateScammer(scammer model.Scammer) error
}

func (r repo) GetScammerList(query model.ScammerQuery) ([]model.ScammertList, error) {

	var scammers []model.ScammertList

	db := r.db.Table("Scammers")

	if query.DateStart != "" && query.DateEnd != "" {
		db = db.Where("created_at BETWEEN ? AND ?", query.DateStart, query.DateEnd)
	}

	if query.BankName != nil {
		db = db.Where("bankname = ?", query.BankName)
	}

	if query.Filter != nil {
		db = db.Where("fullname LIKE ? OR bankname LIKE ? OR bank_account LIKE ?", "%"+*query.Filter+"%", "%"+*query.Filter+"%", "%"+*query.Filter+"%")
	}

	if err := db.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&scammers).
		Order("id desc").
		Error; err != nil {
		return nil, err
	}

	return scammers, nil
}

func (r repo) CreateScammer(scammer model.Scammer) error {

	if err := r.db.Table("Scammers").
		Create(&scammer).
		Error; err != nil {
		return err
	}

	return nil
}
