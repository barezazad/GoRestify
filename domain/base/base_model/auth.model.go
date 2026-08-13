package base_model

// Auth model
type Auth struct {
	Username string `json:"username" bind:"required"`
	Password string `json:"password" bind:"required"`
}

// RefreshTokenRequest is used to exchange a refresh token for a new token pair.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" bind:"required"`
}

// ResourceList model
type ResourceList struct {
	Resources      string   `json:"resources"`
	ResourcesArray []string `json:"resources_array"`
}
