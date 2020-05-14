// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package menu_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/menu"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 menu
	{
		var (
			name = "test"
		)

		r := menu.Create(helper.Context{
			Uid: adminInfo.Id,
		}, menu.CreateMenuParams{
			Name: name,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Menu{}

		assert.Nil(t, r.Decode(&n))

		defer menu.DeleteMenuById(n.Id)

		assert.NotNil(t, n.Id)
		assert.Equal(t, name, n.Name)
		assert.Equal(t, 0, n.Sort)
		assert.Equal(t, "", n.ParentId)
		assert.Equal(t, "", n.Url)
		assert.Equal(t, "", n.Icon)
		assert.Equal(t, []string{}, n.Accession)
	}

	// 创建一个 menu
	{
		var (
			name      = "test"
			accession = []string{"invalid"}
		)

		r := menu.Create(helper.Context{
			Uid: adminInfo.Id,
		}, menu.CreateMenuParams{
			Name:      name,
			Accession: &accession,
		})

		assert.Equal(t, exception.InvalidParams.Code(), r.Status)
		assert.Equal(t, exception.InvalidParams.Error(), r.Message)
	}

	// 非管理员的uid去创建，应该报错
	{

		userInfo, _ := tester.CreateUser()

		defer tester.DeleteUserByUserName(userInfo.Username)

		var (
			name = "test"
		)

		r := menu.Create(helper.Context{
			Uid: userInfo.Id,
		}, menu.CreateMenuParams{
			Name: name,
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.AdminNotExist.Error(), r.Message)
	}
}

func TestCreateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建 menu
	{
		var (
			name = "test"
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&menu.CreateMenuParams{
			Name: name,
		})

		r := tester.HttpAdmin.Post("/v1/menu", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.Menu{}

		assert.Nil(t, res.Decode(&n))

		defer menu.DeleteMenuById(n.Id)

		assert.Equal(t, name, n.Name)
		assert.Equal(t, 0, n.Sort)
	}
}
