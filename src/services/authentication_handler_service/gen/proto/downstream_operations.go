package proto

import (
	"strconv"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
)

type DownStreamOperation func() (interface{}, error)

// ServiceDependentOperations defines a set of downstream service operations this service relies on
type ServiceDependentOperations interface {
	// CallAuthenticationService performs a downstream call to the authentication service
	CallAuthenticationService(client *core_auth_sdk.AuthService) DownStreamOperation
}

func (req *AuthenticateAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		return client.LoginAccount(req.Email, req.Password)
	}
}

func (req *CreateAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		return client.ImportAccount(req.Email, req.Password, false)
	}
}

func (req *DeleteAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		id := strconv.Itoa(int(req.GetId()))
		if err := client.ArchiveAccount(id); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func (req *GetAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		id := strconv.Itoa(int(req.GetId()))
		account, err := client.GetAccount(id)
		if err != nil {
			return nil, err
		}
		return account, nil
	}
}

func (req *LockAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		id := strconv.Itoa(int(req.GetId()))
		if err := client.LockAccount(id); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func (req *UnLockAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		id := strconv.Itoa(int(req.GetId()))
		if err := client.UnlockAccount(id); err != nil {
			return nil, err
		}
		return nil, nil
	}
}


func (req *LogoutAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		return nil, client.LogOutAccount()
	}
}

func (req *UpdateAccountRequest) CallAuthenticationService(client core_auth_sdk.AuthService) DownStreamOperation {
	return func() (interface{}, error) {
		id := strconv.Itoa(int(req.Id))
		return nil, client.Update(id, req.Email)
	}
}
