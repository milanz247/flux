package dto

// LoginPageDTO — props for Pages/Auth/Login.vue.
type LoginPageDTO struct {
	Status string `json:"status,omitempty"`
}

// RegisterPageDTO — props for Pages/Auth/Register.vue.
type RegisterPageDTO struct{}

// ForgotPasswordPageDTO — props for Pages/Auth/ForgotPassword.vue.
type ForgotPasswordPageDTO struct {
	Status string `json:"status,omitempty"`
}

// ResetPasswordPageDTO — props for Pages/Auth/ResetPassword.vue.
type ResetPasswordPageDTO struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

// VerifyEmailPageDTO — props for Pages/Auth/VerifyEmail.vue.
type VerifyEmailPageDTO struct {
	Status string `json:"status,omitempty"`
}

// LoginDTO is the validated login input.
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterDTO is the validated registration input.
type RegisterDTO struct {
	Name                 string `json:"name" validate:"required,min=2,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=72"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

// ForgotPasswordDTO is the validated forgot-password input.
type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordDTO is the validated password-reset input.
type ResetPasswordDTO struct {
	Token                string `json:"token" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,max=72"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}
