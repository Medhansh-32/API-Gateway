package service

import (
	"errors"
	"log"

	"github.com/medhansh-32/api-gateway/internal/models/requests"
	"github.com/medhansh-32/api-gateway/internal/models/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct{
	userService UserService
	jwtService JWTService
}

func NewAuthService(userService UserService, jwtService JWTService) (*AuthenticationService){
	return &AuthenticationService{userService: userService,jwtService: jwtService}
}

func (authenticationService AuthenticationService) Login(loginRequest requests.LoginRequest) (*response.LoginResponse, error){
	
	userName:=loginRequest.GetUserName()
	password:=loginRequest.GetPassword()
	
	if userName == ""{
		return nil, errors.New("User Name Not Found")
	}
	if password == ""{
		return nil, errors.New("Password Not Found")	
	}
	
	user,err:= authenticationService.userService.GetUserByUserName(userName)
	
	if err!=nil{
		return nil, errors.New("User Not Found with username : "+userName)
	}
	

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(password))
	if err!=nil{
		log.Print(password)
		log.Print(user.PasswordHash)
		return nil, errors.New("Wrong Password")
	}

   token,err:= authenticationService.jwtService.GenerateToken(user)

   if err!=nil{
	return nil,err
   }

   return &response.LoginResponse{
		UserId : user.ID,
		UserName : user.Username,
		Role : user.Role,
		Email : user.Email,
		Token : token,
   },nil

}

func (authenticationService AuthenticationService) ValidateToken(token string) (*Claims,error){
	return authenticationService.jwtService.ValidateToken(token) 
}