package v1

import (
	"context"

	divisionpb "example/gen/go/division/v1"
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
