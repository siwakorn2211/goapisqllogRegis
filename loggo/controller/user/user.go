package user

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret []byte

func ReadAll(c *gin.Context) {
	hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	header := c.Request.Header.Get("Authorization")
	// ลบช่องว่างที่อาจเป็นไปได้ก่อนคำว่า "Bearer"
	tokenString := strings.TrimSpace(strings.ReplaceAll(header, "Bearer", ""))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["userId"])
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "claims": claims})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": err.Error()})
		return
	}
}
