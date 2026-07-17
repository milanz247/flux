package controllers

import (
	"flux/app/dto"
	"flux/app/services"
	"flux/framework"
)

// AuthController serves the full authentication flow: login, registration,
// logout, forgot/reset password and email verification.
type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// ShowLogin — GET /login
func (c *AuthController) ShowLogin(req *framework.Request) {
	req.View("Auth/Login", dto.LoginPageDTO{Status: req.Query("status")})
}

// Login — POST /login
func (c *AuthController) Login(req *framework.Request) {
	var input dto.LoginDTO
	if !req.Validate(&input) {
		return
	}

	result, err := c.service.Login(req.Context(), input)
	if err != nil {
		req.Error(err)
		return
	}

	req.App().Auth().StartSession(req, result.Token)
	req.Redirect("/dashboard")
}

// ShowRegister — GET /register
func (c *AuthController) ShowRegister(req *framework.Request) {
	req.View("Auth/Register", dto.RegisterPageDTO{})
}

// Register — POST /register
func (c *AuthController) Register(req *framework.Request) {
	var input dto.RegisterDTO
	if !req.Validate(&input) {
		return
	}

	result, err := c.service.Register(req.Context(), input)
	if err != nil {
		req.Error(err)
		return
	}

	req.App().Auth().StartSession(req, result.Token)
	req.Redirect("/dashboard")
}

// Logout — POST /logout
func (c *AuthController) Logout(req *framework.Request) {
	req.App().Auth().EndSession(req)
	req.RedirectWith("/login", "success", "You have been signed out.")
}

// ShowForgotPassword — GET /forgot-password
func (c *AuthController) ShowForgotPassword(req *framework.Request) {
	req.View("Auth/ForgotPassword", dto.ForgotPasswordPageDTO{Status: req.Query("status")})
}

// ForgotPassword — POST /forgot-password
func (c *AuthController) ForgotPassword(req *framework.Request) {
	var input dto.ForgotPasswordDTO
	if !req.Validate(&input) {
		return
	}

	if err := c.service.ForgotPassword(req.Context(), input); err != nil {
		req.Error(err)
		return
	}

	req.Redirect("/forgot-password?status=sent")
}

// ShowResetPassword — GET /reset-password?token=...&email=...
func (c *AuthController) ShowResetPassword(req *framework.Request) {
	req.View("Auth/ResetPassword", dto.ResetPasswordPageDTO{
		Token: req.Query("token"),
		Email: req.Query("email"),
	})
}

// ResetPassword — POST /reset-password
func (c *AuthController) ResetPassword(req *framework.Request) {
	var input dto.ResetPasswordDTO
	if !req.Validate(&input) {
		return
	}

	if err := c.service.ResetPassword(req.Context(), input); err != nil {
		req.Error(err)
		return
	}

	req.Redirect("/login?status=password-reset")
}

// VerifyEmail — GET /verify-email?token=...
// Without a token it renders the "check your inbox" page.
func (c *AuthController) VerifyEmail(req *framework.Request) {
	token := req.Query("token")
	if token == "" {
		req.View("Auth/VerifyEmail", dto.VerifyEmailPageDTO{Status: req.Query("status")})
		return
	}

	if err := c.service.VerifyEmail(req.Context(), token); err != nil {
		req.Error(err)
		return
	}

	req.Redirect("/dashboard")
}

// ResendVerification — POST /verify-email/resend (auth required)
func (c *AuthController) ResendVerification(req *framework.Request) {
	if err := c.service.SendVerificationEmail(req.Context(), req.UserID()); err != nil {
		req.Error(err)
		return
	}
	req.Redirect("/verify-email?status=sent")
}
