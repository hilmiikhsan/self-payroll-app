package repository

import (
	"context"
	"errors"
	"self-payrol/config"
	"self-payrol/model"

	"gorm.io/gorm"
)

type companyRepository struct {
	Cfg config.Config
}

func NewCompanyRepository(cfg config.Config) model.CompanyRepository {
	return &companyRepository{Cfg: cfg}
}

func (c *companyRepository) Get(ctx context.Context) (*model.Company, error) {
	company := new(model.Company)

	if err := c.Cfg.Database().WithContext(ctx).First(company).Error; err != nil {
		return nil, err
	}

	return company, nil
}

func (c *companyRepository) CreateOrUpdate(ctx context.Context, company *model.Company) (*model.Company, error) {
	companyModel := new(model.Company)

	err := c.Cfg.Database().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(companyModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(company).Error; err != nil {
					return err
				}
				return tx.First(companyModel).Error
			}
			return err
		}

		// TODO: tuliskan baris code untuk update data company
		if err := tx.Model(companyModel).Updates(company).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return companyModel, nil
}

func (c *companyRepository) DebitBalance(ctx context.Context, amount int, note string) error {
	company, err := c.Get(ctx)
	if err != nil {
		return errors.New("company data not found")
	}

	// TODO: tuliskan baris code untuk mengurangi balance
	company.Balance -= amount

	if err := c.Cfg.Database().WithContext(ctx).Model(company).Updates(company).Find(company).Error; err != nil {
		return err

	}

	if err := c.Cfg.Database().WithContext(ctx).Create(&model.Transaction{
		Amount: amount,
		Note:   note,
		Type:   model.TransactionTypeDebit,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (c *companyRepository) AddBalance(ctx context.Context, balance int) (*model.Company, error) {
	company, err := c.Get(ctx)
	if err != nil {
		return nil, errors.New("company data not found")
	}

	// TODO: tuliskan baris code untuk topup balance
	company.Balance += balance

	if err := c.Cfg.Database().WithContext(ctx).Model(company).Updates(company).Find(company).Error; err != nil {
		return nil, err
	}

	if err := c.Cfg.Database().WithContext(ctx).Create(&model.Transaction{
		Amount: balance,
		Note:   "Topup balance company",
		Type:   model.TransactionsTypeCredit,
	}).Error; err != nil {
		return nil, err
	}

	return company, nil
}
