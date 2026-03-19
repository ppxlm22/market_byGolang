package repository

type Repository interface {
	GetUserByUsername(username string) (int, string, error)
}
