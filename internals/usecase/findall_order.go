package usecase

import "github.com/guirialli/go-pos/clean_arch/internals/entity"

type OrderFindAllOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type FindAllOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewFindAllOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (u *FindAllOrdersUseCase) Execute() ([]OrderFindAllOutputDTO, error) {
	orders, err := u.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var outputDTOs []OrderFindAllOutputDTO
	for _, order := range orders {
		outputDTOs = append(outputDTOs, OrderFindAllOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return outputDTOs, nil
}
