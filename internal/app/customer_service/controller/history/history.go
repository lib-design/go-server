// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package history

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"sort"
	"time"
)

type History struct {
	ID       string               `json:"id"`       // 消息 ID
	Sender   schema.ProfilePublic `json:"sender"`   // 消息发送者
	Receiver schema.ProfilePublic `json:"receiver"` // 消息接受者
	Type     ws.TypeResponseUser  `json:"type"`     // 消息类型
	Payload  interface{}          `json:"payload"`  // 消息体
	Date     string               `json:"date"`     // 消息时间
}

type Session struct {
	Id        string               `json:"id"`         // 会话 ID
	User      schema.ProfilePublic `json:"user"`       // 用户信息
	Waiter    schema.ProfilePublic `json:"waiter"`     // 客服信息
	History   []History            `json:"history"`    // 历史消息
	CreatedAt string               `json:"created_at"` // 创建会话的时间
}

func SessionItemToMap(sessionItems []model.CustomerSessionItem) ([]History, error) {
	var (
		result = make([]History, 0)
		err    error
	)

	for _, item := range sessionItems {
		target := History{
			ID: item.Id,
			Sender: schema.ProfilePublic{
				Id:       item.Sender.Id,
				Username: item.Sender.Username,
				Nickname: item.Sender.Nickname,
				Avatar:   item.Sender.Avatar,
			},
			Receiver: schema.ProfilePublic{
				Id:       item.Receiver.Id,
				Username: item.Receiver.Username,
				Nickname: item.Receiver.Nickname,
				Avatar:   item.Receiver.Avatar,
			},
			Payload: item.Payload,
			Date:    item.CreatedAt.Format(time.RFC3339Nano),
		}

		switch item.Type {
		case model.SessionTypeText:
			target.Type = ws.TypeResponseUserMessageText

			var payload ws.MessageTextPayload

			if err = json.Unmarshal([]byte(item.Payload), &payload); err != nil {
				return nil, err
			}

			target.Payload = payload
		case model.SessionTypeImage:
			target.Type = ws.TypeResponseUserMessageImage

			var payload ws.MessageImagePayload

			if err = json.Unmarshal([]byte(item.Payload), &payload); err != nil {
				return nil, err
			}

			target.Payload = payload
		}

		result = append(result, target)
	}

	return result, nil
}

// 获取某个用户的聊天记录
func GetHistory(userID string, txs ...*gorm.DB) ([]History, error) {
	var (
		result = make([]History, 0)
		tx     *gorm.DB
		err    error
	)

	if len(txs) > 0 {
		tx = txs[0]
	}

	if tx == nil {
		tx = database.Db.Begin()
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()

	list := make([]model.CustomerSessionItem, 0)

	query := tx.Model(model.CustomerSessionItem{}).Where("sender_id = ?", userID).Or("receiver_id = ?", userID).Order("created_at DESC").Preload("Sender").Preload("Receiver").Limit(100)

	if err = query.Find(&list).Error; err != nil {
		return nil, err
	}

	for _, info := range list {
		target := History{
			ID: info.Id,
			Sender: schema.ProfilePublic{
				Id:       info.Sender.Id,
				Username: info.Sender.Username,
				Nickname: info.Sender.Nickname,
				Avatar:   info.Sender.Avatar,
			},
			Receiver: schema.ProfilePublic{
				Id:       info.Receiver.Id,
				Username: info.Receiver.Username,
				Nickname: info.Receiver.Nickname,
				Avatar:   info.Receiver.Avatar,
			},
			Payload: info.Payload,
			Date:    info.CreatedAt.Format(time.RFC3339Nano),
		}

		switch info.Type {
		case model.SessionTypeText:
			target.Type = ws.TypeResponseUserMessageText

			var payload ws.MessageTextPayload

			if err = json.Unmarshal([]byte(info.Payload), &payload); err != nil {
				return nil, err
			}

			target.Payload = payload
		case model.SessionTypeImage:
			target.Type = ws.TypeResponseUserMessageImage

			var payload ws.MessageImagePayload

			if err := json.Unmarshal([]byte(info.Payload), &payload); err != nil {
				return nil, err
			}

			target.Payload = payload
		}

		result = append(result, target)
	}

	return result, nil
}

type sessionMap map[string][]Session

// 获取客服最近的聊天记录
func GetWaiterSession(waiterID string, txs ...*gorm.DB) (result []Session, err error) {
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	}

	if tx == nil {
		tx = database.Db.Begin()
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()

	list := make([]model.CustomerSession, 0)

	query := tx.Model(model.CustomerSession{}).
		Where("waiter_id = ?", waiterID).
		Order("created_at DESC").
		Preload("User").Preload("Waiter").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			m := model.CustomerSessionItem{}
			return db.Order(fmt.Sprintf("%s.created_at DESC", m.TableName())).Limit(1)
		}).
		Preload("Items.Sender").
		Preload("Items.Receiver").
		Limit(100)

	if err = query.Find(&list).Error; err != nil {
		return
	}

	var noDuplicationMap = sessionMap{}

	for _, info := range list {
		histories, err := SessionItemToMap(info.Items)

		if err != nil {
			return nil, err
		}

		target := Session{
			//User: info.User,
			User: schema.ProfilePublic{
				Id:       info.User.Id,
				Username: info.User.Username,
				Nickname: info.User.Nickname,
				Avatar:   info.User.Avatar,
			},
			Waiter: schema.ProfilePublic{
				Id:       info.Waiter.Id,
				Username: info.Waiter.Username,
				Nickname: info.Waiter.Nickname,
				Avatar:   info.Waiter.Avatar,
			},
			History:   histories,
			CreatedAt: info.CreatedAt.Format(time.RFC3339Nano),
		}

		result = append(result, target)
	}

	// 去除重复的 session
	for _, data := range result {
		if sessions, ok := noDuplicationMap[data.User.Id]; ok {
			noDuplicationMap[data.User.Id] = append(sessions, data)
		} else {
			noDuplicationMap[data.User.Id] = []Session{data}
		}
	}

	targets := make([]Session, 0)

	// session 去除重
	for _, sessions := range noDuplicationMap {
		// 合并聊天记录
		histories := make([]History, 0)
		for _, session := range sessions {
			histories = append(histories, session.History...)
		}

		// 最新的历史消息在上面
		sort.SliceStable(histories, func(i, j int) bool { return histories[i].Date > histories[j].Date })

		targets = append(targets, Session{
			Id:        sessions[0].Id,
			User:      sessions[0].User,
			Waiter:    sessions[0].Waiter,
			History:   histories,
			CreatedAt: sessions[0].CreatedAt,
		})
	}

	result = targets

	return
}
