package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"strings"
)

type AdminService interface {
	GetAdmin(id int64) (*model.AdminDetail, error)
	GetAdminList(query model.AdminListQuery) (*model.SuccessWithPagination, error)
	GetGroup(id int) (*model.AdminGroupPermissionResponse, error)
	GetGroupList(query model.AdminGroupQuery) (*model.SuccessWithPagination, error)
	Login(data model.LoginAdmin) (*string, error)
	Create(user *model.CreateAdmin) (error, *[]string)
	CreateGroup(data *model.AdminCreateGroup) error
	UpdateAdmin(adminId int64, data model.AdminBody) (error, *[]string)
	UpdateGroup(groupId int64, data *model.AdminUpdateGroup) error
	ResetPassword(adminId int64, body model.AdminUpdatePassword) error
	DeleteGroup(id int64) error
	DeletePermission(perm model.DeletePermission) error
	DeleteAdmin(id int64) error
}

const AdminloginFailed = "Phone Or Password is incorrect"
const AdminNotFound = "Admin not found"
const AdminExist = "Admin already exist"
const AdminPhoneExist = "Phone already exist"
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

func (s *adminService) GetAdmin(id int64) (*model.AdminDetail, error) {

	admin, perList, group, err := s.repo.GetAdmin(id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminNotFound)
		}

		return nil, err
	}

	var result model.AdminDetail
	result.Id = admin.Id
	result.Username = admin.Username
	result.Fullname = admin.Fullname
	result.Phone = admin.Phone
	result.Email = admin.Email
	result.Status = admin.Status
	result.Role = admin.Role
	result.PermissionList = *perList

	if group != nil {
		result.Group = group
	}

	return &result, nil
}

func (s *adminService) GetAdminList(query model.AdminListQuery) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, err
	}

	list, total, err := s.repo.GetAdminList(query)
	if err != nil {
		return nil, err
	}

	result := &model.SuccessWithPagination{
		Message: "Success",
		List:    list,
		Total:   *total,
	}

	return result, nil
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

func (s *adminService) GetGroupList(query model.AdminGroupQuery) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, err
	}

	list, total, err := s.repo.GetGroupList(query)
	if err != nil {
		return nil, err
	}

	result := &model.SuccessWithPagination{
		Message: "Success",
		List:    list,
		Total:   *total,
	}

	return result, nil
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

	if err := helper.CompareAdminPassword(data.Password, user.Password); err != nil {
		return nil, badRequest(AdminloginFailed)
	}

	token, err := helper.CreateJWTAdmin(*user)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}

func (s *adminService) Create(data *model.CreateAdmin) (error, *[]string) {

	username, err := s.repo.CheckAdmin(data.Username)
	if err != nil {
		return err, nil
	}

	if username {
		return badRequest("Username already exist"), nil
	}

	email, err := s.repo.CheckPhone(data.Phone)
	if err != nil {
		return err, nil
	}

	if email {
		return badRequest("Phone already exist"), nil
	}

	checkGroup, err := s.groupRepo.CheckGroupExist(data.AdminGroupId)
	if err != nil {
		return internalServerError(err.Error()), nil
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound), nil
	}

	var perIds []int64
	for _, j := range *data.Permissions {
		perIds = append(perIds, j.Id)
	}

	checkPermission, err := s.perRepo.CheckPerListAndGroupId(data.AdminGroupId, perIds)
	if err != nil {
		return internalServerError(err.Error()), nil
	}

	var idNotFound []string
	for _, j := range *data.Permissions {

		exist := false

		for _, k := range checkPermission {
			if j.Id == k.Id {
				exist = true
			}
		}

		if !exist {
			idNotFound = append(idNotFound, fmt.Sprintf("%d", j.Id))
		}
	}

	if len(idNotFound) > 0 {
		return badRequest(fmt.Sprintf("Permission id %s not found", strings.Join(idNotFound, ","))), nil
	}

	var notToggleList []string
	for _, j := range *data.Permissions {

		for _, k := range checkPermission {
			if j.Id == k.Id {

				if !k.IsRead && j.IsRead {
					notToggleList = append(notToggleList, fmt.Sprintf("%s View ไม่ได้รับอนุญาตให้เปิดใช้งาน.", k.Name))
				}

				if !k.IsWrite && j.IsWrite {
					notToggleList = append(notToggleList, fmt.Sprintf("%s Manage ไม่ได้รับอนุญาตให้เปิดใช้งาน.", k.Name))
				}
			}
		}
	}

	if len(notToggleList) > 0 {
		return nil, &notToggleList
	}

	hashedPassword, err := helper.GenAdminPassword(data.Password)
	if err != nil {
		return internalServerError(err.Error()), nil
	}

	splitFullname := strings.Split(data.Fullname, " ")
	if len(splitFullname) == 1 || strings.Trim(data.Fullname, " ") == "" {
		return badRequest("Fullname must be firstname lastname"), nil
	}

	var firstname, lastname *string
	if len(splitFullname) == 2 {
		firstname = &splitFullname[0]
		lastname = &splitFullname[1]
	}

	if len(splitFullname) == 3 {
		firstname = &splitFullname[1]
		lastname = &splitFullname[2]
	}

	newUser := model.Admin{}
	newUser.Email = data.Email
	newUser.Username = data.Username
	newUser.Fullname = data.Fullname
	newUser.Firstname = *firstname
	newUser.Lastname = *lastname
	newUser.Password = string(hashedPassword)
	newUser.Role = "ADMIN"
	newUser.Status = data.Status
	newUser.Phone = data.Phone
	newUser.AdminGroupId = data.AdminGroupId

	return s.repo.CreateAdmin(newUser, data.Permissions), nil
}

