package additional

import (
	"context"
	additionpb "example/gen/go/additional/v1"
)

var _ additionpb.AdditionServiceServer = (*server)(nil)

type server struct{}

func New() *server {
	return &server{}
}

func (s server) Add(ctx context.Context, request *additionpb.AddRequest) (*additionpb.AddResponse, error) {
	return &additionpb.AddResponse{
		Rez: request.GetA() + request.GetB(),
	}, nil
}
