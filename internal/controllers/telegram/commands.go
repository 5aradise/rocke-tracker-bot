package telegram

import (
	"bot/internal/utils/lang"
	"fmt"

	"gopkg.in/telebot.v4"
)

type command struct {
	cmd         string
	description lang.String
	handler     telebot.HandlerFunc
}

type commands []command

func (cmds commands) set(b *telebot.Bot) error {
	var (
		en    = make([]telebot.Command, 0, len(cmds))
		uaru  = make([]telebot.Command, 0, len(cmds))
		other = make([]telebot.Command, 0, len(cmds))
	)
	for _, cmd := range cmds {
		b.Handle("/"+cmd.cmd, cmd.handler)

		en = append(en, telebot.Command{
			Text:        cmd.cmd,
			Description: cmd.description.In(lang.English),
		})
		uaru = append(uaru, telebot.Command{
			Text:        cmd.cmd,
			Description: cmd.description.In(lang.Ukrainian),
		})
		other = append(other, telebot.Command{
			Text:        cmd.cmd,
			Description: cmd.description.In(lang.Other),
		})
	}
	err := b.SetCommands(english, en)
	if err != nil {
		return fmt.Errorf("setting up commands for english languages: %w", err)
	}
	err = b.SetCommands(ukrainian, uaru)
	if err != nil {
		return fmt.Errorf("setting up commands for ukrainian languages: %w", err)
	}
	err = b.SetCommands(russian, uaru)
	if err != nil {
		return fmt.Errorf("setting up commands for russian languages: %w", err)
	}
	err = b.SetCommands(other)
	if err != nil {
		return fmt.Errorf("setting up commands for other languages: %w", err)
	}
	return nil
}
