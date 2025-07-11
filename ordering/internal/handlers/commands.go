package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/application/commands"
	"eda-in-golang/ordering/orderingpb"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(app application.App) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{
		app: app,
	}
}

// This function registers the ordering module to listen for and handle incoming command messages from the message broker.
func RegisterCommandHandlers(subscriber am.CommandSubscriber, handlers ddd.CommandHandler[ddd.Command]) error {
	// Creates a Message Handler Wrapper
	cmdMsgHandler := am.CommandMessageHandlerFunc(func(ctx context.Context, cmdMsg am.IncomingCommandMessage) (ddd.Reply, error) {
		return handlers.HandleCommand(ctx, cmdMsg)
	})

	// Subscribes to Command Channel
	return subscriber.Subscribe(
		orderingpb.CommandChannel, // Channel: orderingpb.CommandChannel - the message broker channel/topic to listen on
		cmdMsgHandler,             // Handler: cmdMsgHandler - the function that processes incoming messages
		am.MessageFilter{ // Message Filter: Only processes RejectOrderCommand and ApproveOrderCommand messages
			orderingpb.RejectOrderCommand,
			orderingpb.ApproveOrderCommand,
		},
		am.GroupName("ordering-commands"), // Group Name: "ordering-commands" - consumer group for load balancing and fault tolerance
	)
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case orderingpb.RejectOrderCommand:
		return h.doRejectOrder(ctx, cmd)
	case orderingpb.ApproveOrderCommand:
		return h.doApproveOrder(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doRejectOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.RejectOrder)

	return nil, h.app.RejectOrder(ctx, commands.RejectOrder{ID: payload.GetId()})
}

func (h commandHandlers) doApproveOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.ApproveOrder)

	return nil, h.app.ApproveOrder(ctx, commands.ApproveOrder{
		ID:         payload.GetId(),
		ShoppingID: payload.GetShoppingId(),
	})
}
