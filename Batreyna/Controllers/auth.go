package Controllers

import (
	"log"
	"net/http"
	"time"

	"Batreyna/Models"
	"Batreyna/Token"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func CurrentUser(c *gin.Context) {
	user_id, err := Token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Logout(c *gin.Context) {
	token, err := Token.ExtractJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = false
	claims["exp"] = time.Now()
	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Claims = token2.Claims
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := Models.User{}

	user.Username = input.Username
	user.Password = input.Password

	uid, token, err := Models.LoginCheck(user.Username, user.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or password is incorrect."})
		return
	}

	user, _ = Models.GetUserByID(uid)

	c.JSON(http.StatusOK, gin.H{"message": "Login Successful", "jwt": token})

}

type RegisterInput struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Permission int    `json:"permission"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := Models.User{}

	user.Username = input.Username
	user.Password = input.Password
	user.Permission = input.Permission
	_, err := user.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "validated"})
}

func DeleteUser(c *gin.Context) {
	user_id, err := Token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := Models.DB.Unscoped().Delete(&Models.User{}, user.ID).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account Deleted Successfully"})
}
