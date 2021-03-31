package interactor

type HealthCheck struct {
	store Store
}

func NewHealthCheck(store Store) *HealthCheck {
	return &HealthCheck{store: store}
}

func (hc *HealthCheck) CheckStoreConnection() error {
	return hc.store.CheckStoreConnection()
}
