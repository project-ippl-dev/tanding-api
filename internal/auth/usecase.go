package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/mail"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase struct {
	repository *db.Queries
	db         *sql.DB
	rdb        *redis.Client
}

func NewUsecase(repository *db.Queries, db *sql.DB, rdb *redis.Client) Usecase {
	return Usecase{repository: repository, db: db, rdb: rdb}
}

func (u Usecase) register(ctx context.Context, req registerRequest, host string, feHost string) (statusCode int, err error) {
	if _, err := u.repository.AccountFetchUserIDByUsername(ctx, db.AccountFetchUserIDByUsernameParams{
		Username: req.Email,
		Type:     "manual",
	}); err == nil {
		return http.StatusForbidden, fmt.Errorf("user with specific email already exist")
	}

	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in hashing password : %s", err.Error())
	}

	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	userID, err := txQuery.UserCreate(ctx, db.UserCreateParams{
		Name:           req.Name,
		Phone:          req.Phone,
		Photo:          "",
		BornAt:         "",
		BornOn:         sql.NullTime{},
		IdentityNumber: "",
		Gender:         "",
		About:          "",
	})
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("rollback fail in create user tx : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in rollback create user tx : %s", err.Error())
	}
	accountID, err := txQuery.AccountCreate(ctx, db.AccountCreateParams{
		Username: req.Email,
		Password: hashedPassword,
		Type:     "manual",
		UserID:   userID,
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("rollback fail in in rollback create account tx : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in rollback create account tx : %s", err.Error())
	}
	//Create Verification Token
	token, err := tools.JWTCreateToken(tools.JWT{
		ID:        userID.String(),
		AccountID: accountID,
		RoleName:  "user",
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("rollback fail in generate token : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in rollback generate token : %s", err.Error())
	}

	if err = u.rdb.Set(ctx, "verify-"+req.Email, token.AccessToken, 24*time.Hour).Err(); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("rollback fail in set redis keyval : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in rollback set redis keyval : %s", err.Error())
	}

	if err = tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	go func(data sendMailRequest) {
		buffer, err := mail.TemplateToBuffer(data.BodyPath, data.BodyParam)
		if err != nil {
			log.Printf("error in converting html file to buffer : %s \n", err.Error())
		}
		mail.SendMail(mail.Request{
			To:      req.Email,
			Subject: "[Tanding!] Please confirm your email address",
			Body:    buffer,
		})
	}(sendMailRequest{
		BodyParam: bodyParam{
			Name: req.Name,
			Path: feHost + "/register/confirm-email?token=" + token.AccessToken,
			Host: &host,
		},
		BodyPath: "public/templates/mail/regist-verification.html",
	})
	return http.StatusCreated, nil
}

func (u Usecase) verify(ctx context.Context, kind string, decoded tools.JWT, token string) (statusCode int, email string, err error) {
	email, err = u.repository.AccountFetchEmailByID(ctx, db.AccountFetchEmailByIDParams{
		ID:   decoded.AccountID,
		Type: db.AccountTypeManual,
	})
	if err != nil {
		return http.StatusNotFound, "", fmt.Errorf("error in fetch account : %s", err.Error())
	}
	retrievedToken, err := u.rdb.Get(ctx, kind+"-"+email).Result()
	if err != nil {
		return http.StatusNotFound, "", fmt.Errorf("error in get redis key, it's possible because token is already used or expired. error : %s", err.Error())
	}
	if token != retrievedToken {
		return http.StatusUnauthorized, "", fmt.Errorf("error in comparing outcoming token with retrieved token")
	}
	if err = u.repository.AccountUpdateStatusByID(ctx, decoded.AccountID); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error in update status account : %s", err.Error())
	}
	if err = u.rdb.Del(ctx, kind+"-"+email).Err(); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error in del key in redis : %s", err.Error())
	}
	return http.StatusOK, email, nil
}

func (u Usecase) login(ctx context.Context, req loginReq) (statusCode int, response loginResponse, err error) {
	account, err := u.repository.AccountFetchOneByEmail(ctx, db.AccountFetchOneByEmailParams{
		Username: req.Username,
		Type:     db.AccountTypeManual,
	})

	if err != nil {
		return http.StatusNotFound, response, fmt.Errorf("error in fetch account : %s", err.Error())
	}
	if err = tools.HashCheckPassword(account.Password, req.Password); err != nil {
		return http.StatusUnauthorized, response, fmt.Errorf("error in comparing password : %s", err.Error())
	}

	//Fetch And Store Geo IP
	if err = u.storeLoginDetail(ctx, account.UserID); err != nil {
		return http.StatusInternalServerError, response, err
	}

	token, err := tools.JWTCreateToken(tools.JWT{
		ID:        account.UserID.String(),
		AccountID: account.ID,
		RoleName:  account.Role,
	})
	if err != nil {
		return http.StatusInternalServerError, response, fmt.Errorf("error in generate access token : %s", err.Error())
	}

	//Fetch Privilege
	privileges, err := u.repository.PrivilegeFetchByUserID(ctx, account.UserID)
	if err != nil {
		return http.StatusInternalServerError, response, fmt.Errorf("error in fetch user privilege : %s", err.Error())
	}

	return http.StatusOK, loginResponse{
		Data: loginDataResponse{
			Name:  account.Name,
			Photo: account.Photo,
		},
		Token:          token,
		Role:           account.Role,
		Privileges:     privileges,
		CanParticipate: &account.CanParticipate,
		UserID:         account.UserID,
	}, nil
}

