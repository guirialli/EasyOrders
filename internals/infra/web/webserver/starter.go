package webserver

type Starter struct {
	Server WebServer
}

func NewWebServerStarter(server WebServer) *Starter {
	return &Starter{
		Server: server,
	}
}
