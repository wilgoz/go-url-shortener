package shortener

type RedirectService interface {
	Find(shortened string) (*Redirect, error)
	Store(model *Redirect) error
}
