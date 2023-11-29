package events

type Scanner interface {
	Scan(dest ...any) error
}
