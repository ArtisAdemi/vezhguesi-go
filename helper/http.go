package helper

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type HttpError struct {
	Message string `json:"message"`
}

func HTTPError(c *fiber.Ctx, err error, errLocation string) error {
	if err == nil || c == nil {
		log.WithFields(log.Fields{"error": err, "Ctx": c}).Error("unexpected HTTP error handling")
		return nil
	}
	statusCode := http.StatusInternalServerError
	// FIXME: this is brittle and should not be necessary.
	// If we need to differentiate between different possible error types, we should
	// create appropriate error types with clearly defined meaning.
	errStr := strings.ToLower(err.Error())
	for keyword, status := range map[string]int{
		"not found":               http.StatusNotFound,
		"bad request":             http.StatusBadRequest,
		"conflict":                http.StatusConflict,
		"impossible":              http.StatusNotAcceptable,
		"unauthorized":            http.StatusUnauthorized,
		"forbidden":               http.StatusForbidden,
		"invalid":                 http.StatusBadRequest,
		"missing":                 http.StatusBadRequest,
		"method not allowed":      http.StatusMethodNotAllowed,
		"internal error":          http.StatusInternalServerError,
		"unauthorize access role": http.StatusForbidden,
	} {
		if strings.Contains(errStr, keyword) {
			statusCode = status
			break
		}
	}
	err2 := writeJSON(c, statusCode, HttpError{Message: errStr})
	if err2 != nil {
		//just in case if json marshal fail
		c.JSON(fiber.Map{"message": http.StatusInternalServerError})
		log.WithFields(log.Fields{"error": err, "Ctx": c}).Error("unexpected HTTP error handling")
	}
	log.WithFields(log.Fields{"statusCode": statusCode, "error": err, "errLocation": errLocation}).Error("HTTP Error")
	return err2
}

// writeJSON writes the value v to the http response stream as json with standard
// json encoding.
func writeJSON(c *fiber.Ctx, code int, v interface{}) error {
	c.Set("Content-Type", "application/json")
	return c.Status(code).JSON(v)
}