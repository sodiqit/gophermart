//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Orders struct {
	ID        int32 `sql:"primary_key"`
	UserID    int32
	Status    string
	Accrual   *float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
