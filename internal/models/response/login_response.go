package response

type LoginResponse struct{
	UserId int64 `json:"userId"`
	UserName string `json:"userName"`
	Role string `json:"role"`
	Email string `json:"email"`
	Token string `json:"token"`
}