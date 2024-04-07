package models

import (
	"burakozkan138/questionanswerapi/pkg"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RoleType string

const (
	UserRole  RoleType = "user"
	AdminRole RoleType = "admin"
)

type (
	User struct {
		BaseModel
		Fullname            string    `json:"fullname"`
		Username            string    `gorm:"unique" json:"username"`
		Email               string    `gorm:"unique" json:"email"`
		Password            string    `json:"password"`
		Role                RoleType  `gorm:"default:user" json:"role"`
		ProfileImage        string    `gorm:"default:default.jpg" json:"profile_image"`
		Website             string    `json:"website"`
		Location            string    `json:"location"`
		Bio                 string    `json:"bio"`
		Blocked             bool      `gorm:"default:false" json:"blocked"`
		ResetPasswordToken  string    `json:"reset_password_token"`
		ResetPasswordExpire time.Time `json:"reset_password_expire"`
	}
	UserProfile struct {
		BaseModel
		Fullname     string   `json:"fullname"`
		Username     string   `json:"username"`
		Email        string   `json:"email"`
		Role         RoleType `json:"role"`
		ProfileImage string   `json:"profile_image"`
		Website      string   `json:"website"`
		Location     string   `json:"location"`
		Bio          string   `json:"bio"`
		Blocked      bool     `json:"blocked"`
	}

	UserInform struct {
		BaseModel
		Fullname     string   `json:"fullname"`
		Username     string   `json:"username"`
		ProfileImage string   `json:"profile_image"`
		Role         RoleType `json:"role"`
	}

	EditUserValidation struct {
		Fullname string `json:"fullname,omitempty" validate:"omitempty,min=3,max=50"`
		Website  string `json:"website,omitempty" validate:"omitempty,url"`
		Location string `json:"location,omitempty" validate:"omitempty,max=50"`
		Bio      string `json:"bio,omitempty" validate:"omitempty,max=255"`
	}
)

func (u *User) BeforeSave(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		return u.HashPassword(u.Password)
	}

	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.HashPassword(u.Password)
}

func (u *User) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}

	return true
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ToUserProfileResponse() UserProfile {
	return UserProfile{
		BaseModel:    u.BaseModel,
		Fullname:     u.Fullname,
		Username:     u.Username,
		Email:        u.Email,
		Role:         u.Role,
		ProfileImage: u.ProfileImage,
		Website:      u.Website,
		Location:     u.Location,
		Bio:          u.Bio,
		Blocked:      u.Blocked,
	}
}

func (u *User) ToUserInformResponse() UserInform {
	return UserInform{
		BaseModel:    u.BaseModel,
		Fullname:     u.Fullname,
		Username:     u.Username,
		ProfileImage: u.ProfileImage,
		Role:         u.Role,
	}
}

func (u *User) ToUserAuthResponse() (AuthResponse, error) {
	accessToken, err := pkg.CreateAccessToken(u.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, err := pkg.CreateRefreshToken(u.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		User:         u.ToUserProfileResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *User) EditUserCheckFields(editUser EditUserValidation) {
	if len(editUser.Fullname) > 0 {
		u.Fullname = editUser.Fullname
	}

	if len(editUser.Website) > 0 {
		u.Website = editUser.Website
	}

	if len(editUser.Location) > 0 {
		u.Location = editUser.Location
	}

	if len(editUser.Bio) > 0 {
		u.Bio = editUser.Bio
	}
}
