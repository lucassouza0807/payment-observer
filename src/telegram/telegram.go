package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	lucasTypes "pagamento/src/headers"
)

const telegramAPI = "https://api.telegram.org/bot%s/%s"

func SendMessage(text string) error {

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	url := fmt.Sprintf(telegramAPI, token, "sendMessage")

	message := lucasTypes.SendMessageRequest{
		ChatID: os.Getenv("TELEGRAM_CHANEL_ID"),
		Text:   text,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao enviar mensagem, status: %s", resp.Status)
	}

	return nil
}
