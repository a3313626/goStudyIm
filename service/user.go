package service

import (
	"chat/model"
	"chat/serializer"
)

type UserRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	count := 0
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)

	if count != 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "该用户名已被注册,请更换",
		}
	}

	user = model.User{}

	user.UserNmae = service.UserName
	err := user.SetPassword(service.Password)
	if err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    "密码生成失败",
		}
	}

	err = model.DB.Create(&user).Error
	if err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    "注册失败,请重试",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "注册成功",
	}

}
