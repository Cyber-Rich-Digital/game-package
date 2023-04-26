package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
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
	DeleteLineNoifyUserGame(id int64) error
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

	var lineNotifyGame model.LinenotifyGame
	var lineGame model.LinenotifyGameParam

	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("URL_LINE"), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	// Set the query parameters
	q := req.URL.Query()
	q.Add("response_type", lineNotifyGame.ResponseType)
	q.Add("client_id", lineNotifyGame.ClientId)
	q.Add("redirect_uri", lineNotifyGame.RedirectUri)
	q.Add("scope", lineNotifyGame.Scope)
	q.Add("state", "1")
	req.URL.RawQuery = q.Encode()
	fmt.Println("url request:", req.URL.RawQuery)

	// Set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	// Print the response status code and body
	defer res.Body.Close()
	fmt.Println("Response status code:", res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	fmt.Println("Response body:", body)

	response, err := client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {

		res, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var result interface{}
		fmt.Println("Response body:", body)

		if err := json.Unmarshal(res, &result); err != nil {
			return nil, err
		}

		return nil, errors.New("Error some")

	}

	if err != nil {
		return nil, notFound("record NotFound")
	}

	URL_AUTH_LINE := os.Getenv("URL_LINE") + "/oauth/token?grant_type=authorization_code&code=" + lineGame.Code + "&redirect_uri=" + lineNotifyGame.RedirectUri + "&client_id=" + lineNotifyGame.ClientId + "&client_secret=" + lineNotifyGame.ClientSecret

	url := URL_AUTH_LINE
	req, err1 := http.NewRequest("POST", url, nil)
	if err != nil {
		sentry.CaptureException(err1)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	response, err2 := client.Do(req)
	if err != nil {
		sentry.CaptureException(err2)
		return nil, err
	}
	defer response.Body.Close()
	fmt.Println("Response status code:", res.StatusCode)
	responsebody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	fmt.Println("Response body:", string(responsebody))

	return nil, internalServerError(err.Error())

}

func (s *lineNotifyService) CreateNotifyGame(data model.LineNoifyUsergameBody) error {

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

func (s *lineNotifyService) DeleteLineNoifyUserGame(id int64) error {

	_, err := s.repo.GetLinenotifyUserGameById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repo.DeleteLinenotifyGame(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
