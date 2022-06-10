package postal

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Сервер для обмена сообщениями: обычные и PUSH
// Обычные сообщения передаются по протоколу http методами POST, в виде JSON
// PUSH сообщения передаются по протоколу websocket в виде JSON
type Postal struct {
	requestHandlers map[string]RequestHandler
	pushChan        chan pushMessage
}

// Контекст для обработчиков запросов
type RequestContext struct {
	fiberCtx *fiber.Ctx
}

// Обработчик запроса, получает запрос и возвращает ответ
type RequestHandler = func(ctx *RequestContext)

// Структура для отправки PUSH сообщения
type pushMessage struct {
	Channel     string      `json:"channel"`
	MessageName string      `json:"messageName"`
	Message     interface{} `json:"message"`
}

// Создаёт новый экземпляр
func New() *Postal {
	requestHandlers := make(map[string]RequestHandler)
	return &Postal{
		requestHandlers: requestHandlers,
	}
}

// Добавляет обработчик запроса на который ожидается ответ
func (post *Postal) AddRequestHandler(requestName string, handler RequestHandler) {
	post.requestHandlers[requestName] = handler
}

// Отправляет PUSH сообщение в канал
func (post *Postal) SendPush(channelName string, messageName string, message interface{}) {
	post.pushChan <- pushMessage{
		Channel:     channelName,
		MessageName: messageName,
		Message:     message,
	}
}

// Читает запрос в JSON структуру
func (post *RequestContext) ReadJson(message interface{}) {
	post.fiberCtx.BodyParser(message)
}

// Отправляет ответ на запрос
func (post *RequestContext) SendResponse(message interface{}) {
	post.fiberCtx.JSON(message)
}

// Запускает сервис
func (post *Postal) Listen(port int) {
	app := fiber.New()
	app.Post("/requests/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")

		handler := post.requestHandlers[name]
		if handler != nil {
			handler(&RequestContext{
				fiberCtx: ctx,
			})
		}

		return nil
	})

	pushChan := make(chan pushMessage)
	post.pushChan = pushChan

	app.Get("/push", websocket.New(func(ctx *websocket.Conn) {
		msg := <-post.pushChan
		ctx.WriteJSON(msg)
	}))

	app.Listen(fmt.Sprintf(":%d", port))
}
