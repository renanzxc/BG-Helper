package utils

type IAppInterface interface {
	Setup() (err error)
	Run()
	Shutdown()
}
