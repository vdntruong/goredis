package ticket

import "context"

type DataProvider interface {
	InsertFilm(c context.Context, f Film) (*Film, error)
}

type DataCacher interface {
}
