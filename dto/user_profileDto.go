package dto

type UpdateUserProfileDTO struct {
	Name      string `json:"name" validate:"omitempty,min=1,max=50"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	Language  string `json:"language" validate:"omitempty,oneof=es en fr de pt"`
}
