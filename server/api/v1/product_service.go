package v1

import (
	"context"

	apiv1 "github.com/thetnaingtn/dirty-hand/proto/gen/api/v1"
	"github.com/thetnaingtn/dirty-hand/store"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *APIV1Service) CreateProduct(ctx context.Context, req *apiv1.CreateProductRequest) (*apiv1.Product, error) {
	prod := &store.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Cover:       req.GetCover(),
	}
	created, err := s.store.CreateProduct(ctx, prod)
	if err != nil {
		return nil, err
	}
	return toProtoProduct(created), nil
}

func (s *APIV1Service) UpdateProduct(ctx context.Context, req *apiv1.UpdateProductRequest) (*apiv1.Product, error) {
	prod := &store.Product{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Cover:       req.GetCover(),
	}
	updated, err := s.store.UpdateProduct(ctx, prod)
	if err != nil {
		return nil, err
	}
	return toProtoProduct(updated), nil
}

func (s *APIV1Service) ListProducts(ctx context.Context, req *apiv1.ListProductsRequest) (*apiv1.ListProductsResponse, error) {
	prods, err := s.store.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	resp := &apiv1.ListProductsResponse{Products: make([]*apiv1.Product, 0, len(prods))}
	for _, p := range prods {
		cp := p // capture range variable
		resp.Products = append(resp.Products, toProtoProduct(cp))
	}
	return resp, nil
}

func (s *APIV1Service) DeleteProduct(ctx context.Context, req *apiv1.DeleteProductRequest) (*emptypb.Empty, error) {
	if err := s.store.DeleteProduct(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func toProtoProduct(p *store.Product) *apiv1.Product {
	if p == nil {
		return nil
	}
	return &apiv1.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Cover:       p.Cover,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
