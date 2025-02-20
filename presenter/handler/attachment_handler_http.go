package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type attachmentHandler struct {
	useCase usecase.AttachmentUseCase
}

func NewAttachmentHandler(uc usecase.AttachmentUseCase) AttachmentHandler {
	return &attachmentHandler{useCase: uc}
}

func (h *attachmentHandler) Insert(c echo.Context) error {
	ctx := context.Background()

	var input model.AttachmentRequest

	attachableType := c.FormValue("attachable_type")

	if attachableType == "" {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "failed request")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "file not found")
	}

	input.AttachableType = attachableType

	attachment, err := h.useCase.Insert(ctx, input, file)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "success", attachment, 0, 0)
}
