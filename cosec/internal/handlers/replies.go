package handlers

import (
	"context"
	"fmt"

	"eda-in-golang/cosec/internal/models"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/sec"
)

func RegisterReplyHandlers(subscriber am.ReplySubscriber, orchestrator sec.Orchestrator[*models.CreateOrderData]) error {
	fmt.Println("=== [Cosec module] Registering reply handlers")
	replyMsgHandler := am.MessageHandlerFunc[am.IncomingReplyMessage](func(ctx context.Context, replyMsg am.IncomingReplyMessage) error {
		fmt.Println("=== [Cosec module] Handling reply:", replyMsg.ReplyName())
		return orchestrator.HandleReply(ctx, replyMsg)
	})
	return subscriber.Subscribe(orchestrator.ReplyTopic(), replyMsgHandler, am.GroupName("cosec-replies"))
}
