// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notify

type Event string

const (
	EventSendNotifyToAllUser             Event = "EventSendNotifyToAllUser"             // 推送给所有用户
	EventSendNotifyToCustomUser          Event = "EventSendNotifyToCustomUser"          // 推送给指定用户
	EventSendNotifyCheckUserLoginStatus  Event = "EventSendNotifyCheckUserLoginStatus"  // 推送检查用户登录状态
	EventSendNotifyToUserNewNotification Event = "EventSendNotifyToUserNewNotification" // 推送新的系统通知
	EventSendNotifyToUserNewMessage      Event = "EventSendNotifyToUserNewMessage"      // 推送指定用户的个人消息
)

type Notifier interface {
	SendNotifyToAllUser(headings string, content string, data map[string]interface{}) error                      // 向所有用户推送
	SendNotifyToCustomUser(userIds []string, headings string, content string, data map[string]interface{}) error // 推送自定义通知
	SendNotifySystemNotificationToUser(notificationId string) error                                              // 推送系统通知
	SendNotifyUserNewMessage(messageId string) error                                                             // 推送用户消息
	SendNotifyToUserForLoginStatus(userID string) error                                                          // 推送用户登录异常
}

var Notify = NewNotifierOneSignal()
