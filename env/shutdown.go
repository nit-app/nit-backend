package env

func (e *Env) Shutdown() {
	if e.DB != nil {
		_ = e.DB.Close()
	}

	if e.Redis != nil {
		_ = e.Redis.Close()
	}
}

func Shutdown() {
	env.Shutdown()
}
