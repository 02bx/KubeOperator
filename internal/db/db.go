package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

const phaseName = "db"

type InitDBPhase struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func (i *InitDBPhase) Init() error {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		i.User,
		i.Password,
		i.Host,
		i.Port,
		i.Name)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return err
	}
	DB = db
	return nil
}


func (i *InitDBPhase) PhaseName() string {
	return phaseName
}
