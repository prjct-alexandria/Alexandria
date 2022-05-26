package entities

type Version struct {
	ArticleID int64
	Id        int64
	Title     string
	Owners    []string
	Status    string
}

const (
	VersionDraft    string = "draft"
	VersionPending  string = "pending"
	VersionReview   string = "review"
	VersionAccepted string = "accepted"
)
