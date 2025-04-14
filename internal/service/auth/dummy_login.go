package auth

func (s *authService) CreateDummyLogin(role string) (string, error) {
	id := 1 //Фиктивный ID
	token, err := generateJWT(id, role)

	if err != nil {
		return "", err
	}

	return token, nil
}
