package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"strings"
)

type AdminService interface {
	GetGroup(id int) (*model.AdminGroupPermissionResponse, error)
	GetGroupList() (*[]model.GroupList, error)
	Login(data model.LoginAdmin) (*string, error)
	Create(user *model.CreateAdmin) error
	CreateGroup(data *model.AdminCreateGroup) error
	DeleteGroup(id int64) error
	DeletePermission(id int64) error
}

const AdminloginFailed = "Phone Or Password is incorrect"
const AdminNotFound = "Admin not found"
const AdminGroupNotFound = "Group not found"

type adminService struct {
	repo      repository.AdminRepository
	perRepo   repository.PermissionRepository
	groupRepo repository.GroupRepository
}

func NewAdminService(
	repo repository.AdminRepository,
	perRepo repository.PermissionRepository,
	groupRepo repository.GroupRepository,
) AdminService {
	return &adminService{repo, perRepo, groupRepo}
}

func (s *adminService) GetGroup(id int) (*model.AdminGroupPermissionResponse, error) {

	group, err := s.repo.GetGroup(id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminGroupNotFound)
		}

		return nil, err
	}

	return group, nil
}

func (s *adminService) GetGroupList() (*[]model.GroupList, error) {

	group, err := s.repo.GetGroupList()
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *adminService) Login(data model.LoginAdmin) (*string, error) {

	user, err := s.repo.GetAdminByUsername(data)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminloginFailed)
		}

		return nil, internalServerError(err.Error())
	}

	if user == nil {
		return nil, badRequest(AdminloginFailed)
	}

	if err := helper.ComparePassword(data.Password, user.Password); err != nil {
		return nil, badRequest(AdminloginFailed)
	}

	token, err := helper.CreateJWT(*user)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}

func (s *adminService) Create(data *model.CreateAdmin) error {

	username, err := s.repo.CheckUsername(data.Username)
	if err != nil {
		return err
	}

	if username {
		return badRequest("Username already exist")
	}

	email, err := s.repo.CheckPhone(data.Phone)
	if err != nil {
		return err
	}

	if email {
		return badRequest("Phone already exist")
	}

	hashedPassword, err := helper.GenPassword(data.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	newUser := model.Admin{}
	newUser.Email = data.Email
	newUser.Username = data.Username
	newUser.Password = string(hashedPassword)
	newUser.Role = "ADMIN"

	return s.repo.CreateAdmin(newUser)
}

func (s *adminService) CreateGroup(data *model.AdminCreateGroup) error {

	checkGroup, err := s.groupRepo.CheckGroupExist(data.GroupId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound)
	}

	var groupIds []int64
	for _, v := range data.PermissionIds {
		groupIds = append(groupIds, int64(v))
	}

	checkPermission, err := s.perRepo.CheckPerListExist(groupIds)
	if err != nil {
		return internalServerError(err.Error())
	}

	var idNotFound []string
	for _, j := range groupIds {

		exist := false

		for _, k := range checkPermission {
			if j == k {
				exist = true
			}
		}

		if !exist {
			idNotFound = append(idNotFound, fmt.Sprintf("%d", j))
		}
	}

	if len(idNotFound) > 0 {
		return badRequest(fmt.Sprintf("Permission id %s not found", strings.Join(idNotFound, ",")))
	}

	var list []model.AdminPermissionList

	for _, v := range data.PermissionIds {
		list = append(list, model.AdminPermissionList{
			GroupId:      data.GroupId,
			PermissionId: v,
		})
	}

	if err := s.repo.CreateGroup(list); err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeleteGroup(id int64) error {

	if err := s.groupRepo.DeleteGroup(id); err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeletePermission(id int64) error {

	if err := s.perRepo.DeletePermission(id); err != nil {
		return err
	}

	return nil
}
