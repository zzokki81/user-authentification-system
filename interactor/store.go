package interactor

type Store interface {
	CheckStoreConnection() error
}
