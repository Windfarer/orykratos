package sessiontokenexchange

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Exchanger struct {
	ID        uuid.UUID `db:"id" rw:"r"`
	NID       uuid.UUID `db:"nid" rw:"r"`
	FlowID    uuid.UUID `db:"flow_id" rw:"r"`
	SessionID uuid.UUID `db:"session_id" rw:"rw"`
	Code      string    `db:"code" rw:"r"`

	// CreatedAt is a helper struct field for gobuffalo.pop.
	CreatedAt time.Time `db:"created_at"`

	// UpdatedAt is a helper struct field for gobuffalo.pop.
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Exchanger) TableName() string {
	return "session_token_exchangers"
}

type (
	Persister interface {
		CreateSessionTokenExchanger(ctx context.Context, flowID uuid.UUID, code string) error
		GetExchangerFromCode(ctx context.Context, flowID uuid.UUID, code string) (*Exchanger, error)
		UpdateSessionOnExchanger(ctx context.Context, flowID uuid.UUID, sessionID uuid.UUID) error
	}

	PersistenceProvider interface {
		SessionTokenExchangePersister() Persister
	}
)
