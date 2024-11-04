package pghelper

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func ToPGUUID(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: uuid, Valid: true}
}

func ToPGTimestamp(time time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  time,
		Valid: true,
	}
}
