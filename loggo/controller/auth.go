package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"tomoncode/gologin/orm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

var hmacSampleSecret []byte

func Register(c *gin.Context) {
	// ตรวจสอบข้อมูล JSON ที่ส่งมา
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่ามีผู้ใช้ที่มีชื่อผู้ใช้เหมือนกับที่ส่งมาหรือไม่
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exists"})
		return
	}

	// สร้างรหัสผ่านเข้ารหัส
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// สร้างข้อมูลผู้ใช้ในฐานข้อมูล
	user := orm.User{
		Username:  json.Username,
		Passworde: string(encryptedPassword),
		Fullname:  json.Fullname,
		Avatar:    json.Avatar,
	}
	if err := orm.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// สร้าง JSON response แสดงว่าลงทะเบียนผู้ใช้สำเร็จ
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Registered", "userId": user.ID})
}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	// ตรวจสอบข้อมูล JSON ที่ส่งมา
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่ามีผู้ใช้ที่มีชื่อผู้ใช้เหมือนกับที่ส่งมาหรือไม่
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does NotExists"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(userExist.Passworde), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			"exp":    time.Now().Add(time.Minute * 1).Unix(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login Success", "token": tokenString})
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login Failed"})
	}
}
