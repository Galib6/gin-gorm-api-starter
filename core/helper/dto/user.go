package dto

import "mime/multipart"

type (
	UserRegisterRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		Role    string `json:"role,omitempty"`
		Picture string `json:"picture,omitempty"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserNameUpdateRequest struct {
		ID   string `json:"id"`
		Name string `json:"name" binding:"required"`
	}

	UserUpdateRequest struct {
		ID       string `json:"id"`
		Name     string `json:"name" form:"name"`
		Email    string `json:"email" form:"email" binding:"omitempty,email"`
		Role     string `json:"role" form:"role" binding:"omitempty,oneof=admin user"`
		Password string `json:"password" form:"password"`
	}

	UserChangePictureRequest struct {
		Picture *multipart.FileHeader `json:"picture" form:"picture"`
	}
)