func (s *adminService) CreateGroup(data *model.AdminCreateGroup) error {

	checkGroup, err := s.groupRepo.CheckGroupExist(data.GroupId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound)
	}

	var perIds []int64
	for _, v := range data.Permissions {
		perIds = append(perIds, int64(v.Id))
	}

	checkPermission, err := s.perRepo.CheckPerListExist(perIds)
	if err != nil {
		return internalServerError(err.Error())
	}

	var idNotFound []string
	for _, j := range perIds {

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

	for _, v := range data.Permissions {
		list = append(list, model.AdminPermissionList{
			GroupId:      data.GroupId,
			PermissionId: v.Id,
			IsRead:       v.IsRead,
			IsWrite:      v.IsWrite,
		})
	}

	if err := s.repo.CreateGroupAdmin(list); err != nil {
		return err
	}

	return nil
}

func (s *adminService) UpdateAdmin(adminId int64, body model.AdminBody) (error, *[]string) {

	var data model.UpdateAdmin

	if body.GroupId != nil {
		checkGroup, err := s.groupRepo.CheckGroupExist(*body.GroupId)
		if err != nil {
			return internalServerError(err.Error()), nil
		}

		if !checkGroup {
			return notFound(AdminGroupNotFound), nil
		}

		data.AdminGroupId = body.GroupId
	}

	var adminPer []model.AdminPermission
	var oldGroupId *int

	if body.GroupId != nil && body.Permissions != nil {

		getGroupId, err := s.repo.GetAdminGroup(adminId)
		if err != nil {
			return internalServerError(err.Error()), nil
		}

		oldGroupId = &getGroupId.AdminGroupId
		var perIds []int64
		for _, j := range *body.Permissions {
			perIds = append(perIds, j.Id)
		}

		checkPermission, err := s.perRepo.CheckPerListAndGroupId(*body.GroupId, perIds)
		if err != nil {
			return internalServerError(err.Error()), nil
		}

		var idNotFound []string
		for _, j := range *body.Permissions {

			exist := false

			for _, k := range checkPermission {
				if j.Id == k.Id {
					exist = true
				}
			}

			if !exist {
				idNotFound = append(idNotFound, fmt.Sprintf("%v", j.Id))
			}
		}

		if len(idNotFound) > 0 {
			return badRequest(fmt.Sprintf("Permission id %s not found", strings.Join(idNotFound, ","))), nil
		}

		var notToggleList []string
		for _, j := range *body.Permissions {

			for _, k := range checkPermission {
				if j.Id == k.Id {

					if !k.IsRead && j.IsRead {
						notToggleList = append(notToggleList, fmt.Sprintf("%s View ไม่ได้รับอนุญาตให้เปิดใช้งาน.", k.Name))
					}

					if !k.IsWrite && j.IsWrite {
						notToggleList = append(notToggleList, fmt.Sprintf("%s Manage ไม่ได้รับอนุญาตให้เปิดใช้งาน.", k.Name))
					}
				}
			}
		}

		if len(notToggleList) > 0 {
			return nil, &notToggleList
		}

		for _, v := range *body.Permissions {
			adminPer = append(adminPer, model.AdminPermission{
				AdminId:      adminId,
				PermissionId: v.Id,
				IsRead:       v.IsRead,
				IsWrite:      v.IsWrite,
			})
		}
	}

	data.Email = body.Email
	data.Status = body.Status

	splitFullname := strings.Split(body.Fullname, " ")
	if len(splitFullname) == 1 || strings.Trim(body.Fullname, " ") == "" {
		return badRequest("Fullname must be firstname lastname"), nil
	}

	if len(splitFullname) == 2 {
		data.Firstname = splitFullname[0]
		data.Lastname = splitFullname[1]
	}

	if len(splitFullname) == 3 {
		data.Firstname = splitFullname[1]
		data.Lastname = splitFullname[2]
	}

	data.Fullname = body.Fullname

	return s.repo.UpdateAdmin(adminId, oldGroupId, data, &adminPer), nil
}

func (s *adminService) UpdateGroup(groupId int64, data *model.AdminUpdateGroup) error {

	checkGroup, err := s.groupRepo.CheckGroupExist(groupId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound)
	}

	var permissionIds []int64
	for _, v := range data.Permissions {
		permissionIds = append(permissionIds, int64(v.Id))
	}

	checkPermission, err := s.perRepo.CheckPerListExist(permissionIds)
	if err != nil {
		return internalServerError(err.Error())
	}

	var idNotFound []string
	for _, j := range permissionIds {

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

	for _, v := range data.Permissions {
		list = append(list, model.AdminPermissionList{
			GroupId:      groupId,
			PermissionId: v.Id,
			IsRead:       v.IsRead,
			IsWrite:      v.IsWrite,
		})
	}

	if err := s.repo.UpdateGroup(groupId, data.Name, list); err != nil {
		return err
	}

	return nil
}

func (s *adminService) ResetPassword(adminId int64, body model.AdminUpdatePassword) error {

	checkAdmin, err := s.repo.CheckAdminById(adminId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkAdmin {
		return notFound(AdminNotFound)
	}

	newPasword, err := helper.GenAdminPassword(body.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	body.Password = newPasword

	if err := s.repo.UpdatePassword(adminId, body); err != nil {
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

func (s *adminService) DeletePermission(perm model.DeletePermission) error {

	if err := s.perRepo.DeletePermission(perm); err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeleteAdmin(id int64) error {

	if err := s.repo.DeleteAdmin(id); err != nil {
		return err
	}

	return nil
}
