package auth

type Authenticator struct{}

func New() *Authenticator {
	return &Authenticator{}
}
