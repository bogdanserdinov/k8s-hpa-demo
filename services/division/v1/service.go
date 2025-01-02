package v1

import (
	"context"

	divisionpb "example/gen/go/division/v1"
	factorialpb "example/gen/go/x/factorial"
	"example/pkg/factorial"
)

var _ divisionpb.DivisionServiceServer = (*server)(nil)

type server struct{}

func New() *server {
	return &server{}
}

func (s *server) Divide(ctx context.Context, request *divisionpb.DivideRequest) (*divisionpb.DivideResponse, error) {
	return &divisionpb.DivideResponse{
		Rez: float32(request.GetA()) / float32(request.GetB()),
	}, nil
}

func (s *server) Factorial(ctx context.Context, request *factorialpb.FactorialRequest) (*factorialpb.FactorialResponse, error) {
	return &factorialpb.FactorialResponse{
		Res: int64(factorial.Calculate(int(request.GetNum()))),
	}, nil
}
