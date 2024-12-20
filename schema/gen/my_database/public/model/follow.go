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

type Follow struct {
	UserID          int32     `sql:"primary_key"` // Идентификатор пользователя, который подписывается
	FollowingUserID int32     `sql:"primary_key"` // Идентификатор пользователя, на которого подписываются
	CreatedAt       time.Time // Дата и время создания подписки
}
