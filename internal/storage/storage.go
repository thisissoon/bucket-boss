package storage

type Storer interface {
	Lister
}

type Lister interface {
	List(ext string) ([]string, error)
}

type Deleter interface {
	Delete()
}
