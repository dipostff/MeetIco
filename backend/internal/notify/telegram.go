package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type TelegramNotifier interface {
	SendMessage(chatID, text string) error
}

type Bot struct {
	token string
}

func NewBot(token string) TelegramNotifier {
	if token == "" {
		return &Noop{}
	}
	return &Bot{token: token}
}

func (b *Bot) SendMessage(chatID, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(respBody))
	}

	return nil
}

type Noop struct{}

func (n *Noop) SendMessage(chatID, text string) error {
	slog.Info("telegram message (noop)", "chat_id", chatID, "text", text)
	return nil
}
