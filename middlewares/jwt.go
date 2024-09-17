package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		message := map[string]interface{}{
			"message": nil,
		}

		c, err := r.Cookie("token")
		if err != nil {
			log.Println("Missing token cookie:", err)
			message["message"] = "Token is missing"
			helper.Response(w, message, http.StatusUnauthorized)
			return 
		}

		// mengambil value token
		tokenString := c.Value
		claims := &config.JWTClaim{}
		//parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error){
			return config.JWT_KEY, nil
		})
		if err != nil {
			switch err {
			case jwt.ErrTokenSignatureInvalid:
				log.Println("Invalid token signature:", err)
				message["message"] = "Unauthorized"
				helper.Response(w, message, http.StatusUnauthorized)
				return
			case jwt.ErrTokenExpired:
				log.Println("Token have expired", err)
				message["message"] = "Unauthorized"
				helper.Response(w, message, http.StatusUnauthorized)
			default:
				log.Println("Error Parsing token:", err)
				message["message"] = err.Error()
				helper.Response(w, message, http.StatusUnauthorized)
				return
			}
		}

		// memeriksa apakah token available and valid
		if claims, ok := token.Claims.(*config.JWTClaim); ok && token.Valid{
			role := claims.Role
			endpoint := r.URL.Path

			if err := endPoinCanAccess(role, endpoint); err != nil {
				log.Println("Access denied to endpoint:", err)
				message["message"] = err.Error()
				helper.Response(w, message, http.StatusUnauthorized)
				return
			}
		} else {
			message["message"] = "Unauthorized"
			helper.Response(w, message, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
} 

func endPoinCanAccess(role, endpoint string) error {
	var endpoin = []string{
		"/berkahjaya/adminside/barang",
		"/berkahjaya/adminside/barang/inputbarang",
		"/berkahjaya/adminside/barang/updatebarang",
		"/berkahjaya/adminside/barang/deletebarang",
		"/berkahjaya/adminside/hadiah",
		"/berkahjaya/adminside/hadiah/inputhadiah",
		"/berkahjaya/adminside/hadiah/updatehadiah",
		"/berkahjaya/adminside/hadiah/deletehadiah",
		"/berkahjaya/adminside/pengajuan/poin",
		"/berkahjaya/adminside/pengajuan/poin/sendmsgggiftsarrive",
		"/berkahjaya/adminside/pengajuan/poin/finished",
		"/berkahjaya/adminside/pengajuan/poin/verify",
		"/berkahjaya/adminside/pengajuan/poin/verify/cancel",
		"/berkahjaya/adminside/pengajuan/hadiah",
		"/berkahjaya/adminside/hadiah/search",
		"/berkahjaya/adminside/barang/search",
	}

	if role == "Admin" {
		for _, en := range endpoin {
			if endpoint == en {
				return nil
			}
		}
	}

	return fmt.Errorf("access denied to endpoint: %s", endpoint)
}