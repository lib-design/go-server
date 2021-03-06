// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package login

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

func GetLatestLoginLog(c helper.Context) (res schema.Response) {
	var (
		err  error
		data = schema.LogLogin{}
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	logInfo := model.LoginLog{
		Uid: c.Uid,
	}

	query := schema.Query{}

	query.Normalize()

	if err = query.Order(database.Db.Where(&logInfo).Preload("User")).First(&logInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(logInfo, &data.LogLoginPure); err != nil {
		return
	}

	if err = mapstructure.Decode(logInfo.User, &data.User); err != nil {
		return
	}

	data.CreatedAt = logInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = logInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetLoginLog(id string) (res schema.Response) {
	var (
		err  error
		data = schema.LogLogin{}
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	logInfo := model.LoginLog{
		Id: id,
	}

	query := schema.Query{}

	query.Normalize()

	if err = query.Order(database.Db.Where(&logInfo).Preload("User")).First(&logInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(logInfo, &data.LogLoginPure); err != nil {
		return
	}

	if err = mapstructure.Decode(logInfo.User, &data.User); err != nil {
		return
	}

	data.CreatedAt = logInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = logInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetLoginLogRouter = router.Handler(func(c router.Context) {
	id := c.Param("log_id")

	c.ResponseFunc(nil, func() schema.Response {
		return GetLoginLog(id)
	})
})
