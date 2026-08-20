package am

import (
	"strings"
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	CommandHdrPrefix = "COMMAND_"
	CommandNameHdr = CommandHdrPrefix + "NAME"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"
)

type (
	CommandMessageHandler = MessageHandler[CommandMessage]

	Command interface {
		ddd.Command
		Destination() string
	}

	command struct {
		ddd.Command
		destination string
	}
)

func NewCommand(name, destination string, data []byte, options ...ddd.CommandOption) Command {
	return command{
		Command: ddd.NewCommand(name, data, options...),
		destination: destination,
	}
}

func (c command) Destination() string {
	return c.destination
}

type commandMsgHandler struct {
	publisher ReplyPublisher
	handler   ddd.CommandHandler[ddd.Command]
}

func NewCommandMessageHandler(publisher ReplyPublisher, handler ddd.CommandHandler[ddd.Command]) RawMessageHandler {
	return commandMsgHandler{
		publisher: publisher,
		handler:   handler,
	}
}

func (h commandMsgHandler) HandleMessage(ctx context.Context, msg RawMessage) error {
	var commandData api.CommandMessageData

	err := proto.Unmarshal(msg.Data(), &commandData)
	if err != nil {
		return err
	}

	commandMsg := commandMessage{
		id:         msg.ID(),
		name:       msg.MessageName(),
		data:       commandData.GetPayload(),
		metadata:   commandData.GetMetadata().AsMap(),
		occurredAt: commandData.GetOccurredAt().AsTime(),
	}

	destination := commandMsg.Metadata().Get(CommandReplyChannelHdr).(string)

	reply, err := h.handler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publishReply(ctx, destination, h.failure(reply, commandMsg))
	}

	return h.publishReply(ctx, destination, h.success(reply, commandMsg))
}

func (h commandMsgHandler) publishReply(ctx context.Context, destination string, reply ddd.Reply) error {
	return h.publisher.Publish(ctx, destination, reply)
}

func (h commandMsgHandler) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeFailure)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMsgHandler) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeSuccess)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMsgHandler) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	for key, value := range cmd.Metadata() {
		if key == CommandNameHdr {
			continue
		}

		if strings.HasPrefix(key, CommandHdrPrefix) {
			hdr := ReplyHdrPrefix + key[len(CommandHdrPrefix):]
			reply.Metadata().Set(hdr, value)
		}
	}

	return reply
}