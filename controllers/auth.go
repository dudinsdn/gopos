package controllers

import (
	"din/gopos/models"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func isValidEmail(email string) bool {
	rx := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return rx.MatchString(email)
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db}
}

func validateSignupInput(user *models.User) string {
	if strings.TrimSpace(user.Name) == "" {
		return "name is required!"
	}
	if strings.TrimSpace(user.Email) == "" {
		return "email is required!"
	}
	if !isValidEmail(user.Email) {
		return "email format is invalid"
	}
	if strings.TrimSpace(user.Password) == "" {
		return "password is required!"
	}
	if len(user.Password) < 6 {
		return "password must be at least 6 characters!"
	}
	return ""
}

func (ac *AuthController) SignUp(ctx *gin.Context) {
	var input models.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if msg := validateSignupInput(&input); msg != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var existing models.User
	if err := ac.db.Where("email = ?", input.Email).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	if err := ac.db.Create(&input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user": gin.H{
			"id":    input.ID,
			"email": input.Email,
			"name":  input.Name,
		},
	})
}

func validateLoginInput(user *models.User) string {
	if strings.TrimSpace(user.Email) == "" {
		return "email is required!"
	}
	if !isValidEmail(user.Email) {
		return "email format is invalid"
	}
	if strings.TrimSpace(user.Password) == "" {
		return "password is required!"
	}
	return ""
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var input models.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if msg := validateLoginInput(&input); msg != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var user models.User

	if err := ac.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ac *AuthController) Profile(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
