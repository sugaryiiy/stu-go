package user

// Service defines user-related business operations.
type service struct {
	repo repository
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
func (s *service) DeleteByUserName(username string) error {
	s.repo.DeleteByUserName(username)
	return nil
}
