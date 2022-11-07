package services

import (
	"company-service/config"
	"company-service/db"
	"company-service/errors"
	apiError "company-service/errors"
	"company-service/models"
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService
// CompanyService interface
type CompanyService interface {
	CreateCompany(request *models.Company) (*models.Company, *apiError.Error)
	UpdateCompany(request *models.Company, companyID string) *errors.Error
	DeleteCompanyById(email string) *apiError.Error
	GetCompany(companyId string) (*models.Company, error)
}

// companyService struct
type companyService struct {
	Config      *config.Config
	companyRepo db.CompanyRepository
	kafka       db.KafkaRepository
}

// NewCompanyService instantiate an companyService
func NewCompanyService(authRepo db.CompanyRepository, kafkaRepo db.KafkaRepository, conf *config.Config) CompanyService {
	return &companyService{
		Config:      conf,
		companyRepo: authRepo,
		kafka:       kafkaRepo,
	}
}

var validTypes = map[string]bool{"Corporations": true, "NonProfit": true, "Cooperative": true, "Sole Proprietorship": true}

func (a *companyService) CreateCompany(company *models.Company) (*models.Company, *apiError.Error) {
	company.Id = uuid.New().String()
	err := a.companyRepo.IsCompanyNameExist(company.Name)
	if err != nil {
		return nil, apiError.New("company name already exist, please choose another name", http.StatusBadRequest)
	}
	if !validTypes[company.Type] {
		return nil, apiError.New("invalid company type", http.StatusBadRequest)
	}
	company, err = a.companyRepo.CreateCompany(company)
	if err != nil {
		log.Printf("unable to create company: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	_, _ = a.kafka.AddMessageToKafka(company, context.Background())

	return company, nil
}

func (m *companyService) UpdateCompany(request *models.Company, companyID string) *errors.Error {

	if !validTypes[request.Type] {
		return apiError.New("invalid company type", http.StatusBadRequest)
	}

	company := models.Company{
		Name:              request.Name,
		Description:       request.Description,
		AmountOfEmployees: request.AmountOfEmployees,
		IsRegistered:      request.IsRegistered,
		Type:              request.Type,
	}
	//get company where user and company id is defined above then send it for updating
	err := m.companyRepo.UpdateCompany(&company, companyID)
	if err != nil {
		return errors.ErrInternalServerError
	}
	_, _ = m.kafka.AddMessageToKafka(&company, context.Background())
	return nil
}

func (a *companyService) DeleteCompanyById(companyId string) *apiError.Error {
	err := a.companyRepo.DeleteCompanyById(companyId)
	if err != nil {
		return apiError.ErrInternalServerError
	}
	coy := models.Company{
		Id: companyId,
	}
	_, _ = a.kafka.AddMessageToKafka(&coy, context.Background())
	return nil
}

func (s *companyService) GetCompany(companyId string) (*models.Company, error) {
	coy, err := s.companyRepo.FindCompanyById(companyId)
	if err != nil {
		return &models.Company{}, apiError.ErrInternalServerError
	}
	return coy, nil
}
