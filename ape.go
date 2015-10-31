package ape

import (
	"fmt"
	"regexp"

	"github.com/mix3/guiniol"
	"golang.org/x/net/context"
)

type actionFunc func(*Event) error

type action struct {
	help string
	fn   actionFunc
}

type Connection struct {
	connection *guiniol.Connection
	actionMap  map[string]action
}

func New(token string) *Connection {
	conn := &Connection{
		connection: guiniol.NewConnection(token),
		actionMap:  map[string]action{},
	}
	conn.AddAction("help", "this message", func(e *Event) error {
		args := e.Command().Args()
		if 0 < len(args) {
			if v, ok := conn.actionMap[args[0]]; ok {
				m := "```\n"
				m += fmt.Sprintf("%s --- %s\n", args[0], v.help)
				m += "```"
				e.Reply(m)
			} else {
				e.Reply(fmt.Sprintf("command not found: %s", args[0]))
			}
			return nil
		}

		m := "command list\n```"
		for k, v := range conn.actionMap {
			m += fmt.Sprintf("    %s --- %s\n", k, v.help)
		}
		m += "```"
		e.Reply(m)
		return nil
	})
	return conn
}

func (c *Connection) AddAction(command, help string, fn actionFunc) {
	c.actionMap[command] = action{
		help: help,
		fn:   fn,
	}
}

var pattern1 = regexp.MustCompile(`^([^:]+):\s+(.+)`)   // bot_name: *****
var pattern2 = regexp.MustCompile(`^<@(\w+)>:?\s+(.+)`) // @bot_name *****

func matcher(e *guiniol.EventCtx) (string, string) {
	t := e.MessageEvent().Text
	if m := pattern2.FindStringSubmatch(t); len(m) == 3 {
		return m[1], m[2]
	}
	if m := pattern1.FindStringSubmatch(t); len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}

func (c *Connection) Loop() {
	c.connection.RegisterCb(recoverCb(c.actionCb))
	c.connection.Loop()
}

func (c *Connection) actionCb(ctx context.Context, eventCtx *guiniol.EventCtx) {
	if eventCtx.MessageEvent().Subtype == "bot_message" {
		return
	}

	name, m := matcher(eventCtx)
	if name != eventCtx.UserName() && name != eventCtx.UserId() {
		return
	}

	e := newEvent(ctx, eventCtx, m)

	if action, ok := c.actionMap[e.Command().Name()]; ok {
		if err := action.fn(e); err != nil {
			e.Reply(err.Error())
		}
	} else {
		e.Reply("???")
	}
}

func recoverCb(cb guiniol.CallbackFunc) guiniol.Callback {
	return guiniol.CallbackFunc(func(ctx context.Context, eventCtx *guiniol.EventCtx) {
		defer func() {
			if err := recover(); err != nil {
				e := newEvent(ctx, eventCtx, "")
				e.Reply(fmt.Sprintf("%v", err))
			}
		}()

		cb.Next(ctx, eventCtx)
	})
}
