package domain

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrUserNotFound              = echo.NewHTTPError(http.StatusBadRequest, "User tidak ditemukan")
	ErrPhoneOrEmailNotRegistered = echo.NewHTTPError(http.StatusBadRequest, "Email atau nomor telepon tidak terdaftar")
	ErrPasswordNotSet            = echo.NewHTTPError(http.StatusBadRequest, "Password belum diatur")
	ErrPasswordIncorrect         = echo.NewHTTPError(http.StatusBadRequest, "Password yang Anda masukkan salah")
)
