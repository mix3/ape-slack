package ape

import (
	"fmt"
	"strings"

	"github.com/mix3/guiniol"
	"golang.org/x/net/context"
)

type Event struct {
	ctx      context.Context
	eventCtx *guiniol.EventCtx
	command  *Command
}

func (e *Event) Ctx() context.Context {
	return e.ctx
}

func (e *Event) EventCtx() *guiniol.EventCtx {
	return e.eventCtx
}

func (e *Event) Command() *Command {
	return e.command
}

func (e *Event) Message() string {
	return e.eventCtx.MessageEvent().Text
}

func (e *Event) Reply(message string) {
	e.eventCtx.Reply(fmt.Sprintf(
		"<@%s> %s %s",
		e.eventCtx.MessageEvent().User,
		message,
		e.Permalink(),
	))
}

func (e *Event) ReplyWithoutPermalink(message string) {
	e.eventCtx.Reply(fmt.Sprintf(
		"<@%s> %s",
		e.eventCtx.MessageEvent().User,
		message,
	))
}

func (e *Event) SendMessage(channel, message string) {
	e.eventCtx.SendMessage(channel, message)
}

func (e *Event) Permalink() string {
	return e.eventCtx.Permalink()
}

func newEvent(ctx context.Context, eventCtx *guiniol.EventCtx, message string) *Event {
	args := strings.Split(message, " ")
	return &Event{
		ctx:      ctx,
		eventCtx: eventCtx,
		command:  newCommand(args[0], args[1:]),
	}
}
