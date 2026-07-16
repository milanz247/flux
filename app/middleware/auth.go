// Package middleware contains application middleware built on the framework
// Middleware type: func(req *framework.Request, next func()).
package middleware

import (
	"net/http"

	"flux/app/services"
	"flux/framework"
)

// Auth requires a valid JWT (Bearer header or session cookie), loads the
// user and attaches it to the request. Guests are redirected to /login
// (browser) or given a 401 (API/Inertia XHR without a page).
func Auth(app *framework.App, authService *services.AuthService) framework.Middleware {
	return func(req *framework.Request, next func()) {
		token := app.Auth().TokenFromRequest(req)
		if token == "" {
			rejectGuest(req)
			return
		}

		claims, err := app.Auth().ParseToken(token, framework.PurposeAccess)
		if err != nil {
			app.Auth().EndSession(req)
			rejectGuest(req)
			return
		}

		user, err := authService.AuthUserByID(req.Context(), claims.UserID)
		if err != nil {
			app.Auth().EndSession(req)
			rejectGuest(req)
			return
		}

		req.SetUser(user)
		next()
	}
}

// Guest redirects already-authenticated users away from guest-only pages
// (login, register) to the dashboard.
func Guest(app *framework.App, authService *services.AuthService) framework.Middleware {
	return func(req *framework.Request, next func()) {
		token := app.Auth().TokenFromRequest(req)
		if token != "" {
			if claims, err := app.Auth().ParseToken(token, framework.PurposeAccess); err == nil {
				if user, err := authService.AuthUserByID(req.Context(), claims.UserID); err == nil {
					req.SetUser(user)
					req.Redirect("/dashboard")
					return
				}
			}
		}
		next()
	}
}

// Verified requires the authenticated user to have a verified email address.
// Apply after Auth.
func Verified() framework.Middleware {
	return func(req *framework.Request, next func()) {
		user := req.User()
		if user == nil || user.EmailVerifiedAt == "" {
			req.Redirect("/verify-email")
			return
		}
		next()
	}
}

func rejectGuest(req *framework.Request) {
	// Inertia navigation and API clients get a 401; full-page browser
	// requests are redirected to the login screen.
	if req.IsInertia() || req.WantsJSON() {
		req.JSON(http.StatusUnauthorized, map[string]any{
			"success":  false,
			"message":  "Unauthenticated.",
			"redirect": "/login",
		})
		return
	}
	req.Redirect("/login")
}
