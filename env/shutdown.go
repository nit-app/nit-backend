package env

func (e *Env) Shutdown() {
	_ = e.DB.Close()
}

func Shutdown() {
	env.Shutdown()
}
