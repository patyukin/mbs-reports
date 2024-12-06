package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/proto/error_v1"
	reportpb "github.com/patyukin/mbs-pkg/pkg/proto/report_v1"
)

type UseCase interface {
	GetUserReportUseCase(ctx context.Context, in *reportpb.GetUserReportRequest) (*reportpb.GetUserReportResponse, error)
}

type Server struct {
	reportpb.UnimplementedReportServiceServer
	uc UseCase
}

func New(uc UseCase) *Server {
	return &Server{
		uc: uc,
	}
}

func (s *Server) GetUserReport(ctx context.Context, in *reportpb.GetUserReportRequest) (*reportpb.GetUserReportResponse, error) {
	response, err := s.uc.GetUserReportUseCase(ctx, in)
	if err != nil {
		return &reportpb.GetUserReportResponse{
			Error: &error_v1.ErrorResponse{
				Code:        500,
				Message:     "Internal Server Error",
				Description: fmt.Sprintf("failed uc.GetUserReportUseCase: %v", err),
			},
		}, nil
	}

	if response.Error != nil {
		return &reportpb.GetUserReportResponse{Error: response.Error}, nil
	}

	return response, nil
}
