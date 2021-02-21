package base

type Sessioner interface {
	Send(msg Msg)
	Run()
	SaveLoginState(u User)
	GetUser() *base.User
	GetServer() base.Serverer
	Close()
}
