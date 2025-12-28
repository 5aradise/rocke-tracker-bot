package telegram

import (
	"bot/pkg/lang"
	"fmt"

	"gopkg.in/telebot.v4"
)

type commands []struct {
	cmd         string
	description lang.String
	handler     telebot.HandlerFunc
}

func (cmds commands) set(b *telebot.Bot) error {
	langCmds := [4]struct {
		lang       string
		tgLangCode string
		lsLang     lang.Language
		cmds       []telebot.Command
	}{
		{
			lang:       "english",
			tgLangCode: englishLandCode,
			lsLang:     lang.English,
			cmds:       make([]telebot.Command, 0, len(cmds)),
		},
		{
			lang:       "ukrainian",
			tgLangCode: ukrainianLandCode,
			lsLang:     lang.Ukrainian,
			cmds:       make([]telebot.Command, 0, len(cmds)),
		},
		{
			lang:       "russian",
			tgLangCode: russianLandCode,
			lsLang:     lang.Russian,
			cmds:       make([]telebot.Command, 0, len(cmds)),
		},
		{
			lsLang: lang.Other,
			cmds:   make([]telebot.Command, 0, len(cmds)),
		},
	}

	for _, cmd := range cmds {
		b.Handle("/"+cmd.cmd, cmd.handler)

		for i, lc := range langCmds[:len(langCmds)-1] {
			langCmds[i].cmds = append(lc.cmds, telebot.Command{
				Text:        cmd.cmd,
				Description: cmd.description.In(lc.lsLang),
			})
		}
	}

	iother := len(langCmds) - 1
	for _, lc := range langCmds[:iother] {
		err := b.SetCommands(lc.tgLangCode, lc.cmds)
		if err != nil {
			return fmt.Errorf("setting up commands for %s language: %w", lc.lang, err)
		}
	}
	err := b.SetCommands(langCmds[iother].cmds)
	if err != nil {
		return fmt.Errorf("setting up commands for other language: %w", err)
	}
	return nil
}
