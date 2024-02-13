package main

import (
	"fmt"
	AuthController "tomoncode/gologin/controller"
	UserController "tomoncode/gologin/controller/user"
	"tomoncode/gologin/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	Passworde string
	Fullname  string
	Avatar    string
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	orm.InitDB()

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	r.GET("/users/readall", UserController.ReadAll)

	r.Run("localhost:8080")
}
