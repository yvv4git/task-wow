package configs

type Server struct {
	Hostname   string
	Port       int
	DebugLevel bool
	Difficulty int
}

func NewServerConfig(hostname string, port int, enableDebug bool, difficulty int) *Server {
	return &Server{
		Hostname:   hostname,
		Port:       port,
		DebugLevel: enableDebug,
		Difficulty: difficulty,
	}
}
