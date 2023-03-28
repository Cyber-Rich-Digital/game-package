package repository

import (
	"cybergame-api/model"
	"fmt"

	"gorm.io/gorm"
)

func NewTagRepository(db *gorm.DB) TagRepository {
	return &repo{db}
}

type TagRepository interface {
	GetTagsByWebsiteId(data model.TagParam) (*[]model.TagResponse, error)
	GetTagBynameAndWebId(name string, webId int) (*model.Tag, error)
	CreateTag(tag model.TagInsert) (*int, error)
	DeleteTag(id int) error
}

func (r repo) GetTagsByWebsiteId(data model.TagParam) (*[]model.TagResponse, error) {
	fmt.Println(data.DeviceId)

	selectFields := "t.id, t.name, t.created_at, n.total"
	joinNoti := "LEFT JOIN Notifications AS n ON n.device_id = d.id AND n.tag_id = t.id"
	joinDevice := "LEFT JOIN Devices AS d ON d.id = ? AND d.website_id = t.website_id"

	var tags []model.TagResponse
	if err := r.db.Table("Tags t").
		Distinct("t.id").
		Select(selectFields).
		Joins(joinDevice, data.DeviceId).
		Joins(joinNoti).
		Group("t.id, n.total, d.id").
		Where("t.website_id = ?", data.WebsiteId).
		Where("t.deleted_at IS NULL").
		Find(&tags).
		Error; err != nil {
		return nil, err
	}

	return &tags, nil
}

func (r repo) GetTagBynameAndWebId(name string, webId int) (*model.Tag, error) {

	var tags *model.Tag
	if err := r.db.Table("Tags").
		Select("id, website_id, name, created_at").
		Where("name = ? AND website_id = ?", name, webId).
		Where("deleted_at IS NULL").
		Find(&tags).
		Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r repo) CreateTag(tag model.TagInsert) (*int, error) {

	result := r.db.Table("Tags").Create(&tag)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tag.Id, nil
}

func (r repo) DeleteTag(id int) error {

	result := r.db.Table("Tags").Where("id = ?", id).Delete(&model.Tag{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
