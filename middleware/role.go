package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type connectDB struct {
	db *gorm.DB
}

func Role(db *gorm.DB) connectDB {
	return connectDB{db}
}

func (r connectDB) CheckAdmin(perName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.MustGet("role").(string)

		if role != "SUPER_ADMIN" {

			adminId, err := c.Get("adminId")
			if !err {
				c.AbortWithStatusJSON(401, gin.H{
					"message": "Unauthorized",
				})
				return
			}

			var count int64

			if err := r.db.Table("Admin_permissions ap").
				Select("ap.id").
				Joins("LEFT JOIN Permissions p ON p.id = ap.permission_id").
				Where("ap.admin_id = ?", adminId).
				Where("p.name = ?", perName).
				Limit(1).
				Count(&count).Error; err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"message": "Internal Server Error",
				})
				return
			}

			if count == 0 {
				c.AbortWithStatusJSON(403, gin.H{
					"message": "Permission Denied",
				})
				return
			}

		}

		c.Next()

	}
}
