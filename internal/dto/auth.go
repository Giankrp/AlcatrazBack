package dto

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	Token      string `json:"token,omitempty"`
	Require2FA bool   `json:"require_2fa"`
}

type Setup2FAResponseDTO struct {
	Secret string `json:"secret"`
	QRURI  string `json:"qr_uri"`
}

type Enable2FADTO struct {
	Code   string `json:"code" validate:"required,len=6"`
	Secret string `json:"secret" validate:"required"`
}

type Verify2FADTO struct {
	Code string `json:"code" validate:"required"`
}
