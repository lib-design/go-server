// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func GetDetail(c helper.Context, id string) (res schema.Response) {
	var (
		err  error
		data = schema.Address{}
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

		helper.Response(&res, data, err)
	}()

	addressInfo := model.Address{
		Id:  id,
		Uid: c.Uid,
	}

	if err = database.Db.Model(&addressInfo).Where(&addressInfo).First(&addressInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressNotExist
		}
		return
	}

	if err = mapstructure.Decode(addressInfo, &data.AddressPure); err != nil {
		return
	}

	data.CreatedAt = addressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = addressInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetDetailRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param("address_id")

	res = GetDetail(helper.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, id)
}