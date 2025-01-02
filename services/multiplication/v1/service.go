package v1

import (
	"context"

	multiplicationpb "example/gen/go/multiplication/v1"
	factorialpb "example/gen/go/x/factorial"
	"example/pkg/factorial"
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

func (s *server) Factorial(ctx context.Context, request *factorialpb.FactorialRequest) (*factorialpb.FactorialResponse, error) {
	return &factorialpb.FactorialResponse{
		Res: int64(factorial.Calculate(int(request.GetNum()))),
	}, nil
}
