package stripe_client

import (
	"context"

	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/accountlink"
	"github.com/stripe/stripe-go/v72"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
)

type StripeOperations interface {
	DeleteAccount(ctx context.Context, accountId string) error
	CreateAccount(ctx context.Context, request *merchant_service_proto_v1.CreateAccountRequest) (*Response, error)
	GetAccount(ctx context.Context, accountId string) (*stripe.Account, error)
	GetAccountLink(ctx context.Context, accountId string) (*Response, error)
}

type Client struct {
	Key        string
	Logger     core_logging.ILog
	RefreshUrl string
	ReturnUrl  string
}

var _ StripeOperations = (*Client)(nil)

type Response struct {
	Url      string
	StripeId string
}

type ClientParams struct {
	Key        string
	RefreshUrl string
	ReturnUrl  string
}

func NewStripeClient(logger core_logging.ILog, params *ClientParams) (*Client, error) {
	if params == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	stripe.Key = params.Key

	return &Client{
		Key:        params.Key,
		Logger:     logger,
		RefreshUrl: params.RefreshUrl,
		ReturnUrl:  params.ReturnUrl,
	}, nil
}

func (s *Client) DeleteAccount(ctx context.Context, id string) error {
	if id == constants.EMPTY {
		return service_errors.ErrInvalidInputArguments
	}

	_, err := account.Del(id, nil)
	if err != nil {
		s.Logger.For(ctx).Error(err, err.Error())
		return err
	}

	return nil
}

func (s *Client) CreateAccount(ctx context.Context, request *merchant_service_proto_v1.CreateAccountRequest) (*Response, error) {
	params := &stripe.AccountParams{
		Type: stripe.String(string(stripe.AccountTypeStandard)),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		Country: stripe.String("US"),
		Email:   stripe.String(request.Account.BusinessEmail),
	}

	// TODO: implement this as a retryable operation for all 429+ error codes
	acct, err := account.New(params)
	if err != nil {
		s.Logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	return s.GetAccountLink(ctx, acct.ID)
}

func (s *Client) GetAccountLink(ctx context.Context, accountId string) (*Response, error) {
	// call the account link api
	acctLinkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(accountId),
		RefreshURL: stripe.String(s.RefreshUrl),
		ReturnURL:  stripe.String(s.ReturnUrl),
		Type:       stripe.String(constants.STRIPE_ACCOUNT_ONBOARDING),
	}

	acc, err := accountlink.New(acctLinkParams)
	if err != nil {
		return nil, err
	}

	return &Response{
		Url:      acc.URL,
		StripeId: accountId,
	}, nil
}

func (s *Client) GetAccount(ctx context.Context, accountId string) (*stripe.Account, error) {
	if accountId == constants.EMPTY {
		return nil, service_errors.ErrInvalidInputArguments
	}

	return account.GetByID(accountId, nil)
}

func (s *Client) HandleStripeError(err error) (bool, error) {
	if err != nil {
		// Try to safely cast a generic error to a stripe.Error so that we can get at
		// some additional Stripe-specific information about what went wrong.
		if stripeErr, ok := err.(*stripe.Error); ok {
			shouldRetry := s.ClientRetry(stripeErr.Type)
			if stripeErr.HTTPStatusCode >= 429 {
				s.Logger.Error(stripeErr.Err, stripeErr.Error())
				return shouldRetry, stripeErr
			}
		} else {
			s.Logger.Error(err, err.Error())
		}
	}

	return false, nil
}

func (s *Client) ClientRetry(errorType stripe.ErrorType) bool {
	switch errorType {
	case stripe.ErrorTypeAPI:
		return true
	case stripe.ErrorTypeRateLimit:
		return true
	default:
		return false
	}
}
