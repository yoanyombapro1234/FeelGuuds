package grpc

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsPasswordOrEmailInValid checks request parameters for validity
func (s *Server) IsPasswordOrEmailInValid(email string, password string, operationType string) (error, bool) {
	err, ok := s.IsValidEmail(email, operationType)
	if !ok {
		return err, true
	}

	err, ok = s.IsValidPassword(password, operationType)
	if !ok {
		return err, true
	}

	return nil, false
}

// IsValidEmail checks if an email is valid.
func (s *Server) IsValidEmail(email string, operationType string) (error, bool) {
	if email == "" {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(operationType).Inc()
		err := service_errors.ErrInvalidInputArguments
		s.logger.Error("invalid input parameters. please specify a valid email", zap.Error(err))
		return err, false
	}

	return nil, true
}

// IsValidPassword checks if a password if valid;
func (s *Server) IsValidPassword(password string, operationType string) (error, bool) {
	if password == "" {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(operationType).Inc()
		err := service_errors.ErrInvalidInputArguments
		s.logger.Error("invalid input parameters. please specify a valid password", zap.Error(err))
		return err, false
	}

	return nil, true
}

// CheckJwtTokenForInValidity checks jwt token and asserts the token is a valid one.
func (s *Server) CheckJwtTokenForInValidity(ctx context.Context, result interface{}, operation string) (error, bool, string) {
	token := fmt.Sprintf("%v", result)
	if token == "" {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(operation)
		err := status.Errorf(codes.Internal, "issue casting to jwt token")
		s.logger.Error("casting error", zap.Error(err))
		return err, true, ""
	}

	return nil, false, token
}

// GetIdFromResponseObject attempts to cast a generic response to a type int and returns the proper value if no errors occurred.
func (s *Server) GetIdFromResponseObject(ctx context.Context, response interface{}, operationType string) (int, error) {
	// TODO: change to to int64 in order to limit overflow from happening if tons of customer accounts are created
	id, ok := response.(int)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(operationType)
		err := status.Errorf(codes.Internal, "failed to convert result to uint32 id value")
		s.logger.Error("casting error", zap.Error(err))
		return 0, err
	}
	return id, nil
}

// IsValidID checks that the ID passed as part of the request parameters is indeed valid.
func (s *Server) IsValidID(Id uint32, operation string) (error, bool) {
	if Id == 0 {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(operation).Inc()
		err := service_errors.ErrInvalidInputArguments
		s.logger.Error("invalid input parameters. please specify a valid user id", zap.Error(err))
		return err, false
	}

	return nil, true
}

// GetAccountFromResponseObject obtains an account object from the response object; this account is obtained via an attempted casting operation
func (s *Server) GetAccountFromResponseObject(ctx context.Context, ok bool, result interface{}, operation string) (*core_auth_sdk.Account, error) {
	account, ok := result.(*core_auth_sdk.Account)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(operation)

		err := service_errors.ErrFailedToCastAccount
		s.logger.Error(err.Error())
		return nil, err
	}
	return account, nil
}