func (u Usecase) storeLoginDetail(ctx context.Context, userID uuid.UUID) error {
	//var response userIPDetail
	//url := fmt.Sprintf("https://ipgeolocation.abstractapi.com/v1/?api_key=%s&ip_address=%s", config.Configuration().Location, clientIP)
	//if _, err := tools.HTTPRequest(tools.HTTPParams{
	//	URL:      url,
	//	Body:     nil,
	//	Method:   "GET",
	//	Headers:  nil,
	//	Response: &response,
	//}); err != nil {
	//	log.Printf("error in fetch ip detail : " + err.Error())
	//}

	if err := u.repository.LoginDetailCreate(ctx, userID); err != nil {
		return fmt.Errorf("error in store login detail : %s", err.Error())
	}
	return nil
}

func (u Usecase) callback(ctx context.Context, kind string, accessToken string) (statusCode int, response loginResponse, err error) {
	var profile profileData
	statusCode, err = u.hitSocialAuth(kind, accessToken, &profile)
	if err != nil {
		return statusCode, response, err
	}
	account, err := u.repository.AccountFetchOneByEmail(ctx, db.AccountFetchOneByEmailParams{
		Username: profile.Email,
		Type:     db.AccountType(kind),
	})
	if err != nil {
		statusCode, resp, err := u.storeFromCallback(ctx, profile, kind)
		if err != nil {
			return statusCode, response, err
		}
		account = db.AccountFetchOneByEmailRow{
			ID:     resp.AccountID,
			UserID: resp.UserID,
			Role:   db.RoleUser,
			Photo:  resp.Photo,
			Name:   resp.Name,
		}
	}

	//Fetch And Store Geo IP
	if err := u.storeLoginDetail(ctx, account.UserID); err != nil {
		return http.StatusInternalServerError, response, err
	}

	token, err := tools.JWTCreateToken(tools.JWT{
		ID:        account.UserID.String(),
		AccountID: account.ID,
		RoleName:  account.Role,
	})
	if err != nil {
		return http.StatusInternalServerError, response, fmt.Errorf("error in generate token : %s", err.Error())
	}

	//Fetch Privilege
	var privileges []string
	if account.Role != "admin" {
		privileges, err = u.repository.PrivilegeFetchByUserID(ctx, account.UserID)
		if err != nil {
			return http.StatusInternalServerError, response, fmt.Errorf("error in fetch user privilege : %s", err.Error())
		}
	}

	return http.StatusOK, loginResponse{
		Data:           loginDataResponse{Name: account.Name, Photo: account.Photo},
		Token:          token,
		Role:           account.Role,
		Privileges:     privileges,
		CanParticipate: &account.CanParticipate,
		UserID:         account.UserID,
	}, nil
}

func (u Usecase) hitSocialAuth(kind string, token string, response *profileData) (statusCode int, err error) {
	switch kind {
	case string(db.AccountTypeGoogle):
		statusCode, err = tools.HTTPRequest(tools.HTTPParams{
			URL:      config.OAuthGoogle + token,
			Body:     nil,
			Method:   http.MethodGet,
			Headers:  nil,
			Response: &response,
		})
		if err != nil {
			return statusCode, fmt.Errorf("google : %s", err.Error())
		}
	case string(db.AccountTypeFacebook):
		var profileFB profileFacebook
		statusCode, err = tools.HTTPRequest(tools.HTTPParams{
			URL:      config.OAuthFacebook + token,
			Body:     nil,
			Method:   http.MethodGet,
			Headers:  nil,
			Response: &profileFB,
		})
		if err != nil {
			return statusCode, fmt.Errorf("facebook : %s", err.Error())
		}
		*response = profileData{
			ID:       profileFB.ID,
			Email:    profileFB.Email,
			Verified: true,
			Name:     profileFB.Name,
			Photo:    profileFB.Picture.Data.URL,
		}
	default:
		return http.StatusNotFound, fmt.Errorf("login via %s is not found, login only via facebook or twitter", kind)
	}
	return
}

