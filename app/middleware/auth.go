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
		if user := resolveUser(req, app, authService); user != nil {
			req.SetUser(user)
			req.Redirect("/dashboard")
			return
		}
		next()
	}
}

// Identify attaches the authenticated user to the request when a valid
// token is present, but never redirects or rejects — for pages that are
// public but render differently for signed-in users (e.g. the welcome page).
func Identify(app *framework.App, authService *services.AuthService) framework.Middleware {
	return func(req *framework.Request, next func()) {
		if user := resolveUser(req, app, authService); user != nil {
			req.SetUser(user)
		}
		next()
	}
}

// resolveUser looks up the authenticated user for the current request from
// its bearer token or session cookie, or nil if there isn't one — shared by
// Guest and Identify, which differ only in what they do with the result.
func resolveUser(req *framework.Request, app *framework.App, authService *services.AuthService) *framework.AuthUser {
	token := app.Auth().TokenFromRequest(req)
	if token == "" {
		return nil
	}
	claims, err := app.Auth().ParseToken(token, framework.PurposeAccess)
	if err != nil {
		return nil
	}
	user, err := authService.AuthUserByID(req.Context(), claims.UserID)
	if err != nil {
		return nil
	}
	return user
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
		req.Fail(http.StatusUnauthorized, "Unauthenticated.", framework.M{"redirect": "/login"})
		return
	}
	req.Redirect("/login")
}
