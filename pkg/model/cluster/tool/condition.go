package tool

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Condition struct {
	common.BaseModel
	ID        string
	Name      string
	StatusID  string
	Status    string
	Message   string `gorm:"type:text"`
	OrderNum  int
	EndTime   time.Time
	StartTime time.Time
}

func (c *Condition) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c Condition) TableName() string {
	return "ko_tool_condition"
}
