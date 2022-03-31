package softcorp_storage

type Storage interface {
}

type defaultStorage struct {
}

func New() (Storage, error) {
	return &defaultStorage{}, nil
}
