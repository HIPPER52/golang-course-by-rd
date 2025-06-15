package operator

import (
	"course_project/internal/constants"
	"course_project/internal/dto"
	"course_project/internal/repository"
	"course_project/internal/services"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

// GetQueuedDialogs godoc
// @Summary      Get queued dialogs
// @Description  Returns all queued dialogs
// @Tags         operator
// @Produce      json
// @Success      200  {array}  models.QueuedDialog
// @Failure      500  {object}  map[string]interface{}
// @Router       /operator/dialogs/queued [get]
func (h *Handler) GetQueuedDialogs(ctx *fiber.Ctx) error {
	dialogs, err := h.svc.QueuedDialog.ListAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load queued dialogs",
		})
	}
	return ctx.JSON(dialogs)
}

// GetActiveDialogs godoc
// @Summary      Get active dialogs of operator
// @Description  Returns all active dialogs assigned to the logged-in operator
// @Tags         operator
// @Produce      json
// @Success      200  {array}  models.ActiveDialog
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /operator/dialogs/active [get]
func (h *Handler) GetActiveDialogs(ctx *fiber.Ctx) error {
	userIDRaw := ctx.Locals(constants.CONTEXT_USER_ID)
	userIDPtr, ok := userIDRaw.(*string)
	if !ok || userIDPtr == nil || *userIDPtr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	dialogs, err := h.svc.ActiveDialog.FindByOperatorID(ctx.Context(), *userIDPtr)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load active dialogs",
		})
	}
	return ctx.JSON(dialogs)
}

// ListOperators godoc
// @Summary      List all operators
// @Description  Returns all operators in the system
// @Tags         operator
// @Produce      json
// @Success      200  {array}  operator.Operator
// @Failure      500  {object}  map[string]interface{}
// @Router       /operator/list [get]
func (h *Handler) ListOperators(ctx *fiber.Ctx) error {
	operators, err := h.svc.Operator.GetAllOperators(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch operators")
	}
	return ctx.JSON(operators)
}

// CreateOperator godoc
// @Summary      Create operator
// @Description  Creates a new operator account
// @Tags         operator
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateOperatorDTO  true  "Operator Data"
// @Success      201      {object}  operator.Operator
// @Failure      400      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /operator/create [post]
func (h *Handler) CreateOperator(ctx *fiber.Ctx) error {
	var dto dto.CreateOperatorDTO
	if err := ctx.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	op, err := h.svc.Operator.AddOperator(ctx.Context(), dto)
	if err != nil {
		if errors.Is(err, repository.ErrOperatorAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, "Operator already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(op)
}

// GetOperatorStats godoc
// @Summary      Get operator statistics
// @Description  Returns statistics for each operator including dialog counts and average duration
// @Tags         operator
// @Produce      json
// @Success      200  {array}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /operator/stats [get]
func (h *Handler) GetOperatorStats(ctx *fiber.Ctx) error {
	operators, err := h.svc.Operator.GetAllOperators(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load operators")
	}

	type Stats struct {
		ID                     string `json:"id"`
		Username               string `json:"username"`
		Email                  string `json:"email"`
		Role                   string `json:"role"`
		TotalDialogs           int    `json:"total_dialogs"`
		ArchivedDialogs        int    `json:"archived_dialogs"`
		AverageDurationMinutes int    `json:"average_duration_minutes"`
	}

	stats := make([]Stats, 0, len(operators))

	for _, op := range operators {
		activeCount, _ := h.svc.ActiveDialog.CountByOperator(ctx.Context(), op.ID)
		archivedDialogs, _ := h.svc.ArchivedDialog.FindByOperator(ctx.Context(), op.ID)

		totalDuration := 0
		archivedCount := 0

		for _, dlg := range archivedDialogs {
			if !dlg.StartedAt.IsZero() && !dlg.EndedAt.IsZero() {
				duration := dlg.EndedAt.Sub(dlg.StartedAt)
				totalDuration += int(duration.Minutes())
				archivedCount++
			}
		}

		avgDuration := 0
		if archivedCount > 0 {
			avgDuration = totalDuration / archivedCount
		}

		stats = append(stats, Stats{
			ID:                     op.ID,
			Username:               op.Username,
			Email:                  op.Email,
			Role:                   string(op.Role),
			TotalDialogs:           activeCount + len(archivedDialogs),
			ArchivedDialogs:        len(archivedDialogs),
			AverageDurationMinutes: avgDuration,
		})
	}

	return ctx.JSON(stats)
}
