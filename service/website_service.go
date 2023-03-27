package service

import (
	"cyber-api/helper"
	"cyber-api/model"
	"cyber-api/repository"
)

type WebsiteService interface {
	GetWebsiteAndTags(data model.WebsiteParam, userId int) (*model.Website, error)
	GetWebsite(data model.WebsiteParam, userId int) (*model.Website, error)
	GetWebsites(data model.WebsiteQuery) (*model.Pagination, error)
	GetWebsiteTotals(body model.WebsiteDate) (*[]model.WebsiteListResponse, error)
	CreateWebsite(data model.WebsiteBody, userId int) error
	UpdateWebsite(id int, data model.WebsiteBody) error
	DeleteWebsite(id int) error
}

type websiteService struct {
	repo repository.WebsiteRepository
}

var valueNotFound = "Website not found"

func NewWebsiteService(
	repo repository.WebsiteRepository,
) WebsiteService {
	return &websiteService{repo}
}

func (s *websiteService) GetWebsiteAndTags(data model.WebsiteParam, userId int) (*model.Website, error) {

	website, err := s.repo.GetWebsiteAndTags(data.Id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(valueNotFound)
		}

		if err.Error() == "Website not found" {
			return nil, notFound(valueNotFound)
		}

		return nil, internalServerError(err.Error())
	}

	// if website.UserId != userId {
	// 	return nil, badRequest("You don't have permission to access this website")
	// }

	return website, nil
}

func (s *websiteService) GetWebsite(data model.WebsiteParam, userId int) (*model.Website, error) {

	website, err := s.repo.GetWebsite(data.Id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(valueNotFound)
		}

		if err.Error() == "Website not found" {
			return nil, notFound(valueNotFound)
		}

		return nil, internalServerError(err.Error())
	}

	// if website.UserId != userId {
	// 	return nil, badRequest("You don't have permission to access this website")
	// }

	return website, nil
}

func (s *websiteService) GetWebsiteTotals(body model.WebsiteDate) (*[]model.WebsiteListResponse, error) {

	website, err := s.repo.GetWebsiteTotals(body)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(valueNotFound)
		}

		if err.Error() == "Website not found" {
			return nil, notFound(valueNotFound)
		}

		return nil, internalServerError(err.Error())
	}

	return website, nil
}

func (s *websiteService) GetWebsites(data model.WebsiteQuery) (*model.Pagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	website, err := s.repo.GetWebsites(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return website, nil
}

func (s *websiteService) CreateWebsite(data model.WebsiteBody, userId int) error {

	exist, err := s.repo.CheckWebsite(data.DomainName)
	if err != nil {
		return internalServerError(err.Error())
	}

	if exist {
		return badRequest("Website already exist")
	}

	key := helper.GenKey(50)

	var website model.Website
	website.Title = data.Title
	website.DomainName = data.DomainName
	website.UserId = userId
	website.ApiKey = key

	if err := s.repo.CreateWebsite(website); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *websiteService) UpdateWebsite(id int, data model.WebsiteBody) error {

	// if check, err := s.repo.GetWebsitesByIdAndUserIds(id); err != nil {
	// 	return internalServerError(err.Error())
	// } else if !check {
	// 	return notFound("Website not found")
	// }
	check, err := s.repo.CheckWebsite(data.DomainName)
	if err != nil {
		return internalServerError(err.Error())
	}

	if check {
		return badRequest("Website already exist")
	}

	if err := s.repo.UpdateWebsite(id, data); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *websiteService) DeleteWebsite(id int) error {

	if err := s.repo.DeleteWebsite(id); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
