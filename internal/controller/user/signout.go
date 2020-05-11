// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user

import (
	"github.com/axetroy/go-server/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, schema.Response{
		Status:  schema.StatusSuccess,
		Message: "您已登出",
		Data:    true,
	})
}
