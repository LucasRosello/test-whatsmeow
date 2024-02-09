package notification

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/lucasrosello/test-whatsmeow/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var METODO_NO_IMPLEMENTADO = errors.New("Metodo no implementado")

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Notification, error)
	Get(ctx context.Context, id int) (domain.Notification, error)
	Exists(ctx context.Context, notificationCode int) bool
	Save(ctx context.Context, notification domain.Notification) (domain.Notification, error) // Updated return type
	Update(ctx context.Context, notification domain.Notification) error                      // Added method
	Delete(ctx context.Context, id int) error                                                // Added method
}

type repository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	collection := db.Collection("mycollection")

	// Crear un contexto con un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Intentar hacer ping a la base de datos
	err := db.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos: ", err)
	} else {
		log.Println("Conexión a la base de datos realizada con éxito")
	}

	return &repository{
		db:         db,
		collection: collection,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Notification, error) {
	var notifications []domain.Notification

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var notification domain.Notification
		if err := cursor.Decode(&notification); err != nil {
			log.Fatal(err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return notifications, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Notification, error) {
	return domain.Notification{}, METODO_NO_IMPLEMENTADO
}

func (r *repository) Exists(ctx context.Context, notificationCode int) bool {
	return false
}

func (r *repository) Save(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		log.Fatal(err)
		return domain.Notification{}, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	notification.ID = id
	return notification, nil
}

func (r *repository) Update(ctx context.Context, w domain.Notification) error {
	return METODO_NO_IMPLEMENTADO
}

func (r *repository) Delete(ctx context.Context, id int) error {
	return METODO_NO_IMPLEMENTADO
}
