package grpc

import (
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
)

type AuthenticationServiceMock struct {
	GetAccountFunc     func(id string) (*core_auth_sdk.Account, error)
	UpdateAccountFunc  func(id, username string) error
	LockAccountFunc    func(id string) error
	UnlockAccountFunc  func(id string) error
	ArchiveAccountFunc func(id string) error
	ImportAccountFunc  func(username, password string, locked bool) (int, error)
	ExpirePasswordFunc func(id string) error
	LoginAccountFunc   func(username, password string) (string, error)
	SignupAccountFunc  func(username, password string) (string, error)
	LogOutAccountFunc  func() error
}

var ThirdPartyMockService = AuthenticationServiceMock{}

func (m *AuthenticationServiceMock) GetAccount(id string) (*core_auth_sdk.Account, error) {
	return m.GetAccountFunc(id)
}

func (m *AuthenticationServiceMock) Update(id, username string) error {
	return m.UpdateAccountFunc(id, username)
}

func (m *AuthenticationServiceMock) LockAccount(id string) error {
	return m.LockAccountFunc(id)
}

func (m *AuthenticationServiceMock) UnlockAccount(id string) error {
	return m.UnlockAccountFunc(id)
}

func (m *AuthenticationServiceMock) ArchiveAccount(id string) error {
	return m.ArchiveAccountFunc(id)
}

func (m *AuthenticationServiceMock) ImportAccount(username, password string, locked bool) (int, error) {
	return m.ImportAccountFunc(username, password, locked)
}

func (m *AuthenticationServiceMock) ExpirePassword(id string) error {
	return m.ExpirePasswordFunc(id)
}

func (m *AuthenticationServiceMock) LoginAccount(username, password string) (string, error) {
	return m.LoginAccountFunc(username, password)
}

func (m *AuthenticationServiceMock) SignupAccount(username, password string) (string, error) {
	return m.SignupAccountFunc(username, password)
}

func (m *AuthenticationServiceMock) LogOutAccount() error {
	return m.LogOutAccountFunc()
}
