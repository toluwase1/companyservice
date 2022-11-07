package middleware

import (
	"company-service/db"
	"company-service/server/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
)

func IsCompanyOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim := c.MustGet("claim").(Claim)
		log.Println("my claims: ", claim)
		companyId := c.Param("companyID")
		err := db.NewAuthRepo(db.NewDB()).IsBusinessOwner(companyId, claim.UserId)
		if err != nil {
			response.Unauthorized(c, errors.Wrap(err, "you don't own this company").Error())
			return
		}
		c.Next()
	}
}
