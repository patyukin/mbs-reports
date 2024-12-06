package db

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/model"
	reportpb "github.com/patyukin/mbs-pkg/pkg/proto/report_v1"
	"time"
)

type Repository struct {
	db QueryExecutor
}

func (r *Repository) InsertIntoTransactions(ctx context.Context, transactions []model.Transaction) error {
	query := `
INSERT INTO transactions
    (id, payment_id, account_id, user_id, type, amount, currency, description, payment_description, status, send_status, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	for _, transaction := range transactions {
		var description string
		if transaction.Description.Valid {
			description = transaction.Description.String
		}

		createdAt, err := time.Parse(time.RFC3339Nano, transaction.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to parse time: %w", err)
		}

		_, err = r.db.ExecContext(
			ctx,
			query,
			transaction.ID,
			transaction.PaymentID,
			transaction.AccountID,
			transaction.UserID,
			transaction.Type,
			transaction.Amount,
			transaction.Currency,
			description,
			transaction.PaymentDescription,
			transaction.Status,
			transaction.SendStatus,
			createdAt.Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			return fmt.Errorf("failed r.db.ExecContext: %w", err)
		}

	}

	return nil
}

func (r *Repository) SelectReportsByUserID(_ context.Context, in *reportpb.GetUserReportRequest) ([]model.TransactionReport, error) {
	start := in.StartDate
	end := in.EndDate

	query := `
SELECT 
	id, amount, currency, description, payment_description, status, send_status, created_at
FROM transactions 
WHERE user_id = ? AND created_at BETWEEN ? AND ?
ORDER BY created_at
`

	rows, err := r.db.Query(query, in.UserId, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.Query: %w", err)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed rows.Err(): %w", err)
	}

	var reports []model.TransactionReport
	for rows.Next() {
		var report model.TransactionReport
		if err = rows.Scan(
			&report.ID,
			&report.Amount,
			&report.Currency,
			&report.Description,
			&report.PaymentDescription,
			&report.Status,
			&report.SendStatus,
			&report.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed rows.Scan(): %w", err)
		}

		reports = append(reports, report)
	}

	return reports, nil
}
