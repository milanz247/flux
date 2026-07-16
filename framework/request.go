package framework

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthUser is the identity of the authenticated user, attached to the
// request by the auth middleware. It intentionally carries only what the
// framework and shared Inertia props need — full user records stay in the
// application's models.
type AuthUser struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	EmailVerifiedAt string `json:"emailVerifiedAt,omitempty"`
}

// Request wraps the underlying HTTP request and response. It is the only
// object controllers interact with — both for reading input and for writing
// responses (see response.go).
type Request struct {
	app  *App
	gin  *gin.Context
	user *AuthUser

	jsonBody     map[string]any
	jsonBodyRead bool
}

func newRequest(app *App, c *gin.Context) *Request {
	return &Request{app: app, gin: c}
}

// --- Input --------------------------------------------------------------

// Input returns a value from the request body — a JSON field or a form
// field — falling back to the query string, mirroring Laravel's $request->input().
func (r *Request) Input(key string) string {
	if strings.Contains(r.gin.ContentType(), "application/json") {
		if v, ok := r.jsonField(key); ok {
			return v
		}
	} else if v, ok := r.gin.GetPostForm(key); ok {
		return v
	}
	return r.gin.Query(key)
}

// Query returns a query-string parameter.
func (r *Request) Query(key string) string { return r.gin.Query(key) }

// QueryDefault returns a query-string parameter or a fallback when absent.
func (r *Request) QueryDefault(key, fallback string) string {
	return r.gin.DefaultQuery(key, fallback)
}

// Param returns a route parameter, e.g. ":id" in "/users/:id".
func (r *Request) Param(key string) string { return r.gin.Param(key) }

// ParamUint returns a numeric route parameter as uint (0 when invalid).
func (r *Request) ParamUint(key string) uint {
	n, err := strconv.ParseUint(r.gin.Param(key), 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

// Header returns a request header value.
func (r *Request) Header(key string) string { return r.gin.GetHeader(key) }

// Method returns the HTTP method.
func (r *Request) Method() string { return r.gin.Request.Method }

// Path returns the request URL path.
func (r *Request) Path() string { return r.gin.Request.URL.Path }

// FullURL returns the path including the query string.
func (r *Request) FullURL() string { return r.gin.Request.URL.RequestURI() }

// Bind decodes the request body (JSON or form, based on Content-Type) into v.
func (r *Request) Bind(v any) error { return r.gin.ShouldBind(v) }

// File returns an uploaded file by form field name.
func (r *Request) File(name string) (*multipart.FileHeader, error) {
	return r.gin.FormFile(name)
}

// SaveFile stores an uploaded file at dst (directories are created by Gin).
func (r *Request) SaveFile(file *multipart.FileHeader, dst string) error {
	return r.gin.SaveUploadedFile(file, dst)
}

// Context returns the request-scoped context for propagation into services
// and GORM queries.
func (r *Request) Context() context.Context { return r.gin.Request.Context() }

// Cookie returns a request cookie value ("" when absent).
func (r *Request) Cookie(name string) string {
	v, err := r.gin.Cookie(name)
	if err != nil {
		return ""
	}
	return v
}

// IsInertia reports whether this request came from the Flux/Inertia client.
func (r *Request) IsInertia() bool { return r.gin.GetHeader("X-Inertia") == "true" }

// WantsJSON reports whether the client expects a JSON response.
func (r *Request) WantsJSON() bool {
	return r.IsInertia() ||
		strings.Contains(r.gin.GetHeader(headerAccept), "application/json") ||
		strings.Contains(r.gin.ContentType(), "application/json")
}

// --- Auth ---------------------------------------------------------------

// SetUser attaches the authenticated user (called by the auth middleware).
func (r *Request) SetUser(user *AuthUser) { r.user = user }

// User returns the authenticated user, or nil for guests.
func (r *Request) User() *AuthUser { return r.user }

// UserID returns the authenticated user's ID (0 for guests).
func (r *Request) UserID() uint {
	if r.user == nil {
		return 0
	}
	return r.user.ID
}

// App exposes the application container (used by application middleware to
// reach the auth manager, DB, logger, etc.).
func (r *Request) App() *App { return r.app }

// --- Validation ---------------------------------------------------------

// Validate binds the request body into dto (a pointer to a DTO struct) and
// runs struct validation. On failure it writes a 422 response containing a
// field→messages error bag — the Flux Vue client picks these up in
// useForm().errors — and returns false. Typical usage:
//
//	var input dto.CreateUserDTO
//	if !req.Validate(&input) {
//	    return
//	}
func (r *Request) Validate(dto any) bool {
	if err := r.Bind(dto); err != nil {
		r.ValidationError(map[string][]string{
			"_body": {"The request body could not be parsed: " + err.Error()},
		})
		return false
	}

	errors := r.app.validator.Struct(dto)
	if len(errors) > 0 {
		r.ValidationError(errors)
		return false
	}
	return true
}

// --- internals ----------------------------------------------------------

const headerAccept = "Accept"

// jsonField lazily parses a JSON body into a map so Input() can read single
// fields without consuming the body for later Bind() calls.
func (r *Request) jsonField(key string) (string, bool) {
	if !r.jsonBodyRead {
		r.jsonBodyRead = true
		body, err := io.ReadAll(r.gin.Request.Body)
		if err == nil && len(body) > 0 {
			_ = json.Unmarshal(body, &r.jsonBody)
			// Restore the body so Bind()/Validate() can still read it.
			r.gin.Request.Body = io.NopCloser(strings.NewReader(string(body)))
		}
	}
	if r.jsonBody == nil {
		return "", false
	}
	value, ok := r.jsonBody[key]
	if !ok {
		return "", false
	}
	switch v := value.(type) {
	case string:
		return v, true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(v), true
	default:
		return "", false
	}
}
