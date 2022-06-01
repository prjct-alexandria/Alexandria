package entities

type Version struct {
	ArticleID int64
	Id        int64
	Title     string
	Owners    []string
	Status    string
}

// Possible statuses
// Depending on the context,
// Accepted can be interpreted as "ready to merge", or "accepted"
const (
	VersionDraft    string = "draft"
	VersionPending  string = "pending"
	VersionReview   string = "review"
	VersionAccepted string = "accepted"
)
