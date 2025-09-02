package entities

import "time"

type OperationalGuy struct {
	OperationalGuyID       string     `json:"id" bson:"_id,omitempty"`
	Email                  string     `json:"email" bson:"email"`
	Password               string     `json:"-" bson:"password"` // "-" means this field won't be included in JSON responses
	Name                   string     `json:"name" bson:"name"`
	PhoneNumber            string     `json:"phoneNumber" bson:"phoneNumber"`
	Address                string     `json:"address" bson:"address"`
	EmergencyContactNumber string     `json:"emergencyContactNumber" bson:"emergencyContactNumber"`
	CreatedAt              time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt" bson:"updatedAt"`
	LastLoginAt            *time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
}

type OperationalGuyRegistrationRequest struct {
	Email                  string `json:"email" binding:"required,email"`
	Password               string `json:"password" binding:"required,min=5"`
	Name                   string `json:"name" binding:"required,min=2"`
	PhoneNumber            string `json:"phoneNumber" binding:"required,min=10"`
	Address                string `json:"address" binding:"required,min=10"`
	EmergencyContactNumber string `json:"emergencyContactNumber" binding:"required,min=10"`
}

type OperationalGuyLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OperationalGuyLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type OperationalGuyRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type Seller struct {
	SellerID           string     `json:"id" bson:"_id,omitempty"`
	Name               string     `json:"name" bson:"name"`
	PhoneNumber        string     `json:"phoneNumber" bson:"phoneNumber,min=10"`
	Address            string     `json:"address" bson:"address"`
	OTP                int        `json:"otp" bson:"otp,min=6"`
	NumberOfRetriesOTP int        `json:"numberOfRetriesOTP" bson:"numberOfRetriesOTP"`
	OTPGeneratedAt     time.Time  `json:"otpGeneratedAt" bson:"otpGeneratedAt"`
	PIN                int        `json:"pin" bson:"pin,min=6"`
	NumberOfRetriesPIN int        `json:"numberOfRetriesPIN" bson:"numberOfRetriesPIN"`
	LastLoginAt        *time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
	StoreID            string     `json:"storeId" bson:"storeId"`
}

type Customer struct {
	CustomerID         string     `json:"id" bson:"_id,omitempty"`
	Name               string     `json:"name" bson:"name"`
	PhoneNumber        string     `json:"phoneNumber" bson:"phoneNumber"`
	OTP                int        `json:"otp" bson:"otp"`
	NumberOfRetriesOTP int        `json:"numberOfRetriesOTP" bson:"numberOfRetriesOTP"`
	OTPGeneratedAt     time.Time  `json:"otpGeneratedAt" bson:"otpGeneratedAt"`
	PIN                int        `json:"pin" bson:"pin"`
	NumberOfRetriesPIN int        `json:"numberOfRetriesPIN" bson:"numberOfRetriesPIN"`
	LastLoginAt        *time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
}

type SellerRegistrationRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
}
type GetOTP struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
}

type SellerRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SellerVerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CustomerLoginRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
	PIN         int    `json:"pin" binding:"required"`
}

type CustomerLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CustomerRegistrationRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
}

type CustomerRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type AddBasicData struct {
	SellerID string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Address  string `json:"address" bson:"address"`
	PIN      int    `json:"pin" bson:"pin,min=6"`
}

type ResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type CustomerBasicSetupRequest struct {
	UserId  string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
	PIN     int    `json:"pin" bson:"pin,min=6"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Admin struct {
	AdminID                string     `json:"id" bson:"_id,omitempty"`
	Email                  string     `json:"email" bson:"email"`
	Password               string     `json:"-" bson:"password"` // "-" means this field won't be included in JSON responses
	Name                   string     `json:"name" bson:"name"`
	IsFirstLogin           bool       `json:"isFirstLogin" bson:"isFirstLogin"`
	PhoneNumber            string     `json:"phoneNumber" bson:"phoneNumber"`
	Address                string     `json:"address" bson:"address"`
	EmergencyContactNumber string     `json:"emergencyContactNumber" bson:"emergencyContactNumber"`
	CreatedAt              time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt" bson:"updatedAt"`
	LastLoginAt            *time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
}

type AdminRegistrationRequest struct {
	Email                  string `json:"email" binding:"required,email"`
	Password               string `json:"password" binding:"required,min=5"`
	Name                   string `json:"name" binding:"required,min=2"`
	PhoneNumber            string `json:"phoneNumber" binding:"required,min=10"`
	Address                string `json:"address" binding:"required,min=10"`
	EmergencyContactNumber string `json:"emergencyContactNumber" binding:"required,min=10"`
}

type AdminRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
