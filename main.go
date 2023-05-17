package main

import (
	"net/http"
	"time"

	"github.com/afiifatuts/go-auth-jwt/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// var jwtKey = []byte("SECRET_KEY")
var jwtKey = "SECRET_KEY"

func main() {
	//gin router
	r := gin.Default()
	//setup routes
	r.POST("/auth/login", loginHandler)

	//kita melalui login dulu
	userRouter := r.Group("api/v1/users")

	//tambahkan middleware
	userRouter.Use(auth.AuthMiddleware(jwtKey))

	//setup get user profile routes
	userRouter.GET("/:id/profile", profileHandler)

	//start server
	r.Run(":8800")

}

// loginHandler
func loginHandler(c *gin.Context) {
	var user User
	//binding : langsung ke objectnya
	//kalau scan di mapping
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//logic authentication(compare username dan password)

	if user.Username == "enigma" && user.Password == "12345" {
		//bikin code untuk generate token
		token := jwt.New(jwt.SigningMethodHS256)
		//cara claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

		//membuat tokennya
		tokenStr, err := token.SignedString(jwtKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenStr})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func profileHandler(c *gin.Context) {
	//ambil username dari jwt tokennya
	claims := c.MustGet("claims").(jwt.MapClaims)

	username := claims["username"].(string)

	//seharusnya response user dari db, tapi di contoh ini kita return usrname
	c.JSON(http.StatusOK, gin.H{"username": username})
}
