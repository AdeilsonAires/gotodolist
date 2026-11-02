package domain

import "errors"

var (
	ErrNotFound = errors.New("registro não encontrado")
	ErrEmailAlreadyExists = errors.New("email já cadastrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrWrongPassword = errors.New("senha atula incorreta")
	ErrInvalidToken = errors.New("token inválido ou expirado")
	ErrUnauthorized = errors.New("não autenticado")
	ErrForbidden = errors.New("acesso negado")
)