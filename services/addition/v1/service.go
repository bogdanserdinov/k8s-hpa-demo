package v1

import (
	"context"

	additionpb "example/gen/go/addition/v1"
	factorialpb "example/gen/go/x/factorial"
	"example/pkg/factorial"
)

var _ additionpb.AdditionServiceServer = (*server)(nil)

type server struct{}

func New() *server {
	return &server{}
}

func (s *server) Add(ctx context.Context, request *additionpb.AddRequest) (*additionpb.AddResponse, error) {
	return &additionpb.AddResponse{
		Rez: request.GetA() + request.GetB(),
	}, nil
}

func (s *server) Factorial(ctx context.Context, request *factorialpb.FactorialRequest) (*factorialpb.FactorialResponse, error) {
	return &factorialpb.FactorialResponse{
		Res: int64(factorial.Calculate(int(request.GetNum()))),
	}, nil
}
