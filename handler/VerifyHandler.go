package handler

import (
	"WalletVerifyDemo/types"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	isOwnerEth "github.com/nieben/is-owner/eth"
	"net/http"
)

type VerifyHandler struct{}

func (h *VerifyHandler) Message(c echo.Context) error {
	msg, err := isOwnerEth.Message()
	if err != nil {
		c.Logger().Errorf("VerifyHandler Message get message err: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	return c.JSON(http.StatusOK, types.MessageResponse{Message: base64.StdEncoding.EncodeToString(msg)})
}

func (h *VerifyHandler) Verify(c echo.Context) error {
	r := new(types.VerifyRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !common.IsHexAddress(r.Address) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid address"))
	}
	msg, _ := base64.StdEncoding.DecodeString(r.Message)
	signedMsg, _ := base64.StdEncoding.DecodeString(r.SignedMessage)

	rt, err := isOwnerEth.Verify(r.Address, msg, signedMsg)
	if err != nil {
		c.Logger().Errorf("VerifyHandler Verify err: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, types.VerifyResponse{Verified: rt})
}
