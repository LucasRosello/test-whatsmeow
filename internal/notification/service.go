package notification

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3" // Importa el controlador de SQLite3

	"github.com/lucasrosello/test-whatsmeow/internal/domain"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
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
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
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
		sendSimpleMessage()
	} else {
		_, err := s.repository.Save(ctx, notification)
		if err != nil {
			return domain.Notification{}, err
		}
	}

	return notification, nil
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

func sendSimpleMessage() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect

		fmt.Println(client.Store.ID)
		err = client.Connect()
		if err != nil {
			panic(err)
		}

		// Convert the string constant to types.JID
		recipientJID := types.JID{User: "5491160312081", Server: "s.whatsapp.net"}

		var message waProto.Message

		msg := "¿Como esta mi princesita hermosa?"

		message.Conversation = &msg

		// Send a message saying "Hello world!" to 5491160312081
		response, err := client.SendMessage(context.Background(), recipientJID, &message)

		fmt.Println("1ban")
		fmt.Printf("Respuesta: %s Error: %s", response, err)
		fmt.Println("2ban")

	}

	client.Disconnect()
}
