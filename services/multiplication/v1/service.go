package v1

import (
	"context"

	multiplicationpb "example/gen/go/multiplication/v1"
)

var _ multiplicationpb.MultiplicationServiceServer = (*server)(nil)

type server struct{}

func New() *server {
	return &server{}
}

func (s *server) Multiply(ctx context.Context, request *multiplicationpb.MultiplyRequest) (*multiplicationpb.MultiplyResponse, error) {
	return &multiplicationpb.MultiplyResponse{
		Rez: int64(request.GetA()) * int64(request.GetB()),
	}, nil
}
