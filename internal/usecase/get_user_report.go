package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patyukin/mbs-pkg/pkg/proto/error_v1"
	reportpb "github.com/patyukin/mbs-pkg/pkg/proto/report_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) GetUserReportUseCase(ctx context.Context, in *reportpb.GetUserReportRequest) (*reportpb.GetUserReportResponse, error) {
	reports, err := u.db.GetRepo().SelectReportsByUserID(ctx, in)
	if err != nil {
		return &reportpb.GetUserReportResponse{
			Error: &error_v1.ErrorResponse{
				Code:        500,
				Message:     "Internal Server Error",
				Description: fmt.Sprintf("failed to select reports by user id: %v", err),
			},
		}, nil
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"amount", "currency", "description", "payment_description", "status", "created_at"}
	if err = writer.Write(headers); err != nil {
		return &reportpb.GetUserReportResponse{
			Error: &error_v1.ErrorResponse{
				Code:        500,
				Message:     "Internal Server Error",
				Description: fmt.Sprintf("failed to write headers to CSV: %v", err),
			},
		}, nil
	}

	for _, r := range reports {
		am := float64(r.Amount / 100)
		record := []string{
			fmt.Sprintf("%.2f", am),
			r.Currency,
			r.Description,
			r.PaymentDescription,
			r.Status,
			r.CreatedAt,
		}

		if err = writer.Write(record); err != nil {
			return &reportpb.GetUserReportResponse{
				Error: &error_v1.ErrorResponse{
					Code:        500,
					Message:     "Internal Server Error",
					Description: fmt.Sprintf("failed to write record to CSV: %v", err),
				},
			}, nil
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return &reportpb.GetUserReportResponse{
			Error: &error_v1.ErrorResponse{
				Code:        500,
				Message:     "Internal Server Error",
				Description: fmt.Sprintf("failed to flush writer: %v", err),
			},
		}, nil
	}

	now := time.Now()
	objectName := fmt.Sprintf(
		"%04d/%02d/%02d-%s.csv",
		now.Year(),
		int(now.Month()),
		now.Day(),
		uuid.New().String(),
	)

	fileUrl, err := u.mn.UploadCSVBuffer(ctx, objectName, &buf)
	if err != nil {
		return &reportpb.GetUserReportResponse{
			Error: &error_v1.ErrorResponse{
				Code:        500,
				Message:     "Internal Server Error",
				Description: fmt.Sprintf("failed u.mn.UploadCSVBuffer: %v", err),
			},
		}, nil
	}

	log.Debug().Msgf("fileUrl: %v", fileUrl)

	return &reportpb.GetUserReportResponse{Message: fileUrl}, nil
}
