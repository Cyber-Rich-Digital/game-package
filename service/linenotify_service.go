package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"os"

	"github.com/juunini/simple-go-line-notify/notify"
)

func NewLineNotifyService(
	repo repository.LineNotifyRepository,
) LineNotifyService {
	return &lineNotifyService{repo}
}

type LineNotifyService interface {
	//line notify group
	CreateLineNotify(data model.LinenotifyCreateBody) error
	GetLineNotify(data model.LinenotifyListRequest) (*model.SuccessWithPagination, error)
	GetLineNotifyById(data model.LinenotifyParam) (*model.Linenotify, error)
	UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error

	//line notify game
	GetLineNotifyGameById(model.LinenotifyGameParam) (*model.LinenotifyGame, error)
	CreateNotifyGame(data model.LineNoifyUsergameBody) error
	GetLineNoifyUserGameById(model.LineNotifyUserGameParam) (*model.LineNoifyUsergame, error)
}

type lineNotifyService struct {
	repo repository.LineNotifyRepository
}

// CreateSettingWeb implements SettingWebService
func (s *lineNotifyService) CreateLineNotify(data model.LinenotifyCreateBody) error {

	var line model.Linenotify
	line.StartCredit = data.StartCredit
	line.Token = data.Token
	line.NotifyId = data.NotifyId
	line.Status = data.Status

	accessToken := data.Token
	message := os.Getenv("MESSAGE_LINENOTIFY")

	if err := notify.SendText(accessToken, message); err != nil {
		panic(err)
	}

	if err := s.repo.CreateLineNotify(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *lineNotifyService) GetLineNotify(params model.LinenotifyListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&params.Page, &params.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetLineNotify(params)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *lineNotifyService) GetLineNotifyById(data model.LinenotifyParam) (*model.Linenotify, error) {

	line, err := s.repo.GetLineNotifyById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		return nil, internalServerError(err.Error())
	}
	return line, nil
}

func (s *lineNotifyService) UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error {
	if err := s.repo.UpdateLineNotify(id, data); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *lineNotifyService) GetLineNotifyGameById(data model.LinenotifyGameParam) (*model.LinenotifyGame, error) {

	linegame, err := s.repo.GetLinenotifyGameById(data.Id)

	if err != nil {
		return nil, notFound("record NotFound")
	}

	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		return nil, internalServerError(err.Error())
	}
	return linegame, nil
}

func (s *lineNotifyService) CreateNotifyGame(data model.LineNoifyUsergameBody) error {

	//reqURL := "https://notify-bot.line.me/oauth/authorize?response_type=code&client_id=VGtxc8xQ2ghbxd71y6zSB3&redirect_uri=https://cyberrichdigital.com/&scope=notify&state=testdata"
	var bot model.LinenotifyGame
	bot.ResponseType = "code"
	bot.ClientId = "VGtxc8xQ2ghbxd71y6zSB3"
	bot.RedirectUri = "https://cyberrichdigital.com/"
	bot.Scope = "notify"
	bot.State = "1"
	if err := s.repo.CreateLinenotifyGame(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *lineNotifyService) GetLineNoifyUserGameById(data model.LineNotifyUserGameParam) (*model.LineNoifyUsergame, error) {

	botuser, err := s.repo.GetLinenotifyUserGameById(data.Id)

	if err != nil {
		return nil, notFound("record NotFound")
	}

	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		return nil, internalServerError(err.Error())
	}
	return botuser, nil
}
