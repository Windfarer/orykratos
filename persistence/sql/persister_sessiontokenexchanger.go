// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ory/kratos/selfservice/sessiontokenexchange"
	"github.com/ory/x/otelx"
	"github.com/ory/x/sqlcon"
)

var _ sessiontokenexchange.Persister = new(Persister)

func (p *Persister) CreateSessionTokenExchanger(ctx context.Context, flowID uuid.UUID, code string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateSessionTokenExchanger")
	defer otelx.End(span, &err)

	e := sessiontokenexchange.Exchanger{
		NID:    p.NetworkID(ctx),
		FlowID: flowID,
		Code:   code,
	}

	e.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	return p.GetConnection(ctx).Create(e)
}

func (p *Persister) GetExchangerFromCode(ctx context.Context, flowID uuid.UUID, code string) (e *sessiontokenexchange.Exchanger, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetExchangerFromCode")
	defer otelx.End(span, &err)

	e = new(sessiontokenexchange.Exchanger)
	conn := p.GetConnection(ctx)
	if err := conn.Where(
		"flow_id = ? AND nid = ? AND code = ? AND session_id IS NOT NULL AND code <> ''",
		flowID, p.NetworkID(ctx), code).First(e); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return e, nil
}

func (p *Persister) UpdateSessionOnExchanger(ctx context.Context, flowID uuid.UUID, sessionID uuid.UUID) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateSessionOnExchanger")
	defer otelx.End(span, &err)

	conn := p.GetConnection(ctx)
	query := fmt.Sprintf("UPDATE %s SET session_id = ? WHERE flow_id = ? AND nid = ?",
		conn.Dialect.Quote(new(sessiontokenexchange.Exchanger).TableName()),
	)

	return sqlcon.HandleError(conn.RawQuery(query, sessionID, flowID, p.NetworkID(ctx)).Exec())
}
