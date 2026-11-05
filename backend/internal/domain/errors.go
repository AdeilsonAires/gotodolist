package domain

import "errors"

var (
	ErrNotFound           = errors.New("registro não encontrado")
	ErrEmailAlreadyExists = errors.New("email já cadastrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrWrongPassword      = errors.New("senha atual incorreta")
	ErrInvalidToken       = errors.New("token inválido ou expirado")
	ErrUnauthorized       = errors.New("não autenticado")
	ErrForbidden          = errors.New("acesso negado")
	ErrItemNameExists     = errors.New("já existe um item com esse nome na lista")
	ErrListNameExists     = errors.New("já existe uma lista com esse nome")
	ErrUserNameExists     = errors.New("já existe uma usuário com esse nome")
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrInvalidName        = errors.New("nome inválido")
	ErrInvalidIcon        = errors.New("ícone inválido")
	ErrInvalidColor       = errors.New("cor inválida")
	ErrInvalidEmail       = errors.New("email inválido")
	ErrInvalidPassword    = errors.New("senha inválida")
)
