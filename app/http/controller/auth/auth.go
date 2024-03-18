package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/auth"
	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/auth"
	"github.com/arif-x/sqlx-mysql-boilerplate/config"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	hash "github.com/arif-x/sqlx-mysql-boilerplate/pkg/hash"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	JWTTokenAuthed "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
)

// Register method for new user registration.
// @Description new user registration.
// @Summary new user registration.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name"
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password" format(password)
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthWithPermissionResponse
// @Router /api/v1/auth/register [post]
func Register(c *fiber.Ctx) error {
	register := &model.Register{}

	if err := c.BodyParser(register); err != nil {
		return response.BadRequest(c, err)
	}

	password, err := hash.Hash([]byte(register.Password))
	if err != nil {
		return response.InternalServerError(c, err)
	}

	register.Password = password

	repository := repo.NewAuthRepo(database.GetDB())
	user, role_name, permission, err := repository.Register(register)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	token, _, err := GenerateNewAccessToken(user.UUID, user.Username, user.Email, user.Name, user.IsActive, user.EmailVerifiedAt, user.RoleUUID, permission)
	if err != nil {
		return response.InternalServerError(c, errors.New("Internal Error"))
	}

	jwt_expired_at := config.AppCfg().JWTSecretExpireMinutesCount

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           true,
		"message":          fmt.Sprintf("Token will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"data":             token,
		"role_name":        role_name,
		"permission":       permission,
		"token_expired_at": time.Now().Add(time.Duration(jwt_expired_at) * time.Minute),
	})
}

// Login method for user login.
// @Description user login.
// @Summary user login.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "Username Or Email"
// @Param password formData string true "Password" format(password)
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthWithPermissionResponse
// @Router /api/v1/auth/login [post]
func Login(c *fiber.Ctx) error {
	login := &model.Login{}

	if err := c.BodyParser(login); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewAuthRepo(database.GetDB())
	user, role_name, permission, err := repository.Login(login.Username)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return response.InvalidCredential(c, errors.New("No credential"))
		} else {
			return response.InvalidCredential(c, errors.New("No credential"))
		}
	}

	isValid := IsValidPassword([]byte(user.Password), []byte(login.Password))
	if !isValid {
		return response.InvalidCredential(c, errors.New("Incorrect password"))
	}

	token, _, err := GenerateNewAccessToken(user.UUID, user.Username, user.Email, user.Name, user.IsActive, user.EmailVerifiedAt, user.RoleUUID, permission)
	if err != nil {
		return response.InternalServerError(c, errors.New("Internal Error"))
	}

	jwt_expired_at := config.AppCfg().JWTSecretExpireMinutesCount

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           true,
		"message":          fmt.Sprintf("Token will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"data":             token,
		"role_name":        role_name,
		"permission":       permission,
		"token_expired_at": time.Now().Add(time.Duration(jwt_expired_at) * time.Minute),
	})
}

// SendEmail method for user send verificationemail.
// @Description user send verificationemail.
// @Summary user send verificationemail.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthResponse
// @Router /api/v1/auth/send-email [post]
func SendEmail(c *fiber.Ctx) error {
	user := c.Locals("user").(*JWTTokenAuthed.Token)
	claims := user.Claims.(JWTTokenAuthed.MapClaims)
	email_token := claims["email_token"]
	user_email := claims["email"]
	email_verified_at := claims["email_verified_at"]

	if email_verified_at != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Your email has already verified!",
			"data":    nil,
		})
	}

	email, ok := user_email.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to send email!",
			"data":    nil,
		})
	}

	token, ok := email_token.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to send email!",
			"data":    nil,
		})
	}

	err := SendVerificationEmail(email, token)
	if err != nil {
		log.Println("Error sending verification email:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to send email!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Email has been sent!",
		"data":    "OK",
	})
}

