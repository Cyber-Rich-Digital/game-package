package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type menuController struct {
	menuService service.MenuService
}

func newMenuController(
	menuService service.MenuService,
) menuController {
	return menuController{menuService}
}

func MenuController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewPermissionRepository(db)
	service := service.NewMenuService(repo)
	handler := newMenuController(service)

	r = r.Group("/menu")
	r.GET("/", handler.GetMenu)

}

// @Summary Get Menu
// @Description Get Menu
// @Tags Menu
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /menu [get]
func (h menuController) GetMenu(c *gin.Context) {

	list, err := h.menuService.GetMenu()
	if err != nil {
		HandleError(c, err)
		return
	}

	titles := []map[string]interface{}{
		{"title": "คู่มือ", "name": "guide"},

		{"title": "ผู้ดูแล", "name": "admin"},
		{"title": "จัดการผู้ใช้งาน", "name": "admin"},
		{"title": "กลุ่มผู้ใช้งาน", "name": "admin_group"},

		{"title": "สรุปภาพรวม", "name": "summary"},

		{"title": "จัดการเว็บเอเย่น", "name": "agent"},
		{"title": "รายการเว็บเอเย่น", "name": "agent_list"},
		{"title": "รายงานเพิ่ม-ลด เครดิต", "name": "agent_credit"},

		{"title": "จัดการธนาคาร", "name": "bank"},
		{"title": "รายการธนาคาร", "name": "bank_list"},
		{"title": "รายงานธุรกรรมเงินสด", "name": "bank_transaction"},
		{"title": "รายการเดินบัญชีธนาคาร", "name": "bank_account"},
		{"title": "จัดการโปรโมชั่น", "name": "promotion"},

		{"title": "จัดการการตลาด", "name": "marketing"},
		{"title": "รายการลิ้งรับทรัพย์", "name": "marketing_link"},
		{"title": "รายการพันธมิตร", "name": "marketing_partner"},

		{"title": "จัดการกิจกรรม", "name": "activity"},
		{"title": "คืนยอดเสีย", "name": "activity_return"},
		{"title": "กงล้อนำโชค", "name": "activity_lucky"},
		{"title": "เช็คอินรายวัน", "name": "activity_checkin"},
		{"title": "คูปองเงินสด", "name": "activity_coupon"},

		{"title": "จัดการสมาชิกเว็บ", "name": "member"},
		{"title": "รายการสมาชิกเว็บ", "name": "member_list"},
		{"title": "ประวัติฝาก-ถอนสมาชิก", "name": "member_transaction"},
		{"title": "ตั้งค่าช่องทางที่รู้จัก", "name": "member_channel"},
		{"title": "ประวัติการแก้ไขข้อมูล", "name": "member_history"},
		{"title": "รายการมิจฉาชีพ", "name": "member_misconduct"},

		{"title": "รายงาน", "name": "report"},
		{"title": "ยอดสมาชิกผู้ใช้งาน", "name": "report_member"},
		{"title": "ยอดฝาก-ถอน", "name": "report_deposit"},
		{"title": "จำนวนฝาก-ถอนตามเวลา", "name": "report_deposit_time"},
		{"title": "รายงานการแจกโบนัส", "name": "report_bonus"},
		{"title": "จำนวนสมาชิกนับเวลาบันทึก", "name": "report_member_time"},
		{"title": "ยอดสมาชิกตามช่องทางที่รู้จัก", "name": "report_member_channel"},
		{"title": "จำนวนบันทึกรายการตามผู้ใช้งาน", "name": "report_member_user"},

		{"title": "รายงานการตลาด", "name": "marketing_report"},
		{"title": "คืนยอดเสีย", "name": "marketing_report_return"},
		{"title": "ลิงค์รับทรัพย์", "name": "marketing_report_link"},
		{"title": "พันธมิตร", "name": "marketing_report_partner"},

		{"title": "รายงานข้อมูล แพ้-ชนะ", "name": "report_winlose"},

		{"title": "รายงานกิจกรรม", "name": "activity_report"},
		{"title": "กงล้อนำโชค", "name": "activity_report_lucky"},

		{"title": "รายการฝาก-ถอนเสร็จสิ้น", "name": "deposit_withdraw"},
		{"title": "รายการโอนรอดำเนินการ", "name": "waiting_transfer"},
		{"title": "บันทึกรายการฝาก-ถอน", "name": "deposit_withdraw_history"},
		{"title": "อนุมัติฝาก(Auto)", "name": "auto_deposit"},
		{"title": "อนุมัติถอน(Auto)", "name": "auto_withdraw"},

		{"title": "ตั้งค่าระบบ", "name": "setting"},
		{"title": "ข้อมูลเบื้องต้น", "name": "setting_basic"},
		{"title": "แจ้งเตือนกลุ่ม line", "name": "setting_line"},
		{"title": "PushMessage line", "name": "setting_line_push"},
		{"title": "แจ้งเตือน Cyber Notify", "name": "setting_cyber"},

		{"title": "สถานะเรื่องแจ้งแก้ไข", "name": "status_update"},
		{"title": "ใบแจ้งหนี้", "name": "invoice_notice"},
	}

	for i, _ := range list {

		if list[i].Name == titles[0]["name"].(string) {
			list[i].Title = titles[0]["title"].(string)
		}

	}

	c.JSON(201, model.SuccessWithData{Message: "Success", Data: list})
}
