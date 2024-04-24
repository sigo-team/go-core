package services

type Khil struct {
	User
}

func NewKhil(us *UserService) Khil {
	name := us.NameGenerator.Generate()

	khil := Khil{
		User: us.CreateUser(name),
	}
	return khil
}
