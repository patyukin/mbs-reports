package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) ConsumerReportProcess(ctx context.Context, record *kgo.Record) error {
	var transactions []model.Transaction

	log.Debug().Msgf("Received record: %v", string(record.Value))

	if err := json.Unmarshal(record.Value, &transactions); err != nil {
		return fmt.Errorf("failed to unmarshal debezium message: %w", err)
	}

	if err := u.db.GetRepo().InsertIntoTransactions(ctx, transactions); err != nil {
		return fmt.Errorf("failed to insert into reports: %w", err)
	}

	var transactionsStatus []model.TransactionSendStatus
	for _, transaction := range transactions {
		transactionsStatus = append(transactionsStatus, model.TransactionSendStatus{
			ID:         transaction.ID,
			SendStatus: "COMPLETED",
		})
	}

	transactionsStatusBytes, err := json.Marshal(transactionsStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal transactionsStatus: %w", err)
	}

	if err = u.kfk.PublishCreditPaymentSolution(ctx, transactionsStatusBytes); err != nil {
		return fmt.Errorf("failed PublishCreditPaymentSolution: %w", err)
	}

	return nil
}
