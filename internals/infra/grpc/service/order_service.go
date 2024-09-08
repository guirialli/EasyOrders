package service

import (
	"context"
	"github.com/guirialli/go-pos/clean_arch/internals/infra/grpc/pb"
	"github.com/guirialli/go-pos/clean_arch/internals/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase   usecase.CreateOrderUseCase
	FindAllOrdersUseCase usecase.FindAllOrdersUseCase
}

func (s *OrderService) mustEmbedUnimplementedOrderServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, findAllOrdersUseCase usecase.FindAllOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase:   createOrderUseCase,
		FindAllOrdersUseCase: findAllOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, request *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	ordersDTO, err := s.FindAllOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var orders []*pb.Order
	for _, dto := range ordersDTO {
		orders = append(orders, &pb.Order{
			Id:         dto.ID,
			Price:      float32(dto.Price),
			Tax:        float32(dto.Tax),
			FinalPrice: float32(dto.FinalPrice),
		})
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}
