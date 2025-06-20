package auth

import (
	"course_project/cmd/server/parser"
	"course_project/internal/constants/roles"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/services"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{
		svc: svc,
	}
}

type SignUpRequestBody struct {
	Username string     `json:"username" validate:"required,min=2,max=50,alphanum" example:"John Doe"`
	Email    string     `json:"email" validate:"required,min=5,max=50" example:"test@test.com"`
	Password string     `json:"password" validate:"required,min=8,max=20" example:"12345678"`
	Role     roles.Role `json:"role" example:"operator"`
}

type SignUpResponse200Body struct {
	Operator *models.Operator `json:"operator"`
}

// SignUp godoc
// @Summary      Register a new operator
// @Description  Creates a new operator with the given credentials and role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      SignUpRequestBody  true  "Registration payload"
// @Success      200      {object}  SignUpResponse200Body
// @Failure      400      {string}  string  "Bad request"
// @Failure      409      {string}  string  "Operator already exists"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /auth/sign-up [post]
func (h *Handler) SignUp(ctx *fiber.Ctx) error {
	req := &SignUpRequestBody{}
	err := parser.ParseBody(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pwdHash, err := h.svc.Auth.GeneratePasswordHash(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	usr := dto.CreateOperatorDTO{
		Username: req.Username,
		Email:    req.Email,
		PwdHash:  pwdHash,
		Role:     req.Role,
	}

	opr, err := h.svc.Operator.AddOperator(ctx.Context(), usr)

	if errors.Is(err, repository.ErrOperatorAlreadyExists) {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&SignUpResponse200Body{
		Operator: opr,
	})
}

type SignInRequestBody struct {
	Email    string `json:"email" validate:"required,min=5,max=50" example:"test@test.com"`
	Password string `json:"password" validate:"required,min=8,max=20" example:"12345678"`
}

type SignInResponse200Body struct {
	Token      string     `json:"token"`
	OperatorID string     `json:"operator_id"`
	Role       roles.Role `json:"role"`
}

// SignIn godoc
// @Summary      Sign in as an operator
// @Description  Authenticates the operator and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      SignInRequestBody  true  "Login credentials"
// @Success      200      {object}  SignInResponse200Body
// @Failure      400      {string}  string  "Bad request"
// @Failure      401      {string}  string  "Unauthorized"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /auth/sign-in [post]
func (h *Handler) SignIn(ctx *fiber.Ctx) error {
	req := &SignInRequestBody{}

	err := parser.ParseBody(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	oprtr, err := h.svc.Operator.GetOperatorByEmail(ctx.Context(), req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res, err := h.svc.Auth.CompareHashAndPassword(req.Password, oprtr.PwdHash)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !res {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	token, err := h.svc.Auth.CreateAuthToken(oprtr.ID, oprtr.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	return ctx.JSON(&SignInResponse200Body{
		Token:      *token,
		OperatorID: oprtr.ID,
		Role:       oprtr.Role,
	})
}
