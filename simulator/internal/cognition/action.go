package cognition

type CatalogAction string

const (
	ActionSelectMovie      CatalogAction = "SELECT_MOVIE"
	ActionContinueBrowsing CatalogAction = "CONTINUE_BROWSING"
	ActionLeaveSession     CatalogAction = "LEAVE_SESSION"
)
