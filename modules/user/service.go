package user

// Service defines user-related business operations.
type Service interface {
	Create(u User) (User, error)
	GetByID(id int64) (User, error)
	List() ([]User, error)
}
