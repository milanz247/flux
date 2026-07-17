package framework

import (
	"errors"
	"log/slog"
	"maps"
	"net/http"

	"gorm.io/gorm"
)

// M is shorthand for an arbitrary JSON object, like gin.H:
//
//	req.JSON(200, framework.M{"ok": true})
type M = map[string]any

// HTTPError is an error carrying an HTTP status code. Services return these
// (via framework.NewError and friends) and controllers pass them straight to
// req.Error(), which renders the right status.
type HTTPError struct {
	Status  int
	Message string
	cause   error
}

func (e *HTTPError) Error() string { return e.Message }
func (e *HTTPError) Unwrap() error { return e.cause }

// NewError creates an HTTPError with an explicit status.
func NewError(status int, message string) *HTTPError {
	return &HTTPError{Status: status, Message: message}
}

// WrapError attaches a status/message to an underlying error.
func WrapError(status int, message string, cause error) *HTTPError {
	return &HTTPError{Status: status, Message: message, cause: cause}
}

// NotFound is a 404 error.
func NotFound(message string) *HTTPError { return NewError(http.StatusNotFound, message) }

// Unauthorized is a 401 error.
func Unauthorized(message string) *HTTPError { return NewError(http.StatusUnauthorized, message) }

// Forbidden is a 403 error.
func Forbidden(message string) *HTTPError { return NewError(http.StatusForbidden, message) }

// UnprocessableEntity is a 422 error.
func UnprocessableEntity(message string) *HTTPError {
	return NewError(http.StatusUnprocessableEntity, message)
}

// --- Response methods on *Request ----------------------------------------

// View renders an Inertia page. The component name maps directly to a Vue
// file: "Users/Index" renders resources/js/Pages/Users/Index.vue, with the
// DTO's fields injected as props. Shared props (auth, appName, flash) are
// merged in automatically.
func (r *Request) View(component string, props any) {
	pageProps := structToMap(props)

	// Shared props, available to every page (usePage() on the Vue side).
	pageProps["auth"] = M{"user": r.user}
	pageProps["appName"] = r.app.config.App.Name
	pageProps["flash"] = r.consumeFlash()

	r.app.inertia.Render(r.gin, component, pageProps)
}

// JSON writes a raw JSON response with the given status code.
func (r *Request) JSON(status int, data any) {
	r.gin.JSON(status, data)
}

// Success writes a 200 JSON envelope: {"success": true, "data": ...}.
func (r *Request) Success(data any) {
	r.gin.JSON(http.StatusOK, M{"success": true, "data": data})
}

// Fail writes the standard failure envelope: {"success": false, "message":
// ...}. Optional extra maps merge additional fields into the envelope —
// every JSON error Flux emits goes through here so clients can rely on one
// shape.
func (r *Request) Fail(status int, message string, extra ...M) {
	payload := M{"success": false, "message": message}
	for _, m := range extra {
		maps.Copy(payload, m)
	}
	r.gin.JSON(status, payload)
}

// Error renders an error response. HTTPError values keep their status;
// gorm.ErrRecordNotFound becomes 404; anything else is a 500 (logged, and
// the message is hidden from clients unless APP_DEBUG is on).
func (r *Request) Error(err error) {
	status := http.StatusInternalServerError
	message := "Something went wrong."

	var httpErr *HTTPError
	switch {
	case errors.As(err, &httpErr):
		status = httpErr.Status
		message = httpErr.Message
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
		message = "Resource not found."
	default:
		r.app.logger.Error("unhandled error",
			slog.String("error", err.Error()),
			slog.String("path", r.Path()),
			slog.String("method", r.Method()),
		)
		if r.app.config.App.Debug {
			message = err.Error()
		}
	}

	r.Fail(status, message)
}

// ValidationError writes a 422 with a field→messages error bag in the shape
// the Flux Vue client (useForm) understands.
func (r *Request) ValidationError(errors map[string][]string) {
	r.Fail(http.StatusUnprocessableEntity, "The given data was invalid.", M{"errors": errors})
}

// Redirect sends the client to another URL. Non-GET requests redirect with
// 303 See Other so the follow-up is always a GET — required for Inertia
// form submissions to behave like Laravel.
func (r *Request) Redirect(url string) {
	status := http.StatusFound
	if r.Method() != http.MethodGet {
		status = http.StatusSeeOther
	}
	r.gin.Redirect(status, url)
}

// Download sends a file as an attachment with the given filename.
func (r *Request) Download(path, filename string) {
	r.gin.FileAttachment(path, filename)
}

// ServeFile streams a file inline (images, PDFs viewed in the browser).
// Note: req.File(name) reads an *uploaded* file; ServeFile *sends* one.
func (r *Request) ServeFile(path string) {
	r.gin.File(path)
}

// NoContent writes an empty 204 response.
func (r *Request) NoContent() {
	r.gin.Status(http.StatusNoContent)
}
