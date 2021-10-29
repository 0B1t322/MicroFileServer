package auth

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/MicroFileServer/pkg/contextvalue/rolecontext"
	"github.com/MicroFileServer/pkg/contextvalue/subcontext"
	"github.com/MicroFileServer/pkg/utils"

	"github.com/MicahParks/keyfunc"

	ctxtoken "github.com/MicroFileServer/pkg/contextvalue/token"
	"github.com/MicroFileServer/pkg/statuscode"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v4"

	"github.com/MicroFileServer/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Auth struct {
	User		string
	Admin		string

	// Claim where roles
	Claim		string

	jwks 		*keyfunc.JWKs

	auth		endpoint.Middleware
	admin		endpoint.Middleware
}

type Config struct {
	*config.AuthConfig
	Testmode	bool
}

func NewAuth(
	cfg		*Config,
) *Auth {
	a := &Auth{
		User: 	cfg.UserRole,
		Admin:	cfg.AdminRole,
		Claim: cfg.Audience,
	}

	refreshTime := 24*time.Hour

	jwks, err := keyfunc.Get(
		cfg.KeyURL, 
		keyfunc.Options{
			RefreshInterval: &refreshTime,
			RefreshErrorHandler: func(err error) {
				log.WithFields(
					log.Fields{
						"func": "refreshTokernErrorHandle",
						"err": err,
					},
				).Error()
			},
		},
	)
	if err != nil {
		log.WithFields(
			log.Fields{
				"func": "authMiddleware",
				"error": err,
			},
		).Panic("Failed to create jwks")
	}
	a.jwks = jwks

	if cfg.Testmode {
		a.BuildTestAuthMiddleware()
	} else {
		a.BuildAuthMiddleware()
	}

	a.BuildAdminMiddleware()

	return a
}

func (a *Auth) AuthMiddleware() endpoint.Middleware {
	return a.auth
}

func (a *Auth) IsAdmin() endpoint.Middleware {
	return a.admin
}

func (a *Auth) BuildTestAuthMiddleware() {
	a.auth = func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			log.Debug("Test auth")
			_t, err := ctxtoken.GetTokenFromContext(ctx)

			if err == ctxtoken.ErrTokenNotFound {
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Token not found"),
					http.StatusUnauthorized,
				)
			}

			jwtToken := strings.ReplaceAll(_t, "Bearer ","")
			var claims jwt.MapClaims

			_, err = jwt.ParseWithClaims(
				jwtToken,
				&claims,
				func(t *jwt.Token) (interface{}, error) {
					log.Debug("key func")
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					log.Debug("okay")
					return []byte("test"), nil
				},
			)
			if validErr, ok := err.(*jwt.ValidationError); ok {
				switch validErr.Errors {
				case jwt.ValidationErrorExpired:
					// pass
				}
			} else if err != nil {
				log.Error("Err:", err)
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Failed to parse token"),
					http.StatusUnauthorized,
				)
			}

			// if !token.Valid {
			// 	return nil, statuscode.WrapStatusError(
			// 		fmt.Errorf("Token is not valid"),
			// 		http.StatusUnauthorized,
			// 	)
			// }
			log.Debug(claims)
			role, err := a.RoleGetter(claims)
			if err != nil {
				log.WithFields(log.Fields{
					"package" : "middleware/auth",
					"func": "authMiddleware",
					"error" : err,
				}).Error("Failed to get role")

				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Faield to get role"),
					http.StatusUnauthorized,
				)
			}

			ctx = rolecontext.New(
				ctx,
				role,
			)
			log.Debug("role is:", role)

			sub := fmt.Sprint(claims["sub"])

			log.Debug("sub is:", sub)
			ctx = subcontext.New(
				ctx,
				sub,
			)

			return next(ctx, request)
		}
	}
}

func (a *Auth) RoleGetter(claims map[string]interface{}) (string, error) {
	var roles []string
	{
		switch claim := claims[a.Claim].(type) {
		case []interface{}:
			for _, role := range claim {
				roles = append(roles, fmt.Sprint(role))
			}
		case interface{}:
			roles = append(roles, fmt.Sprint(claim))
		}
	}
	if len(roles) == 0 {
		return "", fmt.Errorf("Don't find role")
	}

	if len(roles) == 1 {
		role := roles[0]
		if a.isAdmin(role) {
			return a.Admin, nil
		} else if a.isUser(role) {
			return a.User, nil
		} else {
			return "", fmt.Errorf("Don't find role")
		}
	} else {
		sortedRoles := sort.StringSlice(roles)
		sortedRoles.Sort()
		if role, find := utils.FindString(
			sortedRoles,
			a.Admin,
		); find {
			return role, nil
		}

		if role, find := utils.FindString(
			sortedRoles,
			a.User,
		); find {
			return role, nil
		}
	}

	return "", fmt.Errorf("don't find role")
}

func (a *Auth) isAdmin(role string) bool {
	return role == a.Admin
}

func (a *Auth) isUser(role string) bool {
	return role == a.User
}

func (a *Auth) BuildAuthMiddleware() {
	a.auth = func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			_t, err := ctxtoken.GetTokenFromContext(ctx)
			if err == ctxtoken.ErrTokenNotFound {
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Token not found"),
					http.StatusUnauthorized,
				)
			}

			jwtToken := strings.ReplaceAll(_t, "Bearer ","")
			var claims jwt.MapClaims

			token, err := jwt.ParseWithClaims(jwtToken, &claims, a.jwks.Keyfunc)
			if validErr, ok := err.(*jwt.ValidationError); ok {
				switch validErr.Errors {
				case jwt.ValidationErrorExpired:
					return nil, statuscode.WrapStatusError(
						fmt.Errorf("Token expired"),
						http.StatusUnauthorized,
					)
				}
			} else if err != nil {
				log.Error("Err:", err)
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Failed to parse token"),
					http.StatusUnauthorized,
				)
			}

			if !token.Valid {
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Token is not valid"),
					http.StatusUnauthorized,
				)
			}
			
			role, err := a.RoleGetter(claims)
			if err != nil {
				log.WithFields(log.Fields{
					"package" : "middleware/auth",
					"func": "authMiddleware",
					"error" : err,
				}).Error("Failed to get role")

				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Faield to get role"),
					http.StatusUnauthorized,
				)
			}

			ctx = rolecontext.New(
				ctx,
				role,
			)

			sub := fmt.Sprint(claims["sub"])

			ctx = subcontext.New(
				ctx,
				sub,
			)

			return next(ctx, request)
		}
	}
}

func (a *Auth) BuildAdminMiddleware() {
	a.admin = func(next endpoint.Endpoint) endpoint.Endpoint {
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

			if role != a.Admin {
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("You are not admin"),
					http.StatusForbidden,
				)
			}
			return next(ctx, request)
		}
	}
}