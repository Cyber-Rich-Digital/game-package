package service

import (
	"context"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type MessageService interface {
	GetMessageByTagId(data model.MessageParam) (*[]model.MessageResponse, error)
	CreateMessage(data model.MessageBody) error
	SetRead(data model.MessageRead) error
}

type messageService struct {
	repo        repository.MessageRepository
	websiteRepo repository.WebsiteRepository
	tagRepo     repository.TagRepository
	deviceRepo  repository.DeviceRepository
	notiRepo    repository.NotiRepository
	firebase    *firebase.App
}

func NewMessageService(
	repo repository.MessageRepository,
	websiteRepo repository.WebsiteRepository,
	tagRepo repository.TagRepository,
	deviceRepo repository.DeviceRepository,
	notiRepo repository.NotiRepository,
	firebase *firebase.App,
) MessageService {
	return &messageService{repo, websiteRepo, tagRepo, deviceRepo, notiRepo, firebase}
}

func (s *messageService) GetMessageByTagId(data model.MessageParam) (*[]model.MessageResponse, error) {

	messages, err := s.repo.GetMessagesByTagId(data.TagId)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	if messages == nil {
		return nil, notFound("Tags not found")
	}

	return messages, nil
}

func (s *messageService) CreateMessage(body model.MessageBody) error {

	website, err := s.websiteRepo.GetWebsiteByApiKey(body.ApiKey)
	if err != nil {

		if err.Error() == "record not found" {
			return notFound("Apikey not found")
		}

		return internalServerError(err.Error())
	}

	tag, err := s.tagRepo.GetTagBynameAndWebId(body.Tag, website.Id)
	if err != nil {
		return internalServerError(err.Error())
	}

	var tagData = model.TagInsert{}
	tagData.WebsiteId = website.Id
	tagData.Name = body.Tag

	var data = model.MessageCreate{}
	data.TagId = tag.Id
	data.Message = body.Message
	data.WebsiteId = website.Id

	var tagId *int

	if tag.Id == 0 {

		var err error

		tagId, err = s.tagRepo.CreateTag(tagData)
		if err != nil {
			return internalServerError(err.Error())
		}

		data.TagId = *tagId
	}

	if err := s.repo.CreateMessage(data); err != nil {
		return internalServerError(err.Error())
	}

	tokens, err := s.deviceRepo.GetTokensbyWebsiteId(website.Id)
	if err != nil {
		return internalServerError(err.Error())
	}

	go func() {

		devices, err := s.deviceRepo.GetDevicesByWebId(website.Id)
		if err != nil {
			log.Println(err)
		}

		for _, device := range *devices {

			noti, err := s.notiRepo.GetNotiByTagIdAndDeviceId(tag.Id, device.Id)
			if err != nil && err.Error() != "record not found" {
				log.Println(err)
			}

			if noti != nil {

				if err := s.notiRepo.UpdateNoti(tag.Id, device.Id); err != nil {
					log.Println(err)
				}

			} else {

				var notiData = model.Notifications{}
				notiData.TagId = data.TagId
				notiData.DeviceId = device.Id

				if err := s.notiRepo.CreateNoti(notiData); err != nil {
					log.Println(err)
				}

			}

		}
	}()

	go func() {
		for _, token := range *tokens {
			sendNoti(s, body, token, website.Id, tag.Id)
		}
	}()

	return nil
}

func sendNoti(s *messageService, body model.MessageBody, token model.Device, websiteId int, tagId int) {

	ctx := context.Background()
	client, err := s.firebase.Messaging(ctx)
	if err != nil {
		log.Println(err)
	}

	badge := 1

	var data = map[string]string{
		"message":   body.Message,
		"websiteId": fmt.Sprintf("%d", websiteId),
		"tagId":     fmt.Sprintf("%d", tagId),
	}

	var customData = map[string]interface{}{
		"message":   body.Message,
		"websiteId": fmt.Sprintf("%d", websiteId),
		"tagId":     fmt.Sprintf("%d", tagId),
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: body.Tag,
			Body:  body.Message,
		},
		Data: data,
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge:      &badge,
					CustomData: customData,
				},
			},
		},
		Token: token.FcmToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
	}

	if response != "" {
		fmt.Println("Successfully sent message:", response)
	}
}

func (s *messageService) SetRead(data model.MessageRead) error {

	if err := s.notiRepo.SetRead(data); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