// Verify method for user verify email.
// @Description user verify email.
// @Summary user verify email.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param token path string true "Email Token" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Security ApiKeyAuth
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthWithPermissionResponse
// @Router /api/v1/auth/verify/{token} [get]
func Verify(c *fiber.Ctx) error {
	token := c.Params("token")

	user := c.Locals("user").(*JWTTokenAuthed.Token)
	claims := user.Claims.(JWTTokenAuthed.MapClaims)
	username := claims["username"]
	email_token := claims["email_token"]

	if token != email_token {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid email code!",
			"data":    "OK",
		})
	}

	repository := repo.NewAuthRepo(database.GetDB())
	user_data, role_name, permission, err := repository.Verify(username.(string))

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  false,
			"message": "Can't verify your email!",
			"data":    "OK",
		})
	}

	token_data, _, err := GenerateNewAccessToken(user_data.UUID, user_data.Username, user_data.Email, user_data.Name, user_data.IsActive, user_data.EmailVerifiedAt, user_data.RoleUUID, permission)
	if err != nil {
		return response.InternalServerError(c, errors.New("Internal Error"))
	}

	jwt_expired_at := config.AppCfg().JWTSecretExpireMinutesCount

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           true,
		"message":          fmt.Sprintf("Token has been regenerated and will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"data":             token_data,
		"role_name":        role_name,
		"permission":       permission,
		"token_expired_at": time.Now().Add(time.Duration(jwt_expired_at) * time.Minute),
	})
}

// SendForgotEmailPassword method for user send email to recover password.
// @Description user send email to recover password.
// @Summary user send email to recover password.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "Email/Username" default(superadmin)
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthWithPermissionResponse
// @Router /api/v1/auth/send-password-email [post]
func SendForgotPasswordEmail(c *fiber.Ctx) error {
	forgot_password := &model.ForgotPassword{}

	if err := c.BodyParser(forgot_password); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewAuthRepo(database.GetDB())

	user_data, err := repository.ForgotPassword(forgot_password)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	token_data, forgor_password_token, err := GenerateNewAccessToken(user_data.UUID, user_data.Username, user_data.Email, user_data.Name, user_data.IsActive, user_data.EmailVerifiedAt, user_data.RoleUUID, nil)
	if err != nil {
		return response.InternalServerError(c, errors.New("Internal Error"))
	}

	err = SendPasswordEmail(user_data.Email, forgor_password_token)
	if err != nil {
		log.Println("Error sending forgot password verification email:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to send email!",
			"data":    nil,
		})
	}

	jwt_expired_at := config.AppCfg().JWTSecretExpireMinutesCount

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           true,
		"message":          fmt.Sprintf("Token has been regenerated and will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"data":             token_data,
		"is_email_sent":    true,
		"token_expired_at": time.Now().Add(time.Duration(jwt_expired_at) * time.Minute),
	})
}

// ChackForgotPasswordToken method for check token validity to recover password.
// @Description check token validity to recover password.
// @Summary check token validity to recover password.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param token path string true "Forgot Password Token" default(349af788-7702-442a-ae44-6dd6435b1339)
// @Security ApiKeyAuth
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthWithPermissionResponse
// @Router /api/v1/auth/check-forgot-password-token/{token} [get]
func CheckForgotPasswordToken(c *fiber.Ctx) error {
	token := c.Params("token")

	user := c.Locals("user").(*JWTTokenAuthed.Token)
	claims := user.Claims.(JWTTokenAuthed.MapClaims)
	forgot_password_token := claims["forgot_password_token"]

	if token != forgot_password_token {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid email code!",
			"data":    false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Valid email code!",
		"data":    true,
	})
}

