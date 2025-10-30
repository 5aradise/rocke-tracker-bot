package telegram

import "gopkg.in/telebot.v4"

func (h *Handler) adminMode(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)

	var msg string
	h.adminModeMu.Lock()
	if _, inMode := h.inAdminMode[user.ID]; !inMode {
		h.inAdminMode[user.ID] = struct{}{}
		h.adminModeMu.Unlock()
		msg = youAreInAdminModeMsg.In(userLang)
	} else {
		delete(h.inAdminMode, user.ID)
		h.adminModeMu.Unlock()
		msg = youAreNotInAdminModeMsg.In(userLang)
	}
	return c.Send(msg)
}

func (h *Handler) onText(c telebot.Context) error {
	sender := c.Sender()

	if sender.ID == h.adminID.value {
		return h.handleAdminMessage(c)
	}

	h.adminModeMu.RLock()
	_, inMode := h.inAdminMode[sender.ID]
	h.adminModeMu.RUnlock()
	if inMode {
		return h.handleAdminModeMessage(c)
	}

	return nil
}

func (*Handler) handleAdminMessage(c telebot.Context) error {
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
