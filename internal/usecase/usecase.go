package usecase

import (
	"bytes"
	"context"

	"github.com/patyukin/mbs-reports/internal/db"
)

type MinioClient interface {
	UploadCSVBuffer(ctx context.Context, objectName string, buf *bytes.Buffer) (string, error)
}

type KafkaProducer interface {
	PublishCreditPaymentSolution(ctx context.Context, value []byte) error
}

type UseCase struct {
	db  *db.Registry
	mn  MinioClient
	kfk KafkaProducer
}

func New(db *db.Registry, mn MinioClient, kfk KafkaProducer) *UseCase {
	return &UseCase{
		db:  db,
		mn:  mn,
		kfk: kfk,
	}
}
