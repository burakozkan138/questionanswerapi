package middleware

import (
	"burakozkan138/questionanswerapi/internal/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate, translator = initilazeValidator()
)

func Validation(next http.Handler, obj interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			response := models.NewErrorResponse(false, "Invalid JSON", http.StatusBadRequest, nil, nil)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		json.Unmarshal(bodyBytes, obj)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = validate.Struct(obj)
		if err != nil {
			errors := TranslateErrors(err.(validator.ValidationErrors))

			response := models.NewErrorResponse(false, "Validation error", http.StatusBadRequest, nil, errors)

			w.Header().Set("Content-Type", "application/json")
			w.Write(response.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TranslateErrors(e validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)

	for _, err := range e {
		errors[err.Field()] = err.Translate(translator)
	}

	return errors
}

func initilazeValidator() (*validator.Validate, ut.Translator) {
	val := validator.New()

	english := english.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en_us")

	_ = en_translations.RegisterDefaultTranslations(val, trans)

	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return val, trans
}
