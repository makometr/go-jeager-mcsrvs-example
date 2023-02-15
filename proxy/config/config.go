package config

type Config struct {
	ServiceName string
	Port        string

	WorkerPort int
	WorkerHost string

	JeagerURL string
}
