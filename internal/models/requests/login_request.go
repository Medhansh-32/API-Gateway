package requests

type LoginRequest struct {
    UserName string `json:"userName"`
    Password string `json:"password"`
}

func (L LoginRequest) GetUserName() (string){
	return L.UserName
}

func (L LoginRequest) GetPassword() (string){
	return L.Password
}