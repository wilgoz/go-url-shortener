package shortener

type RedirectRepository interface {
	Find(shortened string) (*Redirect, error)
	Store(model *Redirect) error
}
