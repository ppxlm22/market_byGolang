package repository

type Repository interface {
	GetUserByUsername(username string) (int, string, string, error)
}
