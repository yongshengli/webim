package base

type Serverer interface {
	Run(host, port string, slotContainerLen, slotLen int)
	Unicast(deviceToken string, msg Msg)
	Broadcast(msg Msg) (bool, error)
}
