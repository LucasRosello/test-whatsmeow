package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"

	_ "github.com/mattn/go-sqlite3" // Importa el controlador de SQLite3
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type InterfaceWhatsappClient interface {
}

type WhatsappClient struct {
	client *whatsmeow.Client
}

func NewWhatsappClient() (InterfaceWhatsappClient, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:../../app/data/examplestore.db?_foreign_keys=on", dbLog)
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
	// client.AddEventHandler(eventHandler)

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

	return &WhatsappClient{client: client}, nil
}

// func sendSimpleMessage(jid string /*aca deberia clavar una lista de mensajes y recorrerla*/, messageText string, destinatario string) error {
// 	dbLog := waLog.Stdout("Database", "DEBUG", true)
// 	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
// 	container, err := sqlstore.New("sqlite3", "file:../../app/data/examplestore.db?_foreign_keys=on", dbLog)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
// 	deviceStore, err := container.GetFirstDevice()
// 	if err != nil {
// 		panic(err)
// 	}
// 	clientLog := waLog.Stdout("Client", "DEBUG", true)
// 	client := whatsmeow.NewClient(deviceStore, clientLog)
// 	// client.AddEventHandler(eventHandler)

// 	if client.Store.ID == nil {
// 		// No ID stored, new login
// 		qrChan, _ := client.GetQRChannel(context.Background())
// 		err = client.Connect()
// 		if err != nil {
// 			panic(err)
// 		}
// 		for evt := range qrChan {
// 			if evt.Event == "code" {
// 				// Render the QR code here
// 				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
// 				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
// 				fmt.Println("QR code:", evt.Code)
// 			} else {
// 				fmt.Println("Login event:", evt.Event)
// 			}
// 		}
// 	} else {
// 		// Already logged in, just connect
// 		err = client.Connect()
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Convert the string constant to types.JID
// 		recipientJID := types.JID{User: "5491160312081", Server: "s.whatsapp.net"}

// 		var message waProto.Message

// 		msg := "¿Como esta mi princesita hermosa?"

// 		message.Conversation = &msg

// 		// Send a message saying "Hello world!" to 5491160312081
// 		response, err := client.SendMessage(context.Background(), recipientJID, &message)

// 		fmt.Println("1ban")
// 		fmt.Printf("Respuesta: %s Error: %s", response, err)
// 		fmt.Println("2ban")

// 	}

// 	return nil
// }
