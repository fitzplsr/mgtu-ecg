package middleware

type KeyType struct{}
type keyTypePrivate struct{}

var (
	UserIDKey   KeyType
	UserRoleKey keyTypePrivate
)
