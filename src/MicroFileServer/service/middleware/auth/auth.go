package auth

import (
	"strings"
	"regexp"
	"context"
	"fmt"
	"net/http"

	"github.com/MicroFileServer/pkg/statuscode"

	"github.com/auth0-community/go-auth0"
	"github.com/go-kit/kit/endpoint"
	"gopkg.in/square/go-jose.v2"

	"github.com/MicroFileServer/pkg/config"
	"github.com/MicroFileServer/pkg/contextvalue/rolecontext"
	ctxtoken "github.com/MicroFileServer/pkg/contextvalue/token"
	log "github.com/sirupsen/logrus"
)

func NewGoKitAuth(
	cfg		*config.AuthConfig,
) endpoint.Middleware {
	rolesSet := map[string]struct{}{}

	for _, role := range strings.Split(cfg.RolesConfig.Roles, " ") {
		rolesSet[role] = struct{}{}
	}

	return newGoKitAuth(
		cfg,
		NewRoleGetter(
			"itlab",
			rolesSet,
		),
	)
}

func newGoKitAuth(
	cfg 	*config.AuthConfig,
	f		getRoleFromClaim,
) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			log.Debug("auth middleware")
			client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: cfg.KeyURL}, nil)
			configuration := auth0.NewConfiguration(client, []string{cfg.Audience}, cfg.Issuer, jose.RS256)
			validator := auth0.NewValidator(
				configuration, 
				nil,
			)

			_t, _ := ctxtoken.GetTokenFromContext(ctx)
			r := &http.Request{
				Header: http.Header{
					"Authorization": []string{_t},
				},
			}

			token, err := validator.ValidateRequest(
				r,
			)
			if err != nil {
				log.WithFields(log.Fields{
					"requiredAlgorithm" : "RS256",
					"error" : err,
				}).Debug("Token is not valid!")

				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Token is not valid"),
					http.StatusUnauthorized,
				)
			}
			claims := map[string]interface{}{}
			if err = validator.Claims(r, token, &claims); err != nil {
				log.WithFields(log.Fields{
					"requiredClaims" : "iss, aud, sub, role",
					"error" : err,
				}).Debug("Invalid claims!")
	
				
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Invalid claims"),
					http.StatusUnauthorized,
				)
			}

			role, err := f(claims)
			if err != nil {
				log.WithFields(log.Fields{
					"package" : "middleware/auth",
					"func": "authMiddleware",
					"error" : err,
				}).Debug("Failed to get role")

				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Faield to get role"),
					http.StatusUnauthorized,
				)
			}

			ctx = rolecontext.New(
				ctx,
				role,
			)

			return next(ctx, request)
		}
	}
}

func EndpointAdminMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			log.Debug("admin middleware")
			role, err := rolecontext.GetRoleFromContext(ctx)
			if err != nil {
				log.WithFields(
					log.Fields{
						"package": "middleware/auth",
						"err": err,
					},
				).Panic()
			}

			re := regexp.MustCompile(`\w+.admin`)

			if !re.MatchString(role) {				
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("You are not admin"),
					http.StatusForbidden,
				)
			}

			return next(ctx, request)
		}
	}
}

type getRoleFromClaim func(map[string]interface{}) (string, error)

func NewRoleGetter(
	claimName string,
	rolesSet map[string]struct{},
) getRoleFromClaim {
	return func(claims map[string]interface{}) (string, error) {
		claim, find := claims[claimName]

		if !find {
			return "", fmt.Errorf("Failed to get itlab claim")
		}

		_roles, ok := claim.([]interface{})
		if !ok {
			return "", fmt.Errorf("Failed to cast types")
		}

		roles := sliceOfInterfaceToSliceOfString(_roles)

		for _, role := range roles {
			if _, find := rolesSet[role]; find {
				return role, nil
			}
		}
		
		return "", fmt.Errorf("Failed to get rolse")
	}
}

func sliceOfInterfaceToSliceOfString(slice []interface{}) []string {
	var strs []string

	for _, elem := range slice {
		strs = append(strs, fmt.Sprint(elem))
	}

	return strs
}