package core

type Config struct {
	Database struct {
		Host   string
		Post   int
		User   string
		Pass   string
		Name   string
		Driver string
	}
	App struct {
		Env string
	}
}
