package web

import (
	"encoding/json"
	"github.com/guirialli/go-pos/clean_arch/internals/entity"
	"github.com/guirialli/go-pos/clean_arch/internals/usecase"
	"github.com/guirialli/go-pos/clean_arch/pkg/events"
	"net/http"
)

type OrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewOrderHandler(
	eventDispatcher events.EventDispatcherInterface,
	orderRepository entity.OrderRepositoryInterface,
	orderCreatedEvent events.EventInterface,
) *OrderHandler {
	return &OrderHandler{
		EventDispatcher:   eventDispatcher,
		OrderRepository:   orderRepository,
		OrderCreatedEvent: orderCreatedEvent,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	findAllOrder := usecase.NewFindAllOrdersUseCase(h.OrderRepository)

	output, err := findAllOrder.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
