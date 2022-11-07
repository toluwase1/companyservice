package server

import (
	"company-service/errors"
	_ "company-service/errors"
	"company-service/models"
	"company-service/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) HandleCreateCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Company
		if err := decode(c, &user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.CompanyService.CreateCompany(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Company Creation successful", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) HandleUpdateCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		//_, user, err := GetValuesFromContext(c)
		//if err != nil {
		//	err.Respond(c)
		//	return
		//}
		//id, errr := strconv.ParseUint(c.Param("companyID"), 10, 32)
		id := c.Param("companyID")
		var updateCompanyRequest models.Company
		if err := decode(c, &updateCompanyRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		err := s.CompanyService.UpdateCompany(&updateCompanyRequest, id)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "company updated successfully", http.StatusOK, nil, nil)
	}
}

func (s *Server) HandleGetCompanyDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("companyID")
		company, err := s.CompanyService.GetCompany(id)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "retrieved company successfully", http.StatusOK, gin.H{"company details": company}, nil)
	}
}

func (s *Server) handleDeleteCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		//_, user, err := GetValuesFromContext(c)
		//if err != nil {
		//	err.Respond(c)
		//	return
		//}
		id := c.Param("companyID")
		if err := s.CompanyService.DeleteCompanyById(id); err != nil {
			err.Respond(c)
			return
		}

		response.JSON(c, "company successfully deleted", http.StatusOK, nil, nil)
	}
}
