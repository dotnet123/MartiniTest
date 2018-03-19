package models

func init()  {

}
type IEvent interface {
	Create(model IModel) (int64, error)
}
type IModel interface {

}
type User struct {
	Id int64
	Name string
	IModel
}
type UserHandler struct {
	Next IEvent
}
type UserDal struct {
}
func (this *UserHandler) Create(m IModel) (int64, error){

	user,_ := m.(*User)
	this.Next.Create(m)
	return  user.Id,nil
}
func (this *UserDal) Create(m IModel) (int64, error){

	user,_ := m.(*User)
	return  user.Id,nil
}