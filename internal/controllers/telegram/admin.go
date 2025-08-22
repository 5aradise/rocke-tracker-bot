package telegram

import (
	"bot/internal/utils/lang"

	"gopkg.in/telebot.v4"
)

func (h *Handler) adminMode(c telebot.Context) error {
	user := c.Sender()
	userID := user.ID
	userLang := user.LanguageCode

	var msg string
	h.adminModeMu.Lock()
	if _, inMode := h.inAdminMode[userID]; !inMode {
		h.inAdminMode[userID] = struct{}{}
		h.adminModeMu.Unlock()
		msg = youAreInAdminModeMsg.In(lang.Code(userLang))
	} else {
		delete(h.inAdminMode, userID)
		h.adminModeMu.Unlock()
		msg = youAreNotInAdminModeMsg.In(lang.Code(userLang))
	}
	return c.Send(msg)
}

func (h *Handler) onText(c telebot.Context) error {
	senderID := c.Sender().ID

	if senderID == h.adminID.value {
		return h.handleAdminMessage(c)
	}

	h.adminModeMu.RLock()
	_, inMode := h.inAdminMode[senderID]
	h.adminModeMu.RUnlock()
	if inMode {
		return h.handleAdminModeMessage(c)
	}

	return nil
}

func (h *Handler) handleAdminMessage(c telebot.Context) error {
	msg := c.Message()
	if !msg.IsReply() {
		return nil
	}

	user := msg.ReplyTo.OriginalSender
	_, err := c.Bot().Send(user, c.Text())
	return err
}

func (h *Handler) handleAdminModeMessage(c telebot.Context) error {
	return c.ForwardTo(h.adminID)
}