func (u Usecase) storeFromCallback(ctx context.Context, req profileData, kind string) (statusCode int, response storeFromCallbackResponse, err error) {
	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, response, fmt.Errorf("error in start transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	userID, err := txQuery.UserCreate(ctx, db.UserCreateParams{
		Name:           req.Name,
		Phone:          "",
		Photo:          req.Photo,
		BornAt:         "",
		BornOn:         sql.NullTime{},
		IdentityNumber: "",
		Gender:         "",
		About:          "",
	})
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return http.StatusInternalServerError, response, fmt.Errorf("rollback fail in create user tx : %s", err.Error())
		}
		return http.StatusInternalServerError, response, fmt.Errorf("error in rollback create user tx : %s", err.Error())
	}

	accountID, err := txQuery.AccountCreate(ctx, db.AccountCreateParams{
		Username: req.Email,
		Status:   true,
		Type:     db.AccountType(kind),
		UserID:   userID,
	})
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return http.StatusInternalServerError, response, fmt.Errorf("rollback fail in in rollback create account tx : %s", err.Error())
		}
		return http.StatusInternalServerError, response, fmt.Errorf("error in rollback create account tx : %s", err.Error())
	}

	if err = tx.Commit(); err != nil {
		return http.StatusInternalServerError, response, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusOK, storeFromCallbackResponse{
		Name:      req.Name,
		UserID:    userID,
		AccountID: accountID,
		Photo:     req.Photo,
	}, nil
}

func (u Usecase) check(ctx context.Context, kind string, username string) (db.UserFetchOneRow, error) {
	switch kind {
	case string(db.AccountTypeManual):
		return u.repository.UserFetchOne(ctx, username)
	case string(db.AccountTypeFacebook):
		var profileFB profileFacebook
		if _, err := tools.HTTPRequest(tools.HTTPParams{
			URL:      config.OAuthFacebook + username,
			Body:     nil,
			Method:   http.MethodGet,
			Headers:  nil,
			Response: &profileFB,
		}); err != nil {
			return db.UserFetchOneRow{}, fmt.Errorf("facebook : %s", err.Error())
		}
		return u.repository.UserFetchOne(ctx, profileFB.Email)
	case string(db.AccountTypeGoogle):
		var profile profileData
		if _, err := tools.HTTPRequest(tools.HTTPParams{
			URL:      config.OAuthGoogle + username,
			Body:     nil,
			Method:   http.MethodGet,
			Headers:  nil,
			Response: &profile,
		}); err != nil {
			return db.UserFetchOneRow{}, fmt.Errorf("google : %s", err.Error())
		}
		return u.repository.UserFetchOne(ctx, profile.Email)
	default:
		return db.UserFetchOneRow{}, fmt.Errorf("type not found. Check only based on type manual, google, or facebook")
	}
}

func (u Usecase) forgot(ctx context.Context, username string, host string, feHost string) (statusCode int, err error) {
	user, err := u.repository.UserFetchOneTypeManual(ctx, username)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch user : %s", err.Error())
	}
	token, err := tools.JWTCreateToken(tools.JWT{
		ID:        user.ID.String(),
		AccountID: user.AccountID,
		RoleName:  user.Role,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in generate token : %s", err.Error())
	}

	if err = u.rdb.Set(ctx, "forgot-"+user.Username, token.AccessToken, 24*time.Hour).Err(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in set redis keyval : %s", err.Error())
	}

	go func(data sendMailRequest) {
		buffer, err := mail.TemplateToBuffer(data.BodyPath, data.BodyParam)
		if err != nil {
			log.Printf("error in converting html file to buffer : %s \n", err.Error())
		}
		mail.SendMail(mail.Request{
			To:      user.Username,
			Subject: "[Tanding!] Recover Your Tanding! Account",
			Body:    buffer,
		})
	}(sendMailRequest{
		BodyParam: bodyParam{
			Name: user.Name,
			Path: feHost + "/forgot-password/reset?token=" + token.AccessToken,
			Host: &host,
		},
		BodyPath: "public/templates/mail/forgot-password.html",
	})
	return http.StatusOK, nil
}

