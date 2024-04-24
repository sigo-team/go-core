package services

type Player struct {
	User
	Score int
}

func NewPlayer(us *UserService) Player {
	name := us.NameGenerator.Generate()

	player := Player{
		User:  us.CreateUser(name),
		Score: 0,
	}
	return player
}
