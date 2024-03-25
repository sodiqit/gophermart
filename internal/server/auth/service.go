package auth

type AuthService interface {
	Register(username string, password string) (string, error)
	LogIn(username string, password string) (string, error)
}

type SimpleAuthService struct {
	tokenService TokenService
}

func (s *SimpleAuthService) Register(username string, password string) (string, error) {
	return "", nil
}

func (s *SimpleAuthService) LogIn(username string, password string) (string, error) {
	return "", nil
}

func NewSimpleAuthService(tokenService TokenService) *SimpleAuthService {
	return &SimpleAuthService{
		tokenService: tokenService,
	}
}
