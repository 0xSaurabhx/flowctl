package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateAPIToken(c echo.Context) error {
	user, err := h.getUserInfo(c)
	if err != nil {
		return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
	}

	var req APITokenCreateReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "could not decode request", err, nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return wrapError(ErrValidationFailed, "label is required", err, nil)
	}

	raw, model, err := h.co.CreateAPIToken(c.Request().Context(), user.ID, req.Label)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not create api token", err, nil)
	}

	return c.JSON(http.StatusCreated, APITokenCreateResp{
		APIToken: model,
		Token:    raw,
	})
}

func (h *Handler) HandleListAPITokens(c echo.Context) error {
	user, err := h.getUserInfo(c)
	if err != nil {
		return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
	}

	tokens, err := h.co.ListAPITokens(c.Request().Context(), user.ID)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not list api tokens", err, nil)
	}
	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) HandleRevokeAPIToken(c echo.Context) error {
	user, err := h.getUserInfo(c)
	if err != nil {
		return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
	}

	tokenUUID := c.Param("tokenID")
	if tokenUUID == "" {
		return wrapError(ErrInvalidInput, "token id is required", nil, nil)
	}

	if err := h.co.RevokeAPIToken(c.Request().Context(), user.ID, tokenUUID); err != nil {
		return wrapError(ErrOperationFailed, "could not revoke api token", err, nil)
	}
	return c.NoContent(http.StatusNoContent)
}
