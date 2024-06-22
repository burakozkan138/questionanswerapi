package models

type (
	RegisterValidation struct {
		Fullname        string `json:"fullname" validate:"required,min=3,max=50"`
		Username        string `json:"username" validate:"required,min=3,max=50"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	}

	LoginValidation struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Username string `json:"username" validate:"omitempty,min=3,max=50"`
		Password string `json:"password" validate:"required,min=6"`
	}

	ForgotPasswordValidation struct {
		Email string `json:"email" validate:"required,email"`
	}

	ResetPasswordValidation struct {
		Password        string `json:"password" validate:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	}

	AuthResponse struct {
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		User         UserProfile `json:"user"`
	}

	ctxKey string
)

const USER_CTX_KEY ctxKey = "userID"
