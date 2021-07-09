package stripe_client

import (
	"context"
	"strconv"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/accountlink"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
)

type Client struct {
	Key        string
	Logger     core_logging.ILog
	RefreshUrl string
	ReturnUrl  string
}

type Response struct {
	Url      string
	StripeId uint32
}

type ClientParams struct {
	Key        string
	RefreshUrl string
	ReturnUrl  string
}

func NewStripeClient(logger core_logging.ILog, params ClientParams) *Client {
	return &Client{
		Key:        params.Key,
		Logger:     logger,
		RefreshUrl: params.RefreshUrl,
		ReturnUrl:  params.ReturnUrl,
	}
}

func (s *Client) DeleteAccount(ctx context.Context, acctId uint32) error {
	stripe.Key = s.Key
	id := strconv.Itoa(int(acctId))
	_, err := account.Del(id, nil)
	if err != nil {
		s.Logger.For(ctx).Error(err, err.Error())
		return err
	}

	return nil
}

func (s *Client) CreateAccount(ctx context.Context, request *merchant_service_proto_v1.CreateAccountRequest) (*Response, error) {
	stripe.Key = s.Key
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

	acct, err := account.New(params)
	if err != nil {
		s.Logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	stripeAccountId, err := strconv.Atoi(acct.ID)
	if err != nil {
		s.Logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	// call the account link api
	acctLinkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(acct.ID),
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
		StripeId: uint32(stripeAccountId),
	}, nil
}
