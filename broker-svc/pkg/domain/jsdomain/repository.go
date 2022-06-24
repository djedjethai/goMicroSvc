package jsdomain

type RepoInterf interface {
	RepoResponse() []JsonResponse
	RepoAddData(dt JsonResponse) (bool, error)
}

type Repository struct {
	db []JsonResponse
}

func NewRepository() *Repository {
	return &Repository{
		db: []JsonResponse{},
	}
}

func (r *Repository) RepoResponse() []JsonResponse {
	return r.db
}

func (r *Repository) RepoAddData(dt JsonResponse) (bool, error) {
	r.db = append(r.db, dt)
	return true, nil
}
