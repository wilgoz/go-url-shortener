package shortener

import (
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

// RedirectService defines a common interface for the business layer services
type RedirectService interface {
	Find(shortened string) (*Redirect, error)
	Store(model *Redirect) error
}

var (
	// ErrRedirectNotFound signifies no possible redirects
	ErrRedirectNotFound = errors.New("Redirect Not Found")

	// ErrRedirectInvalid signifies an invalid redirect
	ErrRedirectInvalid = errors.New("Redirect Invalid")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

// NewRedirectService initializes and returns the business logic service handler given the repo
func NewRedirectService(repository RedirectRepository) RedirectService {
	return &redirectService{repository}
}

func (r *redirectService) Find(shortened string) (*Redirect, error) {
	return r.redirectRepo.Find(shortened)
}

func (r *redirectService) Store(model *Redirect) error {
	if err := validate.Validate(model); err != nil {
		return errors.Wrap(
			ErrRedirectInvalid, "service.Redirect.Store",
		)
	}
	model.Shortened = shortid.MustGenerate()
	model.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(model)
}
