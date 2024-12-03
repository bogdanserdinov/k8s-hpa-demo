package v1

import (
	"context"

	subtractionpb "example/gen/go/subtraction/v1"
)

var _ subtractionpb.SubtractionServiceServer = (*server)(nil)

type server struct{}

func New() *server {
	return &server{}
}

func (s *server) Subtract(ctx context.Context, request *subtractionpb.SubtractRequest) (*subtractionpb.SubtractResponse, error) {
	return &subtractionpb.SubtractResponse{
		Rez: request.GetA() - request.GetB(),
	}, nil
}
