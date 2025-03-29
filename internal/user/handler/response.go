package handler

import "roulette/internal/user/model"

type UserResponse struct {
	*model.User
}

func NewUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		user,
	}
}

type UserProfileResponse struct {
	*model.UserProfile
}

func NewUserProfileResponse(userProfile *model.UserProfile) *UserProfileResponse {
	return &UserProfileResponse{
		userProfile,
	}
}
