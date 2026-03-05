package session

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session)
}

type IAfterQuery interface {
	AfterQuery(s *Session)
}

type IBeforeUpdate interface {
	BeforeUpdate(s *Session)
}

type IAfterUpdate interface {
	AfterUpdate(s *Session)
}

type IBeforeDelete interface {
	BeforeDelete(s *Session)
}

type IAfterDelete interface {
	AfterDelete(s *Session)
}

type IBeforeInsert interface {
	BeforeInsert(s *Session)
}

type IAfterInsert interface {
	AfterInsert(s *Session)
}

func CallHookMethod[T any](value any, fn func(method T)) {
	if can, ok := value.(T); ok {
		fn(can)
	}
}
