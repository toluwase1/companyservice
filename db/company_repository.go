package db

import (
	"company-service/models"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// DB provides access to the different db

//go:generate mockgen -destination=../mocks/auth_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db AuthRepository
type CompanyRepository interface {
	CreateCompany(user *models.Company) (*models.Company, error)
	IsCompanyNameExist(name string) error
	FindCompanyByName(email string) (*models.Company, error)
	UpdateCompany(medication *models.Company, companyId string) error
	DeleteCompanyById(companyId string) error
	FindCompanyById(email string) (*models.Company, error)
	IsBusinessOwner(companyId string, userId float64) error
}

type companyRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *GormDB) CompanyRepository {
	return &companyRepo{db.DB}
}

func (a *companyRepo) CreateCompany(company *models.Company) (*models.Company, error) {
	err := a.DB.Create(company).Error
	if err != nil {
		return nil, fmt.Errorf("could not create company: %v", err)
	}
	return company, nil
}

func (a *companyRepo) FindCompanyByName(username string) (*models.Company, error) {
	db := a.DB
	user := &models.Company{}
	err := db.Where("email = ? OR username = ?", username, username).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (a *companyRepo) IsCompanyNameExist(name string) error {
	var count int64
	err := a.DB.Model(&models.Company{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "gorm.count error")
	}
	if count > 0 {
		return fmt.Errorf("name already in use")
	}
	return nil
}

func (m *companyRepo) UpdateCompany(company *models.Company, companyId string) error {
	err := m.DB.Model(&models.Company{}).
		Where("id = ?", companyId).
		Updates(company).Error
	if err != nil {
		return fmt.Errorf("could not update Company: %v", err)
	}
	return nil
}

func (a *companyRepo) DeleteCompanyById(companyId string) error {
	company := &models.Company{}
	err := a.DB.Where("company_id = ?", companyId).Find(company).Error
	if err != nil {
		return fmt.Errorf("could not find company to delete: %v", err)
	}
	err = a.DB.Delete(&models.Company{}, "id = ?", companyId).Error
	if err != nil {
		return fmt.Errorf("could not delete company: %v", err)
	}
	return nil
}

func (a *companyRepo) FindCompanyById(id string) (*models.Company, error) {
	company := &models.Company{}
	err := a.DB.Where("id = ?", id).Find(company).Error
	if err != nil {
		return &models.Company{}, fmt.Errorf("could not find any company with that id %v", err)
	}
	return company, nil
}

func (a *companyRepo) IsBusinessOwner(companyId string, userId float64) error {
	var count int64
	err := a.DB.Model(&models.Company{}).Where("id = ? AND user_id = ?", companyId, userId).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "gorm.count error")
	}
	if count < 1 {
		return fmt.Errorf("not authorized")
	}
	return nil
}
