package env

type Scanner interface {
	Scan(dest ...any) error
}
