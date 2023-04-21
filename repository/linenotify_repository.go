package repository

import (
	"cybergame-api/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func NewLineNotifyRepository(db *gorm.DB) LineNotifyRepository {
	return &repo{db}
}

type LineNotifyRepository interface {
	//linenotify
	GetLineNotify(req model.LinenotifyListRequest) (*model.SuccessWithPagination, error)
	GetLineNotifyById(id int64) (*model.Linenotify, error)
	CreateLineNotify(data model.LinenotifyCreateBody) error
	UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error

	//linenotifygame
	GetLinenotifyGameById(id int64) (*model.LinenotifyGame, error)
	CreateLinenotifyGame(data model.LineNoifyUsergameBody) error
	GetLinenotifyUserGameById(id int64) (*model.LineNoifyUsergame, error)
}

func (r repo) GetLineNotifyById(id int64) (*model.Linenotify, error) {
	var linenotify model.Linenotify

	if err := r.db.Table("line_notify").
		Select("id, start_credit, token, notify_id, status , created_at, updated_at").
		Where("id = ?", id).
		First(&linenotify).
		Error; err != nil {
		return nil, err
	}
	return &linenotify, nil
}

func (r repo) GetLineNotify(req model.LinenotifyListRequest) (*model.SuccessWithPagination, error) {

	var list []model.LinenotifyListResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("line_notify")
	count = count.Select("id")
	if req.Search != "" {
		count = count.Where("id = ?", req.Search)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("line_notify")
		query = query.Select("id,start_credit, token, notify_id, status")
		if req.Search != "" {
			query = query.Where("id = ?", req.Search)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if err = query.
			Limit(req.Limit).
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	if list == nil {
		list = []model.LinenotifyListResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}
func (r repo) CreateLineNotify(data model.LinenotifyCreateBody) error {
	if err := r.db.Table("line_notify").Create(&data).Error; err != nil {
		fmt.Println(data)
		return err
	}

	return nil
}

func (r repo) UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error {
	if err := r.db.Table("line_notify").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetLinenotifyGameById(id int64) (*model.LinenotifyGame, error) {
	var linenotifygame model.LinenotifyGame

	if err := r.db.Table("type_notify").
		Select("id, name, client_id, client_secret, response_type, redirect_uri, scope, state, created_at , updated_at").
		Where("id = ?", id).
		First(&linenotifygame).
		Error; err != nil {
		return nil, err
	}
	return &linenotifygame, nil
}

func (r repo) CreateLinenotifyGame(data model.LineNoifyUsergameBody) error {
	if err := r.db.Table("user_linenotify").Create(&data).Error; err != nil {
		fmt.Println(data)
		return err
	}

	return nil
}

func (r repo) GetLinenotifyUserGameById(id int64) (*model.LineNoifyUsergame, error) {
	var bot model.LineNoifyUsergame

	if err := r.db.Table("user_linenotify").
		Select("id, user_id, type_notify_id, token, created_at, updated_at").
		Where("id = ?", id).
		First(&bot).
		Error; err != nil {
		return nil, err
	}
	return &bot, nil
}
