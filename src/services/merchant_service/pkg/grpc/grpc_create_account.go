package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/stripe_client"
)

func (s *Server) CreateAccount(ctx context.Context, request *merchant_service_proto_v1.CreateAccountRequest) (*merchant_service_proto_v1.CreateAccountResponse, error) {
	sagaSteps := make([]*saga.Step, 0)
	var stripeResponseObj = make(chan *stripe_client.Response, 1)
	var authnAcctId = make(chan uint32, 1)
	var merchantAcctId = make(chan uint64, 1)
	var jwtToken = make(chan string, 1)

	// 1. perform request validations
	operationType := constants.CREATE_MERCHANT_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if request == nil || request.Account == nil {
		s.logger.For(ctx).Error(errors.ErrInvalidInputArguments, errors.ErrInvalidInputArguments.Error())
		return nil, errors.ErrInvalidInputArguments
	}

	merchantAcct := request.Account
	/*
		we perform the create account operation as a distributed tx hence it is imperative
		we configure a set of sagas ... each with their own compensating transactions to roll back the state of the system in the
		event child sagas fail
	*/
	// step 1. create an account through the authentication handler service
	createAuthHandlerSvcAcct := s.sagaCreateAccountThroughAuthenticationHandlerService(merchantAcct, authnAcctId)
	// step 2. create an account through stripe
	createStripeAcct := s.sagaCreateAccountThroughStripe(request, merchantAcct, stripeResponseObj)
	// step 3. save the created account in the merchant service's database
	saveAcctToDb := s.sagaSaveCreatedAccountThroughDB(merchantAcct, merchantAcctId)
	// step 4. obtain the jwt token from the authentication handler service which will be used to authenticate the account
	getJwtTokenFromAuthHandlerSvc := s.sagaGetJwtTokenFromAuthHandlerSvc(merchantAcct, jwtToken)

	sagaSteps = append(sagaSteps,
		createAuthHandlerSvcAcct,
		createStripeAcct,
		saveAcctToDb,
		getJwtTokenFromAuthHandlerSvc)
	if err := s.DbConn.Saga.RunSaga(ctx, "create_account", sagaSteps...); err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	stripeObj := <-stripeResponseObj

	return &merchant_service_proto_v1.CreateAccountResponse{
		AccountId: <-merchantAcctId,
		JwtToken:  <-jwtToken,
		StripeUrl: stripeObj.Url,
	}, nil
}
