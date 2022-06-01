package entities

type Thread struct {
	Id        int64
	ArticleId int64
}

type ReviewThread struct {
	Id       int64
	ReviewId int64
	ThreadId int64
}

type CommitThread struct {
	Id       int64
	CommitId int64
	ThreadId int64
}

type RequestThread struct {
	Id        int64
	RequestId int64
	ThreadId  int64
}
