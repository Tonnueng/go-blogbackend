package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/tonnueng/blogbackend/database"
	"github.com/tonnueng/blogbackend/models"
	"github.com/tonnueng/blogbackend/util"
)
func validateEmail(email string) bool {
	Re := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]+`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data) ;err!=nil{
		fmt.Println("Unable to parse body")
	}
	if len(data["password"].(string))<=6{
		c.Status(400)
		return c.JSON(fiber.Map{
			"ข้อความ":"รหัสต้องมากกว่า 6 ตัวอักษร",
		})
	}

	if !validateEmail(strings.TrimSpace(data["email"].(string))){
		c.Status(400)
		return c.JSON(fiber.Map{
			"ข้อความ":"อีเมลล์ไม่ถูกค้อง",
		})
	}

	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id!=0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"ข้อความ":"อีเมลล์ถูกใช้งานแล้ว",
		})
	}


	user := models.User{
		Firstname: data["first_name"].(string),
		Lastname: data["last_name"].(string),
		Phone: data["phone"].(string),
		Email: strings.TrimSpace(data["email"].(string)),
	}
	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":user,
		"ข้อควาท":"สร้างบัญชีผู้ใช้สำเร็จ", 
	})

}

func Login(c *fiber.Ctx) error {
    var data map[string]string

    if err := c.BodyParser(&data); err != nil {
        log.Printf("Error parsing request body: %v\n", err)
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "message": "Bad Request",
        })
    }

    var user models.User
    database.DB.Where("email=?", data["email"]).First(&user)
    if user.Id == 0 {
        c.Status(fiber.StatusNotFound)
        return c.JSON(fiber.Map{
            "message": "Email address doesn't exist, kindly create an account",
        })
    }

    if err := user.ComparePassword(data["password"]); err != nil {
        log.Printf("Error comparing password: %v\n", err)
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "message": "Incorrect password",
        })
    }

    token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))

    if err != nil {
        log.Printf("Error generating JWT: %v\n", err)
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "message": "Internal Server Error",
        })
    }

    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(time.Hour * 24),
        HTTPOnly: true,
    }
    c.Cookie(&cookie)

    return c.JSON(fiber.Map{
        "message": "Login successful",
        "user":    user,
    })
}

type Claims struct{
	jwt.StandardClaims
}