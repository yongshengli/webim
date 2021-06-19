package base

type Sessioner interface {
	Send(msg Msg)
	Run()
	SaveLoginState(u User)
	GetUser() *User
	GetServer() Serverer
	Close()
}
