package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/XineAurora/instantnotes-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SignUp(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 15)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := models.User{Name: body.Name, PasswordHash: string(hash), Email: body.Email}
	res := h.DB.Create(&user)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) SignIn(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}

	var user models.User
	h.DB.Where("email = ?", body.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid email or password",
			})
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       user.ID,
		"expiration": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to get token",
			})
			return
		}
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authToken", tokenString, int(time.Hour/time.Second)*24*7, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("authToken")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) >= claims["expiration"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		h.DB.First(&user, claims["user"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
