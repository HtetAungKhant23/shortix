package shortener

type Store interface {
	Save(url *URL) error
}