func (u Usecase) reset(ctx context.Context, req resetRequest, decoded tools.JWT, token string) (statusCode int, err error) {
	if req.Password != req.ConfirmPassword {
		return http.StatusUnprocessableEntity, fmt.Errorf("password not matched")
	}
	username, err := u.repository.AccountFetchEmailByID(ctx, db.AccountFetchEmailByIDParams{
		ID:   decoded.AccountID,
		Type: db.AccountTypeManual,
	})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch account : %s", err.Error())
	}

	retrievedToken, err := u.rdb.Get(ctx, "forgot-"+username).Result()
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch data in redis : %s", err.Error())
	}
	if retrievedToken != token {
		return http.StatusForbidden, fmt.Errorf("error in comparing outcoming token with retrieved token")
	}
	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in generate hash password : %s", err.Error())
	}
	if err = u.repository.AccountUpdatePassword(ctx, db.AccountUpdatePasswordParams{
		Password: hashedPassword,
		ID:       decoded.AccountID,
	}); err != nil {
		return http.StatusNotFound, fmt.Errorf("error in update account : %s", err.Error())
	}

	if err = u.rdb.Del(ctx, "forgot-"+username).Err(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in del key redis : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) resend(ctx context.Context, username string, kind string) (statusCode int, err error) {
	switch kind {
	case "verify", "forgot", "binding":
		break
	default:
		return http.StatusNotFound, fmt.Errorf("only resend token only for verify, forgot, or binding")
	}
	user, err := u.repository.UserFetchOne(ctx, username)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch user by username : %s", err.Error())
	}
	token, err := tools.JWTCreateToken(tools.JWT{
		ID:        user.ID.String(),
		AccountID: user.AccountID,
		RoleName:  user.Role,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in create token : %s", err.Error())
	}
	if err = u.rdb.Set(ctx, kind+"-"+username, token.AccessToken, 24*time.Hour).Err(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in set redis value : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u Usecase) binding(ctx context.Context, arg bindingParams) (statusCode int, err error) {
	switch arg.Kind {
	case string(db.AccountTypeFacebook), string(db.AccountTypeGoogle):
		var profile profileData
		statusCode, err = u.hitSocialAuth(arg.Kind, arg.Request.Token, &profile)
		if err != nil {
			return statusCode, err
		}
		if _, err = u.repository.AccountFetchOneByEmail(ctx, db.AccountFetchOneByEmailParams{
			Username: profile.Email,
			Type:     db.AccountType(arg.Kind),
		}); err == nil {
			return http.StatusBadRequest, fmt.Errorf("account already binded or used in another account")
		}
		if _, err = u.repository.AccountCreate(ctx, db.AccountCreateParams{
			Username: profile.Email,
			Type:     db.AccountType(arg.Kind),
			UserID:   uuid.MustParse(arg.Decoded.ID),
			Status:   true,
		}); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in create account : %s", err.Error())
		}
	case string(db.AccountTypeManual):
		hashedPassword, err := tools.HashPassword(arg.Request.Password)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in hash password : %s", err.Error())
		}
		if _, err = u.repository.AccountFetchOneByEmail(ctx, db.AccountFetchOneByEmailParams{
			Username: arg.Request.Email,
			Type:     db.AccountType(arg.Kind),
		}); err == nil {
			return http.StatusBadRequest, fmt.Errorf("account already binded or used in another account")
		}
		accountID, err := u.repository.AccountCreate(ctx, db.AccountCreateParams{
			Username: arg.Request.Email,
			Type:     db.AccountType(arg.Kind),
			UserID:   uuid.MustParse(arg.Decoded.ID),
			Password: hashedPassword,
			Status:   false,
		})
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in create account : %s", err.Error())
		}
		token, err := tools.JWTCreateToken(tools.JWT{
			ID:        arg.Decoded.ID,
			AccountID: accountID,
			RoleName:  arg.Decoded.RoleName,
		})
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in create token : %s", err.Error())
		}
		if err := u.rdb.Set(ctx, "binding-"+arg.Request.Email, token.AccessToken, 24*time.Hour).Err(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in set redis value : %s", err.Error())
		}

		name, err := u.repository.UserFetchNameByID(ctx, uuid.MustParse(arg.Decoded.ID))
		if err != nil {
			return http.StatusNotFound, fmt.Errorf("error in get user name : %s", err.Error())
		}

		go func(data sendMailRequest) {
			buffer, err := mail.TemplateToBuffer(data.BodyPath, data.BodyParam)
			if err != nil {
				log.Printf("error in converting html file to buffer : %s \n", err.Error())
			}
			mail.SendMail(mail.Request{
				To:      arg.Request.Email,
				Subject: "[Tanding!] Binding Tanding! Account",
				Body:    buffer,
			})
		}(sendMailRequest{
			BodyParam: bodyParam{
				Name: name,
				Path: token.AccessToken,
				Host: &arg.Host,
			},
			BodyPath: "public/templates/mail/binding-account.html",
		})
	default:
		return http.StatusNotFound, fmt.Errorf("only binding by username and password, google, or facebook")
	}
	return http.StatusOK, nil
}
