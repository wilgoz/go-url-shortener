package shortener

// RedirectService defines a common interface for the business layer services
type RedirectService interface {
	Find(shortened string) (*Redirect, error)
	Store(model *Redirect) error
}
