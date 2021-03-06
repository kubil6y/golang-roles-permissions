package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

// envelope type allows us to send quick json responses
type envelope map[string]interface{}

// outOK() is used to send OK responses
func (app *application) outOK(message interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ok":   true,
		"data": message,
	}
}

// outERR() is used to send ERROR responses
func (app *application) outERR(message interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ok":    false,
		"error": message,
	}
}

func (app *application) writeJSON(
	w http.ResponseWriter, status int, data interface{}, headers http.Header) error {

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b = append(b, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	s := params.ByName("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id < 0 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Errorf("%s", err)
			}
		}()

		fn()
	}()
}

func (app *application) intSliceToSet(nums []int64) []int64 {
	var result []int64
	cache := map[int64]bool{}

	for _, num := range nums {
		cache[num] = true
	}

	for k := range cache {
		result = append(result, k)
	}
	return result
}

// QUERY STRING METHODS BEGIN //////////////////////////////
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func (app *application) readInt(qs url.Values, v *validator.Validator, key string, defaultValue int) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "invalid value")
		return defaultValue
	}
	return i
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func ContainsIS(nums []int64, target int64) bool {
	for _, v := range nums {
		if v == target {
			return true
		}
	}
	return false
}

// CleanedIS() compares input slice to target slice,
// returns only elements that do not exist on target slice
func CleanedIS(target, input []int64) []int64 {
	var result []int64
	for _, v := range input {
		if !ContainsIS(target, v) {
			result = append(result, v)
		}
	}
	return result
}

func ContainsPermission(list []data.Permission, target data.Permission) bool {
	for _, v := range list {
		if v.ID == target.ID {
			return true
		}
	}
	return false
}
