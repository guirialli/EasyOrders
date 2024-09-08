package usecase

import (
	"github.com/guirialli/go-pos/clean_arch/internals/entity"
	"github.com/guirialli/go-pos/clean_arch/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepositoryInterface,
	orderCreated events.EventInterface,
	orderCreatedEventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		EventDispatcher: orderCreatedEventDispatcher,
		OrderCreated:    orderCreated,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	err := order.CalculateFinalPrice()
	if err != nil {
		return OrderOutputDTO{}, err
	}
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}
	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}
	c.OrderCreated.SetPayload(dto)
	if err := c.EventDispatcher.Dispatcher(c.OrderCreated); err != nil {
		return OrderOutputDTO{}, err
	}
	return dto, nil
}
