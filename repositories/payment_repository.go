package repositories

import (
	"autoworkshop/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(p *models.Payment) error
	Delete(id string) error
	FindByInvoice(invoiceID string) ([]models.Payment, error)
}

type paymentRepo struct{ db *gorm.DB }

func NewPaymentRepository(db *gorm.DB) PaymentRepository { return &paymentRepo{db} }

func (r *paymentRepo) Create(p *models.Payment) error { return r.db.Create(p).Error }

func (r *paymentRepo) Delete(id string) error {
	return r.db.Delete(&models.Payment{}, "id = ?", id).Error
}

func (r *paymentRepo) FindByInvoice(invoiceID string) ([]models.Payment, error) {
	var list []models.Payment
	return list, r.db.Where("invoice_id = ?", invoiceID).Find(&list).Error
}
