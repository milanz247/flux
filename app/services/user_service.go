// Package services holds all business logic. Services talk to GORM directly
// (no repository layer), receive context from the controller, and return
// DTOs — never raw models.
package services

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"

	"flux/app/dto"
	"flux/app/models"
	"flux/framework"
)

// UserService manages user CRUD.
type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUsers returns a page of users, optionally filtered by a search term
// matching name or email.
func (s *UserService) GetUsers(ctx context.Context, page, perPage int, search string) ([]dto.UserDTO, dto.PaginationDTO, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	query := s.db.WithContext(ctx).Model(&models.User{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, dto.PaginationDTO{}, err
	}

	var users []models.User
	err := query.
		Order("id DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&users).Error
	if err != nil {
		return nil, dto.PaginationDTO{}, err
	}

	pagination := dto.PaginationDTO{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
	}
	return toUserDTOs(users), pagination, nil
}

// GetUser returns a single user by ID.
func (s *UserService) GetUser(ctx context.Context, id uint) (dto.UserDTO, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return dto.UserDTO{}, err
	}
	return toUserDTO(user), nil
}

// CreateUser creates a user from validated input.
func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserDTO) (dto.UserDTO, error) {
	if err := s.ensureEmailAvailable(ctx, input.Email, 0); err != nil {
		return dto.UserDTO{}, err
	}

	hash, err := framework.HashPassword(input.Password)
	if err != nil {
		return dto.UserDTO{}, err
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}
	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return dto.UserDTO{}, err
	}
	return toUserDTO(user), nil
}

// UpdateUser updates a user's name, email and (optionally) password.
func (s *UserService) UpdateUser(ctx context.Context, id uint, input dto.UpdateUserDTO) (dto.UserDTO, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return dto.UserDTO{}, err
	}

	if err := s.ensureEmailAvailable(ctx, input.Email, id); err != nil {
		return dto.UserDTO{}, err
	}

	user.Name = input.Name
	if user.Email != input.Email {
		user.Email = input.Email
		user.EmailVerifiedAt = nil // changed address must re-verify
	}
	if input.Password != "" {
		hash, err := framework.HashPassword(input.Password)
		if err != nil {
			return dto.UserDTO{}, err
		}
		user.Password = hash
	}

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return dto.UserDTO{}, err
	}
	return toUserDTO(user), nil
}

// DeleteUser soft-deletes a user (gorm.Model's DeletedAt).
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CountUsers returns the total number of users (dashboard stat).
func (s *UserService) CountUsers(ctx context.Context) (int64, error) {
	var total int64
	err := s.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error
	return total, err
}

// CountVerifiedUsers returns how many users verified their email.
func (s *UserService) CountVerifiedUsers(ctx context.Context) (int64, error) {
	var total int64
	err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("email_verified_at IS NOT NULL").Count(&total).Error
	return total, err
}

// CountUsersSince returns users created after the given time.
func (s *UserService) CountUsersSince(ctx context.Context, since time.Time) (int64, error) {
	var total int64
	err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("created_at >= ?", since).Count(&total).Error
	return total, err
}

// RecentUsers returns the newest users (dashboard widget).
func (s *UserService) RecentUsers(ctx context.Context, limit int) ([]dto.UserDTO, error) {
	var users []models.User
	err := s.db.WithContext(ctx).Order("id DESC").Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return toUserDTOs(users), nil
}

func (s *UserService) ensureEmailAvailable(ctx context.Context, email string, excludeID uint) error {
	var count int64
	query := s.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return framework.UnprocessableEntity("The email has already been taken.")
	}
	return nil
}

func toUserDTO(user models.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Verified:  user.EmailVerifiedAt != nil,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04"),
	}
}

func toUserDTOs(users []models.User) []dto.UserDTO {
	dtos := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, toUserDTO(user))
	}
	return dtos
}
