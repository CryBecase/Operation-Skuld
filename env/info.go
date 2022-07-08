package env

type Info struct {
	Server  string
	Version string
	BuildAt int64
	Envv    Env
}

func NewInfo(server, version string, buildAt int64, envv Env) Info {
	return Info{
		Server:  server,
		Version: version,
		BuildAt: buildAt,
		Envv:    envv,
	}
}
