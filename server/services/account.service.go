package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rolandwarburton/playwright-server/errors"
	database "github.com/rolandwarburton/playwright-server/models"
	"gorm.io/gorm"
)

type IAccountService interface {
	GetAccount(db *gorm.DB, object *database.Account) (gin.H, *errors.RestError)
}

func GetAccount(db *gorm.DB, object *database.Account) (gin.H, *errors.RestError) {
	dbc := db.First(&object)

	if dbc.Error != nil {
		return nil, errors.NotFound(fmt.Sprintf("account not found: %v", dbc.Error))
	}

	return gin.H{
		"count":  dbc.RowsAffected,
		"result": object,
	}, nil
}
