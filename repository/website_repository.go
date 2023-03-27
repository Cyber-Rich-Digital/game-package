package repository

import (
	"cyber-api/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func NewWebsiteRepository(db *gorm.DB) WebsiteRepository {
	return &repo{db}
}

type WebsiteRepository interface {
	CheckWebsite(domainName string) (bool, error)
	GetWebsiteByApiKey(apiKey string) (*model.Website, error)
	GetWebsiteByDomain(domainName string) (*model.Website, error)
	GetWebsite(id int) (*model.Website, error)
	GetWebsiteAndTags(id int) (*model.Website, error)
	GetWebsitesByUserIds(id []int) (*[]model.WebsiteList, error)
	GetWebsites(data model.WebsiteQuery) (*model.Pagination, error)
	GetWebsiteTotals(data model.WebsiteDate) (*[]model.WebsiteListResponse, error)
	CreateWebsite(data model.Website) error
	UpdateWebsite(id int, data model.WebsiteBody) error
	DeleteWebsite(id int) error
}

func (r repo) GetWebsiteByDomain(domainName string) (*model.Website, error) {

	var result *model.Website

	if err := r.db.Table("Websites").
		Select("id, user_id").
		Where("domain_name = ?", domainName).
		Where("deleted_at IS NULL").
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Website not found")
	}

	return result, nil
}

func (r repo) GetWebsiteByApiKey(apiKey string) (*model.Website, error) {

	var result *model.Website

	if err := r.db.Table("Websites").
		Select("id").
		Where("api_key = ?", apiKey).
		Where("deleted_at IS NULL").
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Website not found")
	}

	return result, nil
}

func (r repo) CheckWebsite(domainName string) (bool, error) {

	var count int64

	if err := r.db.Table("Websites").
		Select("id").
		Where("domain_name = ?", domainName).
		Where("deleted_at IS NULL").
		Limit(1).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r repo) GetWebsite(id int) (*model.Website, error) {

	var website model.Website
	if err := r.db.Table("Websites").
		Select("id, title, domain_name, api_key, user_id, created_at").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&website).
		Error; err != nil {
		return nil, err
	}

	if website.Id == 0 {
		return nil, errors.New("Website not found")
	}

	return &website, nil
}

func (r repo) GetWebsiteAndTags(id int) (*model.Website, error) {

	var website model.Website
	if err := r.db.Table("Websites w").
		Select("w.id, w.title, w.domain_name, w.api_key, w.user_id, w.created_at").
		Joins("LEFT JOIN Tags AS t ON t.website_id = w.id AND t.deleted_at IS NULL").
		Where("w.id = ?", id).
		Where("w.deleted_at IS NULL").
		First(&website).
		Error; err != nil {
		return nil, err
	}

	if website.Id == 0 {
		return nil, errors.New("Website not found")
	}

	if err := r.db.Table("Tags").
		Select("id, name, website_id").
		Where("website_id = ?", id).
		Where("deleted_at IS NULL").
		Find(&website.Tags).
		Error; err != nil {
		return nil, err
	}

	return &website, nil
}

func (r repo) GetWebsitesByUserIds(id []int) (*[]model.WebsiteList, error) {

	var website *[]model.WebsiteList
	if err := r.db.Table("Websites").
		Select("id, title, domain_name, api_key, user_id, created_at").
		Where("id IN (?)", id).
		Where("deleted_at IS NULL").
		Find(&website).
		Error; err != nil {
		return nil, err
	}

	return website, nil
}

func (r repo) GetWebsites(data model.WebsiteQuery) (*model.Pagination, error) {

	var list []model.WebsiteResponse
	var total int64
	var err error

	selectFields := "w.id, w.title, w.domain_name, w.api_key, w.created_at, w.updated_at, COUNT(m.id) AS total"
	join := "LEFT OUTER JOIN Messages AS m ON w.id = m.website_id"
	group := "m.website_id, w.id, w.title, w.domain_name, w.api_key, w.created_at, w.updated_at"
	whereVal := fmt.Sprintf("%%%s%%", data.Search)

	// Get list of websites //

	query := r.db.Table("Websites w")
	query = query.Select(selectFields).
		Joins(join).
		Group(group)

	if data.Role == "USER" {
		query = query.Where("w.user_id = ?", data.UserId)
	}

	if data.Search != "" {
		query = query.Where("w.title LIKE ?", whereVal).
			Or("w.domain_name LIKE ?", whereVal)
	}

	// Sort by created_at //

	if data.Sort == 1 {
		query = query.Order("w.created_at ASC")
	} else {
		query = query.Order("w.created_at DESC")
	}

	if err = query.
		Where("w.deleted_at IS NULL").
		Limit(data.Limit).
		Offset(data.Page * data.Limit).
		Scan(&list).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset)  //

	count := r.db.Table("Websites w")
	count = count.Select(selectFields).
		Joins(join).
		Group(group)

	if data.Role == "USER" {
		count = count.Where("w.user_id = ?", data.UserId)
	}

	if data.Search != "" {
		count = count.Where("w.title LIKE ?", whereVal).
			Or("w.domain_name LIKE ?", whereVal)
	}

	if err = count.
		Where("w.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	// End count total records for pagination purposes (without limit and offset)  //

	var website model.Pagination
	website.List = list
	website.Total = total

	return &website, nil
}

func (r repo) GetWebsiteTotals(data model.WebsiteDate) (*[]model.WebsiteListResponse, error) {

	t := time.Unix(data.Date, 0)

	var list *[]model.WebsiteListResponse

	selectFields := "w.id, COUNT(m.id) AS total"
	join := "LEFT OUTER JOIN Messages AS m ON  w.id = m.website_id AND m.created_at >= ?"
	group := "m.website_id, w.id"

	if err := r.db.Table("Websites w").
		Select(selectFields).
		Joins(join, t.Format("2006-01-02 15:04:05")).
		Group(group).
		Where("w.user_id = ?", data.UserId).
		Where("w.deleted_at IS NULL").
		Find(&list).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r repo) CreateWebsite(data model.Website) error {

	if err := r.db.Table("Websites").Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateWebsite(id int, data model.WebsiteBody) error {

	if err := r.db.Table("Websites").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteWebsite(id int) error {

	if err := r.db.Table("Websites").Where("id = ?", id).Delete(&model.Website{}).Error; err != nil {
		return err
	}

	return nil
}
