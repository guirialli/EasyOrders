package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"github.com/guirialli/go-pos/clean_arch/internals/usecase"

	"github.com/guirialli/go-pos/clean_arch/internals/infra/graph/model"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
	dto := usecase.OrderInputDTO{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:         output.ID,
		Price:      float64(output.Price),
		Tax:        float64(output.Tax),
		FinalPrice: float64(output.FinalPrice),
	}, nil
}

// FindAllOrder is the resolver for the findAllOrder field.
func (r *mutationResolver) FindAllOrder(ctx context.Context) ([]*model.Order, error) {
	outputArr, err := r.FindAllOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	orders := make([]*model.Order, len(outputArr))
	for _, output := range outputArr {
		orders = append(orders, &model.Order{
			ID:         output.ID,
			Price:      float64(output.Price),
			Tax:        float64(output.Tax),
			FinalPrice: float64(output.FinalPrice),
		})
	}
	return orders, nil

}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
