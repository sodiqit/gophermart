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

type Withdraws struct {
	ID        int32 `sql:"primary_key"`
	UserID    int32
	Amount    float64
	OrderID   string
	CreatedAt time.Time
}
