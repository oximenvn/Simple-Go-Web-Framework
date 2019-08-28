package main

type Config struct {
	Database struct {
		Host   string
		User   string
		Pass   string
		Name   string
		Driver string
	}
	App struct {
		Env string
	}
}
