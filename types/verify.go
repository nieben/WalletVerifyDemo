package types

type MessageResponse struct {
	Message string `json:"message"`
}

type VerifyRequest struct {
	Address       string `json:"address" validate:"required,len=42"`
	Message       string `json:"message" validate:"required,base64"`
	SignedMessage string `json:"signedMessage" validate:"required,base64"`
}

type VerifyResponse struct {
	Verified bool `json:"verified"`
}
