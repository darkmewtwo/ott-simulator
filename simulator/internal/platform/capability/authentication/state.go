package authentication

type AuthenticationState struct {
	Registered bool
	LoggedIn   bool

	UserID string
	Token  string
}
