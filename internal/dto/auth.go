// Alcatraz - Secure open source Password Manager and Storage System
// Copyright (C) 2026 Gian Carlo Ruiz Patiño
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package dto contains the Data Transfer Objects (request/response contracts) for the Alcatraz API.
package dto

type RegisterDTO struct {
	Email              string `json:"email" validate:"required,email"`
	Password           string `json:"password" validate:"required,min=8"`
	ProtectedMasterKey         string `json:"protected_master_key" validate:"required"`
	MasterKeyIV                string `json:"master_key_iv" validate:"required"`
	MasterKeySalt              string `json:"master_key_salt" validate:"required"`
	RecoveryKey                string `json:"recovery_key" validate:"required,min=8"`
	RecoveryProtectedMasterKey string `json:"recovery_protected_master_key" validate:"required"`
	RecoveryKeyIV              string `json:"recovery_key_iv" validate:"required"`
	RecoveryKeySalt            string `json:"recovery_key_salt" validate:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	Token              string `json:"token,omitempty"`
	Require2FA         bool   `json:"require_2fa"`
	ProtectedMasterKey string `json:"protected_master_key,omitempty"`
	MasterKeyIV        string `json:"master_key_iv,omitempty"`
	MasterKeySalt      string `json:"master_key_salt,omitempty"`
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

type ChangeMasterPasswordDTO struct {
	OldPassword        string `json:"old_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	ProtectedMasterKey string `json:"protected_master_key" validate:"required"`
	MasterKeyIV        string `json:"master_key_iv" validate:"required"`
	MasterKeySalt      string `json:"master_key_salt" validate:"required"`
}

type FetchRecoveryDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDTO struct {
	Email              string `json:"email" validate:"required,email"`
	RecoveryKey        string `json:"recovery_key" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	ProtectedMasterKey string `json:"protected_master_key" validate:"required"`
	MasterKeyIV        string `json:"master_key_iv" validate:"required"`
	MasterKeySalt      string `json:"master_key_salt" validate:"required"`
}