// ChangeForgotPassword method for change password from forgot password.
// @Description change password from forgot password.
// @Summary change password from forgot password.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param password formData string true "New Password" format(password)
// @Security ApiKeyAuth
// @Failure 400,401,403,500 {object} response.ErrorResponse "Error"
// @Success 200 {object} response.AuthWithPermissionResponse
// @Router /api/v1/auth/change-forgot-password [post]
func ChangeForgotPassword(c *fiber.Ctx) error {
	user := c.Locals("user").(*JWTTokenAuthed.Token)
	claims := user.Claims.(JWTTokenAuthed.MapClaims)
	username := claims["username"]
	// forgot_password_token := claims["forgot_password_token"]

	change_forgot_password := &model.ChangeForgotPassword{}

	if err := c.BodyParser(change_forgot_password); err != nil {
		return response.BadRequest(c, err)
	}

	// if change_forgot_password.Password != forgot_password_token {
	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"status":  false,
	// 		"message": "Invalid email code!",
	// 		"data":    false,
	// 	})
	// }

	password, err := hash.Hash([]byte(change_forgot_password.Password))
	if err != nil {
		return response.InternalServerError(c, err)
	}

	change_forgot_password.Password = password

	repository := repo.NewAuthRepo(database.GetDB())
	user_data, role_name, permission, err := repository.ChangeForgotPassword(username.(string), change_forgot_password.Password)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  false,
			"message": "Can't change your password!",
			"data":    "OK",
		})
	}

	token_data, _, err := GenerateNewAccessToken(user_data.UUID, user_data.Username, user_data.Email, user_data.Name, user_data.IsActive, user_data.EmailVerifiedAt, user_data.RoleUUID, permission)
	if err != nil {
		return response.InternalServerError(c, errors.New("Internal Error"))
	}

	jwt_expired_at := config.AppCfg().JWTSecretExpireMinutesCount

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           true,
		"message":          fmt.Sprintf("Token has been regenerated and will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"data":             token_data,
		"role_name":        role_name,
		"permission":       permission,
		"token_expired_at": time.Now().Add(time.Duration(jwt_expired_at) * time.Minute),
	})
}

func GenerateNewAccessToken(UserID uuid.UUID, Username string, Email string, Name string, IsActive bool, EmailVerifiedAt *time.Time, RoleUUID string, Permission []string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	forgot_password_token := uuid.New().String()

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = UserID.String()
	claims["username"] = Username
	claims["email"] = Email
	claims["name"] = Name
	claims["is_active"] = IsActive
	claims["email_verified_at"] = EmailVerifiedAt
	claims["email_token"] = uuid.New().String()
	claims["forgot_password_token"] = forgot_password_token
	claims["role_uuid"] = RoleUUID
	claims["permission"] = Permission
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.AppCfg().JWTSecretExpireMinutesCount)).Unix()

	t, err := token.SignedString([]byte(config.AppCfg().JWTSecretKey))
	if err != nil {
		return "", "", err
	}

	return t, forgot_password_token, nil
}

func SendVerificationEmail(emailAddress, token string) error {
	emailConfig := &email.Email{
		From:    os.Getenv("SMTP_EMAIL_FROM"),
		To:      []string{emailAddress},
		Subject: "Email Verification",
		Text:    []byte(fmt.Sprintf("Click the following link to verify your email: %s/api/v1/auth/verify/%s", os.Getenv("APP_EMAIL_REDIRECT_URL"), token)),
	}

	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_ADDRESS"))
	return emailConfig.Send(os.Getenv("SMTP_ADDRESS")+":"+fmt.Sprint(os.Getenv("SMTP_PORT")), auth)
}

func SendPasswordEmail(emailAddress, token string) error {
	emailConfig := &email.Email{
		From:    os.Getenv("SMTP_EMAIL_FROM"),
		To:      []string{emailAddress},
		Subject: "Email Verification",
		Text:    []byte(fmt.Sprintf("Click the following link to change your password: %s/api/v1/auth/forgot-password/%s", os.Getenv("APP_EMAIL_REDIRECT_URL"), token)),
	}

	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_ADDRESS"))
	return emailConfig.Send(os.Getenv("SMTP_ADDRESS")+":"+fmt.Sprint(os.Getenv("SMTP_PORT")), auth)
}

func GeneratePasswordHash(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func IsValidPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false
	}

	return true
}
