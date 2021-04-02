package interactor

type Store interface {
	CheckStoreConnection() error
	Transaction() (Store, error)
}
