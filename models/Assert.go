package models


type IEvent interface {
	Create(model interface{}) (int64, error)
}
type IModel interface {

}
type User struct {
	Id int64
	Name string
	IModel
}
type UserHandler struct {
}

func (this *UserHandler) Create(user *User) (int64, error){
	//dummy := (*User)(unsafe.Pointer(model))
	//user,ok := model.(*User)
	//if !ok{
	//	//return
	//}

	//u:= (&(model.(*User))).(User)
	return  user.Id,nil
}