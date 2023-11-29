package database

import (
	"time"

	"github.com/rolandwarburton/playwright-server/errors"
)

type Tab struct {
	ID        string    `gorm:"primarykey;size:36;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"notNull;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"notNull;type:timestamp" json:"updated_at"`
}

// this can be called inside the controller
func (tab *Tab) Validate() *errors.RestError {
	// TODO actually implement this
	if tab.ID == "123" {
		return errors.BadRequest("bad request")
	}
	return nil
}
