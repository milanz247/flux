package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"flux/app/dto"
	"flux/app/models"
	"flux/framework"
)

// AuthService implements registration, login, password reset and email
// verification on top of the framework's JWT auth manager.
type AuthService struct {
	db     *gorm.DB
	auth   *framework.AuthManager
	mailer framework.Mailer
	appURL string
}

func NewAuthService(app *framework.App) *AuthService {
	return &AuthService{
		db:     app.DB(),
		auth:   app.Auth(),
		mailer: app.Mailer(),
		appURL: app.Config().App.URL,
	}
}

// AuthResult is what a successful login/registration yields.
type AuthResult struct {
	User  dto.UserDTO
	Token string
}

// Register creates the account, sends the verification mail and signs the
// user in.
func (s *AuthService) Register(ctx context.Context, input dto.RegisterDTO) (AuthResult, error) {
	var existing int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", input.Email).Count(&existing).Error; err != nil {
		return AuthResult{}, err
	}
	if existing > 0 {
		return AuthResult{}, framework.UnprocessableEntity("The email has already been taken.")
	}

	hash, err := framework.HashPassword(input.Password)
	if err != nil {
		return AuthResult{}, err
	}

	user := models.User{Name: input.Name, Email: input.Email, Password: hash}
	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return AuthResult{}, err
	}

	if err := s.SendVerificationEmail(ctx, user.ID); err != nil {
		// Non-fatal: the user can request a new link from the verify page.
		slog.Warn("failed to send verification email", slog.String("error", err.Error()))
	}

	return s.issueSession(user)
}

// Login verifies credentials and issues a JWT session.
func (s *AuthService) Login(ctx context.Context, input dto.LoginDTO) (AuthResult, error) {
	var user models.User
	err := s.db.WithContext(ctx).Where("email = ?", input.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || (err == nil && !framework.CheckPassword(user.Password, input.Password)) {
		return AuthResult{}, framework.UnprocessableEntity("These credentials do not match our records.")
	}
	if err != nil {
		return AuthResult{}, err
	}
	return s.issueSession(user)
}

// SendVerificationEmail emails a signed verification link to the user.
func (s *AuthService) SendVerificationEmail(ctx context.Context, userID uint) error {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	if user.EmailVerifiedAt != nil {
		return nil
	}

	token, err := s.auth.IssueEmailVerifyToken(user.ID, user.Email)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/verify-email?token=%s", s.appURL, token)
	body := fmt.Sprintf(
		"Hi %s,\n\nPlease verify your email address by opening the link below:\n\n%s\n\nThis link expires in 24 hours.\n\n— %s",
		user.Name, link, "Flux",
	)
	return s.mailer.Send(user.Email, "Verify your email address", body)
}

// VerifyEmail validates a verification token and marks the address verified.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	claims, err := s.auth.ParseToken(token, framework.PurposeEmailVerify)
	if err != nil {
		return framework.UnprocessableEntity("This verification link is invalid or has expired.")
	}

	var user models.User
	if err := s.db.WithContext(ctx).First(&user, claims.UserID).Error; err != nil {
		return err
	}
	// The link only counts for the address it was issued for.
	if user.Email != claims.Email {
		return framework.UnprocessableEntity("This verification link is invalid or has expired.")
	}
	if user.EmailVerifiedAt != nil {
		return nil
	}

	now := time.Now()
	return s.db.WithContext(ctx).Model(&user).Update("email_verified_at", &now).Error
}

// ForgotPassword emails a signed reset link. It intentionally succeeds even
// for unknown addresses so the endpoint cannot be used to probe accounts.
func (s *AuthService) ForgotPassword(ctx context.Context, input dto.ForgotPasswordDTO) error {
	var user models.User
	err := s.db.WithContext(ctx).Where("email = ?", input.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	token, err := s.auth.IssuePasswordResetToken(user.ID, user.Email)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/reset-password?token=%s&email=%s", s.appURL, token, user.Email)
	body := fmt.Sprintf(
		"Hi %s,\n\nYou requested a password reset. Open the link below to choose a new password:\n\n%s\n\nThis link expires in 1 hour. If you did not request this, no action is needed.\n\n— %s",
		user.Name, link, "Flux",
	)
	return s.mailer.Send(user.Email, "Reset your password", body)
}

// ResetPassword validates a reset token and stores the new password.
func (s *AuthService) ResetPassword(ctx context.Context, input dto.ResetPasswordDTO) error {
	claims, err := s.auth.ParseToken(input.Token, framework.PurposePasswordReset)
	if err != nil {
		return framework.UnprocessableEntity("This password reset link is invalid or has expired.")
	}

	var user models.User
	if err := s.db.WithContext(ctx).First(&user, claims.UserID).Error; err != nil {
		return err
	}
	if user.Email != claims.Email {
		return framework.UnprocessableEntity("This password reset link is invalid or has expired.")
	}

	hash, err := framework.HashPassword(input.Password)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Model(&user).Update("password", hash).Error
}

// AuthUserByID loads the framework identity for the auth middleware.
func (s *AuthService) AuthUserByID(ctx context.Context, id uint) (*framework.AuthUser, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	authUser := &framework.AuthUser{ID: user.ID, Name: user.Name, Email: user.Email}
	if user.EmailVerifiedAt != nil {
		authUser.EmailVerifiedAt = user.EmailVerifiedAt.Format(time.RFC3339)
	}
	return authUser, nil
}

func (s *AuthService) issueSession(user models.User) (AuthResult, error) {
	token, err := s.auth.IssueAccessToken(user.ID, user.Email)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: toUserDTO(user), Token: token}, nil
}
