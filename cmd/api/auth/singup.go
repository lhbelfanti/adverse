package auth

import (
	"context"
	"golang.org/x/crypto/bcrypt"

	"ahbcc/cmd/api/users"
	"ahbcc/internal/log"
)

// SignUp registers a new user in the system
type SignUp func(ctx context.Context, user users.UserDTO) error

// MakeSignUp creates a new SignUp
func MakeSignUp(userExists users.UserExists, insertUser users.Insert) SignUp {
	return func(ctx context.Context, user users.UserDTO) error {
		exists, err := userExists(ctx, user.Username)
		if err != nil {
			log.Error(ctx, err.Error())
			return FailedToRetrieveIfTheUserExists
		}

		if exists {
			log.Error(ctx, FailedToSignUpBecauseTheUserAlreadyExists.Error())
			return FailedToSignUpBecauseTheUserAlreadyExists
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(ctx, err.Error())
			return FailedToGenerateHashFromPassword
		}

		user.Password = string(hash)

		err = insertUser(ctx, user)
		if err != nil {
			log.Error(ctx, err.Error())
			return FailedToInsertUserIntoDatabase
		}

		return nil
	}
}
