package sessiontokenexchange

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/ory/herodot"
	"github.com/ory/kratos/driver/config"
	"github.com/ory/kratos/session"
	"github.com/ory/kratos/x"
)

const (
	RouteExchangeCodeForSessionToken = "/self-service/exchange-code-for-session-token"
)

type (
	handlerDependencies interface {
		PersistenceProvider
		config.Provider
		x.WriterProvider
		session.PersistenceProvider
	}

	HandlerProvider interface {
		SessionTokenExchangeHandler() *Handler
	}
	Handler struct {
		d handlerDependencies
	}
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{d: d}
}

func (h *Handler) RegisterPublicRoutes(public *x.RouterPublic) {
	public.GET(RouteExchangeCodeForSessionToken, h.exchangeCode)
}

func (h *Handler) RegisterAdminRoutes(admin *x.RouterAdmin) {
	admin.GET(RouteExchangeCodeForSessionToken, x.RedirectToPublicRoute(h.d))
}

// Exchange Session Token Parameters
//
// swagger:parameters getLoginFlow
//
//nolint:deadcode,unused
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type exchangeSessionToken struct {
	// The Login Flow ID
	//
	// The value for this parameter comes from `flow` URL Query parameter sent to your
	// application (e.g. `/login?flow=abcde`).
	//
	// required: true
	// in: query
	ID string `json:"id"`

	// The Session Token Exchange Code
	//
	// required: true
	// in: query
	SessionTokenExchangeCode string `json:"code"`
}

// The Response for Registration Flows via API
//
// swagger:model successfulCodeExchangeResponse
type codeExchangeResponse struct {
	// The Session Token
	//
	// A session token is equivalent to a session cookie, but it can be sent in the HTTP Authorization
	// Header:
	//
	// 		Authorization: bearer ${session-token}
	//
	// The session token is only issued for API flows, not for Browser flows!
	Token string `json:"session_token,omitempty"`

	// The Session
	//
	// The session contains information about the user, the session device, and so on.
	// This is only available for API flows, not for Browser flows!
	//
	// required: true
	Session *session.Session `json:"session"`
}

// swagger:route GET /self-service/login/exchange-session-token frontend exchangeSessionToken
//
// # Exchange Session Token
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: successfulNativeLogin
//	  403: errorGeneric
//	  404: errorGeneric
//	  410: errorGeneric
//	  default: errorGeneric
func (h *Handler) exchangeCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	flowID := r.URL.Query().Get("id")
	code := r.URL.Query().Get("code")
	ctx := r.Context()

	if code == "" || flowID == "" {
		h.d.Writer().WriteError(w, r, herodot.ErrBadRequest.WithReason(`both "id" and "code" query params must be set`))
		return
	}

	flowUUID, err := uuid.FromString(flowID)
	if err != nil {
		h.d.Writer().WriteError(w, r, herodot.ErrBadRequest.WithReason(`"id" must be a UUID`))
		return
	}

	e, err := h.d.SessionTokenExchangePersister().GetExchangerFromCode(ctx, flowUUID, code)
	if err != nil {
		h.d.Writer().WriteError(w, r, herodot.ErrNotFound.WithReason(`no session yet for "id" and "code" combination`))
		return
	}

	sess, err := h.d.SessionPersister().GetSession(ctx, e.SessionID, session.ExpandDefault)
	if err != nil {
		h.d.Writer().WriteError(w, r, err)
		return
	}

	h.d.Writer().Write(w, r, &codeExchangeResponse{
		Token:   sess.Token,
		Session: sess,
	})
}
