package messages

type Message string

var (
	Unauthorized                 Message = "Unauthorized"
	WrongLoginOrPassword         Message = "Wrong login or password"
	MissingOrInvalidToken        Message = "Missing or invalid token"
	RefreshTokenNotFound         Message = "Refresh token not found"
	InternalServerError          Message = "Internal server error"
	InvalidTefreshToken          Message = "Invalid refresh token"
	CouldNotGenerateAccessToken  Message = "Could not generate access token"
	CouldNotGenerateRefreshToken Message = "Could not generate refresh token"
	InvalidPassword              Message = "Invalid password"
	InvalidRole                  Message = "Invalid role"
	InvalidPayload               Message = "Invalid payload"
	NotFound                     Message = "Not found"
)
