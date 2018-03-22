package models

import (
	"github.com/dotnet123/fasthttptest/util"

	"github.com/dotnet123/fasthttptest/ext"
	"errors"
)

func init() {
	util.InitHandler(&UserHandler{})
}

type User struct {
	Id   int64
	Name string
}
type Query struct {
	Take  int32
	Skip  int32
	Count int32
}
type UserQuery struct {
	Query
	Name   ext.NullableString
	Qty    int32
	Return []User
}

type UserHandler struct {
}
type UserDal struct {
}

func (UserHandler) Create(user *User, dal UserDal) (int64, error) {

	i, _ := dal.Create(user)
	return user.Id + 3 + i, nil
}
func (UserHandler) Select(query *UserQuery) (userLst []User,count int64,err error) {
	a := User{int64(1), "test"}
	b := make([]User, 0)
	b = append(b, a)
	query.Return = b
	return b,0,nil
}
func (*UserDal) Create(user *User) (int64, error) {

	return user.Id + 2, errors.New("55666")
}


