package operator

import (
	"course_project/internal/dto"
	"course_project/internal/services"
	"course_project/internal/services/operator"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetQueuedDialogs(ctx *fiber.Ctx) error {
	dialogs, err := h.svc.QueuedDialog.ListAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load queued dialogs",
		})
	}
	return ctx.JSON(dialogs)
}

func (h *Handler) GetActiveDialogs(ctx *fiber.Ctx) error {
	dialogs, err := h.svc.ActiveDialog.ListAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load active dialogs",
		})
	}
	return ctx.JSON(dialogs)
}

func (h *Handler) ListOperators(ctx *fiber.Ctx) error {
	operators, err := h.svc.Operator.GetAllOperators(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch operators")
	}
	return ctx.JSON(operators)
}

func (h *Handler) CreateOperator(ctx *fiber.Ctx) error {
	var dto dto.CreateOperatorDTO
	if err := ctx.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	op, err := h.svc.Operator.AddOperator(ctx.Context(), dto)
	if err != nil {
		if errors.Is(err, operator.ErrOperatorAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, "Operator already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(op)
}

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
