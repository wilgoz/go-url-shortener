package shortener

// RedirectRepository defines a common interface for repositories
type RedirectRepository interface {
	Find(shortened string) (*Redirect, error)
	Store(model *Redirect) error
}
