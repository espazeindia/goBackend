package entities

import "time"

type OperationalGuyLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type OperationalGuyLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type OperationalGuyRegistrationRequest struct {
	Email                  string `json:"email" binding:"required,email"`
	Password               string `json:"password" binding:"required,min=6"`
	Name                   string `json:"name" binding:"required,min=2"`
	PhoneNumber            string `json:"phoneNumber" binding:"required,min=10"`
	Address                string `json:"address" binding:"required,min=10"`
	EmergencyContactNumber string `json:"emergencyContactNumber" binding:"required,min=10"`
}

type OperationalGuyRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type OperationalGuy struct {
	ID                     string    `json:"id" bson:"_id,omitempty"`
	Email                  string    `json:"email" bson:"email"`
	Password               string    `json:"-" bson:"password"` // "-" means this field won't be included in JSON responses
	Name                   string    `json:"name" bson:"name"`
	IsFirstLogin           bool      `json:"isFirstLogin" bson:"isFirstLogin"`
	PhoneNumber            string    `json:"phoneNumber" bson:"phoneNumber"`
	Address                string    `json:"address" bson:"address"`
	EmergencyContactNumber string    `json:"emergencyContactNumber" bson:"emergencyContactNumber"`
	CreatedAt              time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt" bson:"updatedAt"`
	LastLoginAt            time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
}

// Seller entities
type SellerLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SellerLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SellerRegistrationRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Name         string `json:"name" binding:"required,min=2"`
	PhoneNumber  string `json:"phoneNumber" binding:"required,min=10"`
	Address      string `json:"address" binding:"required,min=10"`
	BusinessName string `json:"businessName" binding:"required,min=2"`
	BusinessType string `json:"businessType" binding:"required"`
}

type SellerRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Seller struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"-" bson:"password"`
	Name         string    `json:"name" bson:"name"`
	IsFirstLogin bool      `json:"isFirstLogin" bson:"isFirstLogin"`
	PhoneNumber  string    `json:"phoneNumber" bson:"phoneNumber"`
	Address      string    `json:"address" bson:"address"`
	BusinessName string    `json:"businessName" bson:"businessName"`
	BusinessType string    `json:"businessType" bson:"businessType"`
	IsVerified   bool      `json:"isVerified" bson:"isVerified"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
	LastLoginAt  time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
}

// Customer entities
type CustomerLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CustomerLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CustomerRegistrationRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required,min=2"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
	Address     string `json:"address" binding:"required,min=10"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"`
}

type CustomerRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Customer struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"-" bson:"password"`
	Name         string    `json:"name" bson:"name"`
	IsFirstLogin bool      `json:"isFirstLogin" bson:"isFirstLogin"`
	PhoneNumber  string    `json:"phoneNumber" bson:"phoneNumber"`
	Address      string    `json:"address" bson:"address"`
	DateOfBirth  string    `json:"dateOfBirth" bson:"dateOfBirth"`
	IsVerified   bool      `json:"isVerified" bson:"isVerified"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
	LastLoginAt  time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
}
