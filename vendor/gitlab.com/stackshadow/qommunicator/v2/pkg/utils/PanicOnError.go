package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type perr struct {
	Err           error // the original error
	StatusOnError int   // you can override the status code
}

func (err perr) Error() string {
	return err.Err.Error()
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicOnErrWithStatus(err error, statusOnError int) {
	if err != nil {
		panic(perr{
			Err:           err,
			StatusOnError: statusOnError,
		})
	}
}

func PanicWithString(message string) {
	panic(errors.New(message))
}

type ResponseError struct {
	Error      string `json:"error"`
	StackTrace string `json:"stacktrace"`
}

func (err ResponseError) MarshalZerologObject(e *zerolog.Event) {
	e.Str("error", err.Error)
	e.Str("stack", err.StackTrace)
}

func RespondOnPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		RespondWithPanicInterface(r, w)
	}
}

// interf can be an normal error, or an perr-struct ( created by PanicOnErrWithStatus)
func RespondWithPanicInterface(interf interface{}, w http.ResponseWriter) {

	var stackTrace []byte = make([]byte, 1024)
	runtime.Stack(stackTrace, true)

	var newResponseError ResponseError
	newResponseError.StackTrace = string(stackTrace)

	// write the status code
	if perr, success := interf.(perr); success {
		// This need to be called before write
		if perr.StatusOnError == 0 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(perr.StatusOnError)
		}
	} else {
		// This need to be called before write
		w.WriteHeader(http.StatusInternalServerError)
	}

	// an error should be returned
	if err, success := interf.(error); success {
		newResponseError.Error = err.Error()
	}

	jsonBytes, err := json.Marshal(&newResponseError)
	if err == nil {
		w.Write(jsonBytes)
		log.Error().EmbedObject(newResponseError).Send()
	} else {
		w.Write([]byte{})
	}
}
