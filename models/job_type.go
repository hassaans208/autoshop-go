package models

type JobType struct {
	ID              uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientCompanyID string `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_jt_name_cc" json:"-"`
	Name            string `gorm:"type:varchar(100);not null;uniqueIndex:idx_jt_name_cc" json:"name"`
	Subtypes        string `gorm:"type:json" json:"subtypes"` // JSON array of strings
}
