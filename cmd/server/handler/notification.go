package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucasrosello/test-whatsmeow/internal/domain"
	"github.com/lucasrosello/test-whatsmeow/internal/notification"
)

type Notification struct {
	notificationService notification.Service
}

func NewNotification(e notification.Service) *Notification {
	return &Notification{
		notificationService: e,
	}
}

func (e *Notification) Get() gin.HandlerFunc {
	type response struct {
		Data domain.Notification `json:"data"`
	}

	return func(c *gin.Context) {

		ctx := context.Background()
		sel, err := e.notificationService.Get(ctx, int(0))
		if err != nil {
			//c.JSON(404, web.NewError(404, "Notification not found"))
			c.JSON(404, "Ejemplo no encontrado (o no se implemento el metodo, revisar codigo)")
			return
		}
		c.JSON(201, &response{sel})
	}
}

func (e *Notification) GetAll() gin.HandlerFunc {
	type response struct {
		Data []domain.Notification `json:"data"`
	}

	return func(c *gin.Context) {
		log.Println("se llego al handler")
		ctx := context.Background()
		notifications, err := e.notificationService.GetAll(ctx)
		if err != nil {
			c.JSON(404, err.Error())
			return
		}

		c.JSON(201, &response{notifications})
	}
}

func (e *Notification) Store() gin.HandlerFunc {
	type request struct {
		Sender     string    `json:"sender"`
		Receiver   string    `json:"receiver"`
		Message    string    `json:"message"`
		DateToSend time.Time `json:"dateToSend"`
	}

	type response struct {
		Data interface{} `json:"data"`
	}

	return func(c *gin.Context) {
		var req request

		err := c.Bind(&req)
		if err != nil {
			c.JSON(422, "json decoding: "+err.Error())
			return
		}

		noti := domain.Notification{
			Sender:     req.Sender,
			Receiver:   req.Receiver,
			Message:    req.Message,
			DateToSend: req.DateToSend,
		}

		ctx := context.Background()
		notification, err := e.notificationService.Store(ctx, noti)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(201, &response{notification})
	}
}

func (e *Notification) Update() gin.HandlerFunc {
	type request struct {
		Sender string `json:"sender"`
	}

	type response struct {
		Data interface{} `json:"data"`
	}

	return func(c *gin.Context) {
		var req request

		//paramID := c.Param("id")

		if err := c.Bind(&req); err != nil {
			c.JSON(422, "json decoding: "+err.Error())
			return
		}

		noti := domain.Notification{
			Sender: req.Sender,
		}

		ctx := context.Background()
		Notification, err := e.notificationService.Update(ctx, noti)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, &response{Notification})
	}
}

func (e *Notification) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, "invalid ID")
			return
		}

		ctx := context.Background()
		err = e.notificationService.Delete(ctx, int(id))
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, "The Notification has been deleted")
	}
}
