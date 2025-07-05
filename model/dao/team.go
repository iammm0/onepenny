package dao

import "github.com/google/uuid"

// Team 是示例的组队/团队模型，你可以根据实际业务调整
type Team struct {
	BaseModel

	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Members     []User    `gorm:"many2many:team_members;"`
}
