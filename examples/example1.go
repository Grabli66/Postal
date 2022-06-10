package main

import (
	"fmt"
	"time"

	"rapidfingers.com/postal"
)

func main() {
	post := postal.New()

	post.AddRequestHandler("sum", func(ctx *postal.RequestContext) {
		message := &struct {
			N1 int `json:"n1"`
			N2 int `json:"n2"`
		}{}

		ctx.ReadJson(message)

		sum := message.N1 + message.N2

		ctx.SendResponse(&struct {
			Sum int `json:"sum"`
		}{
			Sum: sum,
		})

		post.SendPush("push", "sum", &struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("Sum: %d", sum),
		})
	})

	go func() {
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-ticker.C:
				post.SendPush("notify", "info", &struct {
					Name string `json:"Name"`
				}{
					Name: "Please send data to calculate",
				})
			}
		}
	}()

	post.Listen(26701)
}
