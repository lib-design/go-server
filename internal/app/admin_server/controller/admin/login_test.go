// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func init() {
	admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "123456",
		Name:     "admin",
	}, true)
}

func TestLogin(t *testing.T) {
	// 登陆超级管理员-失败
	{
		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin123",
		})

		assert.Equal(t, exception.InvalidAccountOrPassword.Code(), r.Status)
		assert.Equal(t, exception.InvalidAccountOrPassword.Error(), r.Message)
	}

	// 登陆超级管理员-成功
	{
		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "123456",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := r.Decode(&adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, token.StateAdmin); er != nil {
			t.Error(er)
		} else {
			// 判断UID是否与用户一致
			assert.Equal(t, adminInfo.Id, c.Uid)
		}
	}
}

func TestLoginRouter(t *testing.T) {
	// 登陆无效的管理员账号
	{
		body, _ := json.Marshal(&admin.SignInParams{
			Username: "admin",
			Password: "invalid_password",
		})

		r := tester.HttpAdmin.Post("/v1/login", body, nil)

		assert.Equal(t, http.StatusOK, r.Code)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, exception.InvalidAccountOrPassword.Code(), res.Status)
		assert.Equal(t, exception.InvalidAccountOrPassword.Error(), res.Message)
	}

	// 登陆正确的管理员账号
	{
		body, _ := json.Marshal(&admin.SignInParams{
			Username: "admin",
			Password: "123456",
		})

		r := tester.HttpAdmin.Post("/v1/login", body, nil)

		assert.Equal(t, http.StatusOK, r.Code)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := res.Decode(&adminInfo); err != nil {
			t.Error(err)
		}

		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, token.StateAdmin); er != nil {
			t.Error(er)
		} else {
			// 判断UID是否与用户一致
			assert.Equal(t, adminInfo.Id, c.Uid)
		}
	}
}
