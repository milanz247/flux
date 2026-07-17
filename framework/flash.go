package framework

import (
	"encoding/base64"
	"encoding/json"
)

// Flash is a one-shot message carried across a redirect in a short-lived
// cookie — the Flux equivalent of Laravel's session flash. The next View
// exposes it to Vue as the shared "flash" prop (useFlash() client-side) and
// clears it, so it renders exactly once and never pollutes the URL the way
// ?status= query params do.
type Flash struct {
	Type    string `json:"type"` // "success", "error", ... — the client decides styling
	Message string `json:"message"`
}

const flashCookie = "flux_flash"

// RedirectWith flashes a one-shot message and redirects, e.g.:
//
//	req.RedirectWith("/login", "success", "Your password has been reset.")
func (r *Request) RedirectWith(url, flashType, message string) {
	if payload, err := json.Marshal(Flash{Type: flashType, Message: message}); err == nil {
		r.gin.SetCookie(flashCookie, base64.RawURLEncoding.EncodeToString(payload),
			60, "/", "", secureCookies(r.app.config), true)
	}
	r.Redirect(url)
}

// consumeFlash returns the pending flash message, clearing its cookie so the
// message shows at most once.
func (r *Request) consumeFlash() *Flash {
	raw := r.Cookie(flashCookie)
	if raw == "" {
		return nil
	}
	r.gin.SetCookie(flashCookie, "", -1, "/", "", secureCookies(r.app.config), true)

	payload, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil
	}
	var flash Flash
	if err := json.Unmarshal(payload, &flash); err != nil {
		return nil
	}
	return &flash
}
