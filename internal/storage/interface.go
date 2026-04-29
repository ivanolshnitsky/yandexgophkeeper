package storage

// Storage — интерфейс хранилища данных.
type Storage interface {
	Save(user string, d Data) error
	List(user string) ([]Data, error)
	GetByID(user, id string) (Data, bool, error)
	Update(user string, d Data) (bool, error)
	Delete(user, id string) (bool, error)
	Close() error
}
