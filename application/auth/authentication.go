package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/lripardo/lrw/domain/api"
	"github.com/lripardo/lrw/domain/auth"
	"net/http"
)

const UserNotFoundMessage = "user not found"

var (
	PasswordCost                 = api.NewKey("AUTH_APP_PASSWORD_COST", "gte=4,lte=31", "10")
	AppShowUserExists            = api.NewKey("AUTH_APP_SHOW_USER_EXISTS", api.Boolean, "true")
	AppShowUserNotFound          = api.NewKey("AUTH_APP_SHOW_USER_NOTFOUND", api.Boolean, "true")
	AppAllowLoginFromNotVerified = api.NewKey("AUTH_APP_ALLOW_LOGIN_FROM_NOT_VERIFIED", api.Boolean, "false")
	AppInputValidator            = api.NewKey("AUTH_APP_INPUT_VALIDATOR", "required", "default")
)

type App struct {
	userVerify     *auth.UserVerify
	resetPassword  *auth.ResetPassword
	validate       *validator.Validate
	authentication *auth.Authentication
	inputValidator api.InputValidator
	userRepository auth.UserRepository

	showUserExists            bool
	passwordCost              int
	showUserNotFound          bool
	allowLoginFromNotVerified bool
}

type LoginUserData struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"password-default"`
	Cookie   bool   `json:"cookie"`
	Token    string `json:"token"`
}

type RegisterUserData struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Password  string `json:"password" validate:"password-default"`
	Token     string `json:"token"`
}

type SendVerifyData struct {
	Email string `validate:"required,email,max=255"`
	Token string `json:"token"`
}

type VerifyData struct {
	JWT   string `json:"jwt" validate:"required,jwt"`
	Token string `json:"token"`
}

type SendResetPasswordData struct {
	Email string `validate:"required,email,max=255"`
	Token string `json:"token"`
}

type ResetPasswordData struct {
	JWT      string `json:"jwt" validate:"required,jwt"`
	Password string `json:"password" validate:"password-default"`
	Token    string `json:"token"`
}

type ChangePasswordData struct {
	CurrentPassword string `json:"current_password" validate:"password-default"`
	Password        string `json:"password" validate:"password-default"`
	Token           string `json:"token"`
}

func (u *App) Register(ctx api.Context) *api.Response {
	input := RegisterUserData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	exists, err := u.userRepository.UserExists(input.Email)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if exists {
		if u.showUserExists {
			return api.ResponseConflict()
		}
		api.D("user not exists")
		return api.ResponseInvalidInput()
	}
	user, err := auth.NewUser(input.Email, input.FirstName, input.LastName, input.Password, u.passwordCost)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if err := u.userRepository.CreateUser(user); err != nil {
		return api.ResponseInternalError(err)
	}
	u.userVerify.Start(user)
	return api.ResponseOk()
}

func (u *App) internalLogin(email, password string) (*auth.User, *api.Response) {
	user, err := u.userRepository.ReadUser(email, "email", "password", "verified_on")
	if err != nil {
		return nil, api.ResponseInternalError(err)
	}
	if user == nil {
		if u.showUserNotFound {
			return nil, api.ResponseNotFound()
		}
		api.D(UserNotFoundMessage)
		return nil, api.ResponseUnauthorized()
	}
	if !user.IsPassword(password) {
		return nil, api.ResponseUnauthorized()
	}
	if !u.allowLoginFromNotVerified && user.VerifiedOn == nil {
		return nil, api.ResponseForbidden()
	}
	return user, nil
}

func (u *App) Login(ctx api.Context) *api.Response {
	input := LoginUserData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	user, response := u.internalLogin(input.Email, input.Password)
	if response != nil {
		return response
	}
	authData, err := u.authentication.SignUser(ctx, user, input.Cookie)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if err := u.userRepository.UpdateTimestamp(user.Email, "last_login"); err != nil {
		return api.ResponseInternalError(err)
	}
	return api.ResponseOkWithData(map[string]interface{}{
		"auth": authData,
	})
}

func (u *App) Logout(ctx api.Context) *api.Response {
	user := auth.GetUserContext(ctx)
	if err := u.authentication.SignOutUser(ctx, user); err != nil {
		return api.ResponseInternalError(err)
	}
	return api.ResponseOk()
}

func (u *App) Info(ctx api.Context) *api.Response {
	user := auth.GetUserContext(ctx)
	return api.ResponseOkWithData(map[string]interface{}{
		"user": user,
	})
}

func (u *App) Verify(ctx api.Context) *api.Response {
	input := VerifyData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	claims, err := u.userVerify.Verify(input.JWT)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if claims == nil {
		return api.ResponseUnauthorized()
	}
	if err := u.userRepository.MakeVerified(claims.Subject); err != nil {
		return api.ResponseInternalError(err)
	}
	return api.ResponseOk()
}

func (u *App) SendVerify(ctx api.Context) *api.Response {
	input := SendVerifyData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	user, err := u.userRepository.ReadUser(input.Email, "email", "first_name", "verified_on")
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if user == nil {
		if u.showUserNotFound {
			return api.ResponseNotFound()
		}
		api.D(UserNotFoundMessage)
		return api.ResponseOk()
	}
	u.userVerify.Start(user)
	return api.ResponseOk()
}

func (u *App) internalChangePassword(user *auth.User, password string) *api.Response {
	if err := u.authentication.DeleteUserContext(user.Email); err != nil {
		return api.ResponseInternalError(err)
	}
	if err := user.HashPassword(password, u.passwordCost); err != nil {
		return api.ResponseInternalError(err)
	}
	if err := u.userRepository.UpdatePassword(user); err != nil {
		return api.ResponseInternalError(err)
	}
	return nil
}

func (u *App) ResetPassword(ctx api.Context) *api.Response {
	input := ResetPasswordData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	claims, err := u.resetPassword.ResetPassword(input.JWT)
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if claims == nil {
		api.D("invalid claims found, token empty, expired or invalid")
		return api.ResponseUnauthorized()
	}
	user, err := u.userRepository.ReadUser(claims.Subject, "email", "password", "last_change_password")
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if allow := u.resetPassword.AllowChangePassword(user); !allow {
		api.D("password changed recently, wait until the token reset password life expires plus 1 minute")
		return api.ResponseForbidden()
	}
	if err := u.userRepository.MakeVerified(user.Email); err != nil {
		return api.ResponseInternalError(err)
	}
	if response := u.internalChangePassword(user, input.Password); response != nil {
		return response
	}
	return api.ResponseOk()
}

func (u *App) SendResetPassword(ctx api.Context) *api.Response {
	input := SendResetPasswordData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	user, err := u.userRepository.ReadUser(input.Email, "email", "first_name")
	if err != nil {
		return api.ResponseInternalError(err)
	}
	if user == nil {
		if u.showUserNotFound {
			return api.ResponseNotFound()
		}
		api.D(UserNotFoundMessage)
		return api.ResponseOk()
	}
	u.resetPassword.Start(user)
	return api.ResponseOk()
}

func (u *App) ChangePassword(ctx api.Context) *api.Response {
	input := ChangePasswordData{}
	if response := u.inputValidator.Read(ctx, &input, u.validate); response != nil {
		return response
	}
	userContext := auth.GetUserContext(ctx)
	user, response := u.internalLogin(userContext.Email, input.CurrentPassword)
	if response != nil {
		return response
	}
	if response := u.internalChangePassword(user, input.Password); response != nil {
		return response
	}
	return api.ResponseOk()
}

func (u *App) Routes() []api.Route {
	root := api.NewRootRoute("api/v1")

	authRoute := root.Append(api.Route{
		Path: "auth",
	})

	register := authRoute.Append(api.Route{
		Path:     "register",
		Methods:  []string{http.MethodPost},
		Handlers: []api.Handler{u.Register},
	})

	login := authRoute.Append(api.Route{
		Path:     "login",
		Methods:  []string{http.MethodPost},
		Handlers: []api.Handler{u.Login},
	})

	verify := authRoute.Append(api.Route{
		Path:     "verify",
		Methods:  []string{http.MethodPost},
		Handlers: []api.Handler{u.Verify},
	})

	sendVerify := authRoute.Append(api.Route{
		Path:     "send-verify",
		Methods:  []string{http.MethodPost},
		Handlers: []api.Handler{u.SendVerify},
	})

	sendResetPassword := authRoute.Append(api.Route{
		Path:     "send-reset-password",
		Methods:  []string{http.MethodPost},
		Handlers: []api.Handler{u.SendResetPassword},
	})

	resetPassword := authRoute.Append(api.Route{
		Path:     "reset-password",
		Methods:  []string{http.MethodPost},
		Handlers: []api.Handler{u.ResetPassword},
	})

	// authenticated handlers
	logout := authRoute.Append(api.Route{
		Path:    "logout",
		Methods: []string{http.MethodPost},
		Handlers: []api.Handler{
			u.authentication.Authenticate(),
			u.Logout,
		},
	})

	info := authRoute.Append(api.Route{
		Path:    "",
		Methods: []string{http.MethodGet},
		Handlers: []api.Handler{
			u.authentication.Authenticate(),
			u.Info,
		},
	})

	changePassword := authRoute.Append(api.Route{
		Path:    "change-password",
		Methods: []string{http.MethodPost},
		Handlers: []api.Handler{
			u.authentication.Authenticate(),
			u.ChangePassword,
		},
	})

	return []api.Route{
		register,
		login,
		info,
		verify,
		sendVerify,
		logout,
		sendResetPassword,
		resetPassword,
		changePassword,
	}
}

func NewApp(configuration api.Configuration,
	authentication *auth.Authentication,
	repository auth.UserRepository,
	verify *auth.UserVerify,
	reset *auth.ResetPassword,
	factory api.InputValidationFactory,
) api.App {
	validate := api.NewValidator(api.IsPasswordDefault())
	inputValidator := configuration.String(AppInputValidator)
	input := factory.InputValidator(inputValidator)

	showExists := configuration.Bool(AppShowUserExists)
	passwordCost := configuration.Int(PasswordCost)
	showNotFound := configuration.Bool(AppShowUserNotFound)
	allowLoginNotVerified := configuration.Bool(AppAllowLoginFromNotVerified)

	return &App{
		validate:                  validate,
		userVerify:                verify,
		resetPassword:             reset,
		authentication:            authentication,
		inputValidator:            input,
		userRepository:            repository,
		showUserExists:            showExists,
		passwordCost:              passwordCost,
		showUserNotFound:          showNotFound,
		allowLoginFromNotVerified: allowLoginNotVerified,
	}
}
