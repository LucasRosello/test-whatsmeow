package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mau.fi/whatsmeow"

	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/gin-gonic/gin"
	"github.com/lucasrosello/test-whatsmeow/cmd/server/handler"
	"github.com/lucasrosello/test-whatsmeow/internal/notification"
	_ "github.com/mattn/go-sqlite3" // Importa el controlador de SQLite3
	"github.com/mdp/qrterminal"
)

// func eventHandler(evt interface{}) {
// 	switch v := evt.(type) {
// 	case *events.Message:
// 		fmt.Println("Received a message!", v.Message.GetConversation())

// 	}
// }

func main() {
	// Establecer un contexto con un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Conectar a MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Obtener una referencia a la base de datos
	db := mongoClient.Database("mydatabase")

	fmt.Println(db)

	var whatsappClient *whatsmeow.Client

	whatsappClient, _ = sendSimpleMessage("", "", "")

	router := gin.Default()

	notificationRepository := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepository, whatsappClient)
	notificationHandler := handler.NewNotification(notificationService)
	NotificationsRoutes := router.Group("/api/v1/notification")
	{
		NotificationsRoutes.GET("/", notificationHandler.GetAll())
		NotificationsRoutes.GET("/:id", notificationHandler.Get())
		NotificationsRoutes.POST("/", notificationHandler.Store())
		NotificationsRoutes.PATCH("/:id", notificationHandler.Update())
		NotificationsRoutes.DELETE("/:id", notificationHandler.Delete())
	}

	router.Run() // Update the URL to match your Docker Compose configuration

}

func sendSimpleMessage(jid string /*aca deberia clavar una lista de mensajes y recorrerla*/, messageText string, destinatario string) (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:../../app/data/sessions.db?_foreign_keys=on", dbLog)
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

	}

	return client, nil
}
