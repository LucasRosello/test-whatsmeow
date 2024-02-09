package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	waProto "go.mau.fi/whatsmeow/binary/proto"

	_ "github.com/mattn/go-sqlite3" // Importa el controlador de SQLite3

	"github.com/lucasrosello/test-whatsmeow/internal/domain"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Notification, error)
	GetAll(ctx context.Context) ([]domain.Notification, error)
	Store(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	Update(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository     Repository
	whatsappClient *whatsmeow.Client
}

func NewService(repository Repository, whatsappClient *whatsmeow.Client) Service {
	return &service{
		repository:     repository,
		whatsappClient: whatsappClient,
	}
}

func (s *service) Get(ctx context.Context, id int) (domain.Notification, error) {

	return s.repository.Get(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]domain.Notification, error) {
	log.Println("se llego al servicio")
	return s.repository.GetAll(ctx)
}

func (s *service) Store(ctx context.Context, notification domain.Notification) (domain.Notification, error) {

	if notification.DateToSend.Before(time.Now()) {
		fmt.Println("La fecha DateToSend es anterior a ahora.")
		// Convert the string constant to types.JID
		recipientJID := types.JID{User: notification.Receiver, Server: "s.whatsapp.net"}

		message := waProto.Message{
			Conversation: &notification.Message,
		}

		// Send a message saying "Hello world!" to 5491160312081
		response, err := s.whatsappClient.SendMessage(context.Background(), recipientJID, &message)

		fmt.Printf("STORE: %s Error: %s", response, err)

		return domain.Notification{}, nil
	} else {
		_, err := s.repository.Save(ctx, notification)
		if err != nil {
			return domain.Notification{}, err
		}
	}

	return domain.Notification{}, nil

}

func (s *service) Update(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	err := s.repository.Update(ctx, notification)
	if err != nil {
		return domain.Notification{}, err
	}

	return notification, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())

	}
}
