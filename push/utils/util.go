package util

import (
	"push/proto"
)

// 保证微服务独立性，离线消息的存储也给message吧
func storeOfflineMessage(receiver []string, pushMessage *proto.PushMessage) error {

	return nil
}
