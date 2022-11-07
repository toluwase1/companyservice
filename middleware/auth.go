package middleware

import (
	"company-service/server/response"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"net/url"
)

const AuthHeader = "Authorization"

type Claim struct {
	UserId float64 `json:"user_id"`
	Email  string  `json:"email"`
	Exp    string  `json:"exp"`
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get(AuthHeader)
		if accessToken == "" {
			response.Unauthorized(c, "access token required")
			return
		}

		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "localhost:8081",
				Path:   "/api/v1/auth/authenticate"},
			Header: http.Header{
				"Authorization": []string{accessToken},
			},
		}
		fmt.Println(req.URL.String())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			response.Unauthorized(c, err.Error())
			return
		}

		if resp.StatusCode != http.StatusOK {
			response.Unauthorized(c, fmt.Sprintf("status code: %d", resp.StatusCode))
			return
		}

		var claims jwt.MapClaims
		err = json.NewDecoder(resp.Body).Decode(&claims)

		data := claims["data"].(map[string]interface{})
		id := data["id"].(float64)
		email := data["email"].(string)

		log.Println(id, email)
		log.Printf("%T", data)
		claim := Claim{
			UserId: id,
			Email:  email,
		}
		if err != nil {
			return
		}
		c.Set("claim", claim)
		
		c.Next()
	}
}
