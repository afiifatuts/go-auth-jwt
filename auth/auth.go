package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// midleware mereturn handlefunc
func AuthMiddleware(jwtKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenstr := ctx.GetHeader("Authorization")

		//kalau tokennya kosong abort
		if tokenstr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorize"})
			ctx.Abort()
			return
		}

		//kalau tokennya ada di parse kemudian return tokennya
		// token str beda dengan token
		token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (any, error) { return jwtKey, nil })

		//jika gagal validasi token atau ada error maka abort
		if !token.Valid || err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorize"})
			ctx.Abort()
			return
		}

		//jika berhasil
		claims := token.Claims.(jwt.MapClaims)
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
