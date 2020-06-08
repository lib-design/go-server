package worker

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"log"
)

func textMessageFromUserHandler(msg ws.Message) (err error) {
	waiterClient := ws.WaiterPoll.Get(msg.To)
	userClient := ws.UserPoll.Get(msg.From)

	if waiterClient == nil || userClient == nil {
		return
	}

	if err = waiterClient.WriteJSON(ws.Message{
		From:    msg.From,
		To:      msg.To,
		Type:    msg.Type,
		Payload: msg.Payload,
	}); err != nil {
		return
	}

	// 发送成功，写入聊天记录
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	hash := util.MD5(userClient.UUID + waiterClient.UUID)

	session := model.CustomerSession{
		Id:       hash,
		Uid:      userClient.GetProfile().Id,
		WaiterID: waiterClient.GetProfile().Id,
	}

	// 获取会话
	if err := tx.Model(model.CustomerSession{}).Where(&session).First(&session).Error; err != nil {
		return err
	}

	var raw []byte

	if raw, err = json.Marshal(msg.Payload); err != nil {
		return
	}

	sessionItem := model.CustomerSessionItem{
		Id:         session.Id,
		Type:       model.SessionTypeText,
		ReceiverID: waiterClient.GetProfile().Id,
		SenderID:   userClient.GetProfile().Id,
		Payload:    string(raw),
	}

	// 讲聊天记录写入会话
	if err := tx.Create(&sessionItem).Error; err != nil {
		return err
	}

	return
}

// 处理来之用户端的消息
func MessageFromUserHandler() {
	for {
		// 从客服池中取消息
		msg := <-ws.WaiterPoll.Broadcast

	typeCondition:
		switch ws.TypeResponseWaiter(msg.Type) {
		// 发送数据给客服
		case ws.TypeResponseWaiterMessageText:
			if err := textMessageFromUserHandler(msg); err != nil {
				log.Println(err)
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
}
