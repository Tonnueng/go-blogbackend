package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tonnueng/blogbackend/util"
)

func IsAuthenticate(c *fiber.Ctx)error {
	cookie := c.Cookies("jwt")

	if _, err := util.Parsejwt(cookie); err !=nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"ข้อความ":"ไม่ได้รับการตรวจสอบตัวตน",
		})
	}
	return c.Next()
	
}



