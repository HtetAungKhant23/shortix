package shortener

type Service struct {
	store   Store
	codeLen int
}

func New(store Store) *Service {
	return &Service{
		store:   store,
		codeLen: 6,
	}
}
