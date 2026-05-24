package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func Logger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			traceID := req.Header.Get("X-Trace-Id")
			if traceID == "" {
				id, _ := uuid.NewRandom()
				traceID = id.String()
			}
			c.Set("traceID", traceID)

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			logger.Info("request",
				slog.String("trace_id", traceID),
				slog.String("method", req.Method),
				slog.String("uri", req.RequestURI),
				slog.Int("status", res.Status),
				slog.Duration("latency", time.Since(start)),
				slog.String("user_agent", req.UserAgent()),
			)

			return err
		}
	}
}

func ErrorHandler(logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := "Internal Server Error"
		
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message.(string)
		}

		logger.Error("server error",
			slog.String("trace_id", c.Get("traceID").(string)),
			slog.String("error", err.Error()),
		)

		c.JSON(code, map[string]interface{}{
			"error": map[string]string{
				"code":    http.StatusText(code),
				"message": msg,
			},
		})
	}
}
