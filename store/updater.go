package store

type UpdaterStore interface {
	GetLatestReference() (string, error)
}
