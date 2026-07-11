package authentication

type State struct {
	Registered bool
	LoggedIn   bool

	UserID string
	Token  string
}
