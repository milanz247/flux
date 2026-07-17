// Package routes wires controllers to URLs — the equivalent of Laravel's
// routes/web.php. Dependencies flow constructor-style: App → Service →
// Controller.
package routes

import (
	"flux/app/controllers"
	"flux/app/middleware"
	"flux/app/services"
	"flux/framework"
)

// Register builds the service and controller graph and declares every route.
func Register(app *framework.App) {
	route := app.Router()

	// Services (constructor injection — no repositories, no globals).
	authService := services.NewAuthService(app)

	// Controllers.
	authController := controllers.NewAuthController(authService)
	dashboardController := controllers.NewDashboardController()
	homeController := controllers.NewHomeController()

	// Public marketing page — visible to guests and authenticated users alike,
	// rendering a different nav depending on session state.
	public := route.Group("", middleware.Identify(app, authService))
	public.Get("/", homeController.Welcome)

	// Guest-only pages.
	guest := route.Group("", middleware.Guest(app, authService))
	guest.Get("/login", authController.ShowLogin)
	guest.Post("/login", authController.Login)
	guest.Get("/register", authController.ShowRegister)
	guest.Post("/register", authController.Register)
	guest.Get("/forgot-password", authController.ShowForgotPassword)
	guest.Post("/forgot-password", authController.ForgotPassword)
	guest.Get("/reset-password", authController.ShowResetPassword)
	guest.Post("/reset-password", authController.ResetPassword)

	// Email verification links arrive from the mail client, signed-in or not.
	route.Get("/verify-email", authController.VerifyEmail)

	// Authenticated app.
	auth := route.Group("", middleware.Auth(app, authService))
	auth.Get("/dashboard", dashboardController.Index)
	auth.Post("/logout", authController.Logout)
	auth.Post("/verify-email/resend", authController.ResendVerification)
}
