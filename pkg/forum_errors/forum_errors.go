package forum_errors

const (
	UserAlreadyExists = iota
	CantFindUser
	ConflictingData
	Unknown
	ForumAlreadyExists
	CantFindForum
)

type UserError struct {
	Reason error
	Code   int32
}

func (e *UserError) Error() string {
	return e.Reason.Error()
}

type ForumError struct {
	Reason error
	Code   int32
}

func (e *ForumError) Error() string {
	return e.Reason.Error()
}
