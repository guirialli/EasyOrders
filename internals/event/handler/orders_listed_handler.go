package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/guirialli/go-pos/clean_arch/pkg/events"
	"github.com/streadway/amqp"
)

type OrdersListedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrdersListedHandler(rabbitMQChannel *amqp.Channel) *OrdersListedHandler {
	return &OrdersListedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrdersListedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Orders listed: %v", event.GetPayload())

	jsonOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		fmt.Printf("Error marshaling payload: %v", err)
		return
	}

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	err = h.RabbitMQChannel.Publish(
		"amq.direct",
		"",
		false,
		false,
		msgRabbitmq,
	)
	if err != nil {
		fmt.Printf("Error publishing message: %v", err)
	}
}
