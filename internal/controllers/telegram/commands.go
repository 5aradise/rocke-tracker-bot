package telegram

import (
	"fmt"

	"gopkg.in/telebot.v4"
)

type command struct {
	cmd         string
	description langString
	handler     telebot.HandlerFunc
}

type commands []command

func (cmds commands) set(b *telebot.Bot) error {
	var (
		uaru  = make([]telebot.Command, 0, len(cmds))
		other = make([]telebot.Command, 0, len(cmds))
	)
	for _, cmd := range cmds {
		b.Handle("/"+cmd.cmd, cmd.handler)

		other = append(other, telebot.Command{
			Text:        cmd.cmd,
			Description: cmd.description.other,
		})
		uaru = append(uaru, telebot.Command{
			Text:        cmd.cmd,
			Description: cmd.description.uaru,
		})
	}
	err := b.SetCommands(other)
	if err != nil {
		return fmt.Errorf("setting up commands for other languages: %w", err)
	}
	err = b.SetCommands(ukrainian, uaru)
	if err != nil {
		return fmt.Errorf("setting up commands for ukrainian languages: %w", err)
	}
	err = b.SetCommands(russian, uaru)
	if err != nil {
		return fmt.Errorf("setting up commands for russian languages: %w", err)
	}
	return nil
}
