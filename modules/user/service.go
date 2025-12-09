package user

// Service defines user-related business operations.
type Service interface {
	Create(u User) (*User, error)
	GetByID(id int64) (*User, error)
	List() ([]User, error)
	GetByUsername(username string) (*User, error)
}
type service struct {
	repo Repository
}

func (s *service) Create(u User) (*User, error) {
	return s.repo.Create(u)
}

func (s *service) List() ([]User, error) {
	return s.repo.List()
}
func (s *service) GetByID(id int64) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByUsername(username string) (*User, error) {
	return s.repo.GetByUsername(username)
}
