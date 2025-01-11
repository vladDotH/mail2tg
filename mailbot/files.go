package mailbot

import (
	"log"
	"os"

	"github.com/google/uuid"
)

func (bot *Bot) SaveMessage(html string) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return id, err
	}

	err = os.WriteFile(UUID2Path(id), []byte(html), 0644)

	if err != nil {
		log.Printf("Error while saving email: %e", err)
	}

	return id, err
}
