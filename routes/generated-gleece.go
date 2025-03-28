/*
--
This file is automatically generated. Any manual changes to this file may be overwritten.
It includes routes and handlers by the Gleece API Routes Generator.
--
Authors: Haim Kastner & Yuval Pomerchik
Generated by: Gleece Routes Generator
Generated Date: 2025-03-27
Target Engine: Gin (https://github.com/gin-gonic/gin)
--
Usage:
Refer to the Gleece documentation for details on how to use the generated routes and handlers.
--
Repository: https://github.com/gopher-fleece/gleece
--
*/
package routes
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gopher-fleece/runtime"
	RequestAuth "github.com/haimkastner/go-api-units-example/security"
	UnitsControllerImport "github.com/haimkastner/go-api-units-example/controllers"
	Param0responseQuantity "github.com/haimkastner/unitsnet-go/units"
	Param1data "github.com/haimkastner/unitsnet-go/units"
	// import extension placeholder
)
var validatorInstance = validator.New()
var urlParamRegex *regexp.Regexp
type SecurityListRelation string
const (
	SecurityListRelationAnd SecurityListRelation = "AND"
)
type SecurityCheckList struct {
	Checks   []runtime.SecurityCheck
	Relation SecurityListRelation
}
// type declarations extension placeholder
func registerEnumValidation(validate *validator.Validate, validationName string, allowedValues []string) {
	// Convert the array to a map for O(1) lookup
	lookup := make(map[string]struct{})
	for _, val := range allowedValues {
		lookup[val] = struct{}{}
	}
	// Register the custom validation
	validate.RegisterValidation(validationName, func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		_, exists := lookup[field]
		return exists
	})
}
func extractValidationErrorMessage(err error, fieldName *string) string {
	if err == nil {
		return ""
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	var errStr string
	for _, validationErr := range validationErrors {
		fName := validationErr.Field()
		if fieldName != nil {
			fName = *fieldName
		}
		errStr += fmt.Sprintf("Field '%s' failed validation with tag '%s'. ", fName, validationErr.Tag())
	}
	return errStr
}
func getStatusCode(controller runtime.Controller, hasReturnValue bool, err error) int {
	if controller.GetStatus() != nil {
		return int(*controller.GetStatus())
	}
	if err != nil {
		return http.StatusInternalServerError
	}
	if hasReturnValue {
		return http.StatusOK
	}
	return http.StatusNoContent
}
func bindAndValidateBody[TOutput any](ctx *gin.Context, contentType string, validation string, output **TOutput) error {
	var err error
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil || len(bodyBytes) == 0 {
		if strings.Contains(validation, "required") {
			return fmt.Errorf("body is required but was not provided")
		}
		return nil
	}
	var deserializedOutput TOutput
	switch contentType {
	case "application/json":
		err = json.Unmarshal(bodyBytes, &deserializedOutput)
	default:
		return fmt.Errorf("content-type %s is not currently supported by the validation subsystem", contentType)
	}
	if err != nil {
		return err
	}
	// Validate the unmarshaled data recursively
	if err = validateDataRecursive(deserializedOutput, ""); err != nil {
		return err
	}
	*output = &deserializedOutput
	return nil
}
func validateDataRecursive(data interface{}, path string) error {
	val := reflect.ValueOf(data)
	// Handle pointers by dereferencing
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		return validateDataRecursive(val.Elem().Interface(), path)
	}
	// Handle different types
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		// For slices/arrays, validate each element recursively
		for i := 0; i < val.Len(); i++ {
			elemPath := path
			if path != "" {
				elemPath = fmt.Sprintf("%s[%d]", path, i)
			} else {
				elemPath = fmt.Sprintf("[%d]", i)
			}
			// Get the element - handle case where element might be nil
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			// Validate the element recursively
			if err := validateDataRecursive(elem.Interface(), elemPath); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		// For maps, validate each value recursively
		for _, key := range val.MapKeys() {
			elemPath := path
			if path != "" {
				elemPath = fmt.Sprintf("%s.%v", path, key.Interface())
			} else {
				elemPath = fmt.Sprintf("%v", key.Interface())
			}
			elemVal := val.MapIndex(key)
			if elemVal.Kind() == reflect.Ptr && elemVal.IsNil() {
				continue
			}
			if err := validateDataRecursive(elemVal.Interface(), elemPath); err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		// Validate structs with the validator
		if err := validatorInstance.Struct(data); err != nil {
			if path != "" {
				return fmt.Errorf("validation error at %s: %w", path, err)
			}
			return err
		}
		return nil
	default:
		// Primitive types don't need validation
		return nil
	}
}
func toGinUrl(url string) string {
	processedUrl := urlParamRegex.ReplaceAllString(url, ":$1")
	processedUrl = strings.ReplaceAll(processedUrl, "//", "/")
	if processedUrl == "" {
		return "/"
	}
	if !strings.HasPrefix(processedUrl, "/") {
		processedUrl = "/" + processedUrl
	}
	return processedUrl
}
func authorize(ctx *gin.Context, checksLists []SecurityCheckList) *runtime.SecurityError {
	var lastError *runtime.SecurityError
	for _, list := range checksLists {
		if list.Relation != SecurityListRelationAnd {
			panic(
				"Encountered a security list relation of type '%s' - this is unexpected and indicates a bug in Gleece itself." +
					"Please open an issue at https://github.com/gopher-fleece/gleece/issues",
			)
		}
		// Iterate over each security list
		encounteredErrorInList := false
		for _, check := range list.Checks {
			secErr := RequestAuth.GleeceRequestAuthorization(ctx, check)
			if secErr != nil {
				lastError = secErr
				encounteredErrorInList = true
				break
			}
		}
		// If no error was encountered, validation is considered successful
		// otherwise, we continue over to the next iteration whilst keeping track of the last error
		if !encounteredErrorInList {
			return nil
		}
	}
	// If we got here it means authentication has failed
	return lastError
}
func handleAuthorizationError(ctx *gin.Context, authErr *runtime.SecurityError, operationId string) {
	statusCode := int(authErr.StatusCode)
	if authErr.CustomError != nil {
		// For now, we support JSON only
		ctx.JSON(statusCode, authErr.CustomError.Payload)
		return
	}
	stdError := runtime.Rfc7807Error{
		Type:     http.StatusText(statusCode),
		Detail:   authErr.Message,
		Status:   statusCode,
		Instance: "/gleece/authorization/error/" + operationId,
	}
	ctx.JSON(statusCode, stdError)
}
func wrapValidatorError(validatorErr error, operationId string, fieldName string) runtime.Rfc7807Error {
	return runtime.Rfc7807Error{
		Type: http.StatusText(http.StatusUnprocessableEntity),
		Detail: fmt.Sprintf(
			"A request was made to operation '%s' but parameter '%s' did not pass validation - %s",
			operationId,
			fieldName,
			extractValidationErrorMessage(validatorErr, &fieldName),
		),
		Status:   http.StatusUnprocessableEntity,
		Instance: fmt.Sprintf("/gleece/validation/error/%s", operationId),
	}
}
// function declarations extension placeholder
type MiddlewareFunc func(ctx *gin.Context) bool
type ErrorMiddlewareFunc func(ctx *gin.Context, err error) bool
var beforeOperationMiddlewares []MiddlewareFunc
var afterOperationSuccessMiddlewares []MiddlewareFunc
var onErrorMiddlewares []ErrorMiddlewareFunc
var onInputValidationMiddlewares []ErrorMiddlewareFunc
var onOutputValidationMiddlewares []ErrorMiddlewareFunc
func RegisterMiddleware(executionType runtime.MiddlewareExecutionType, middlewareFunc MiddlewareFunc) {
	switch executionType {
	case runtime.BeforeOperation:
		beforeOperationMiddlewares = append(beforeOperationMiddlewares, middlewareFunc)
	case runtime.AfterOperationSuccess:
		afterOperationSuccessMiddlewares = append(afterOperationSuccessMiddlewares, middlewareFunc)
	}
}
func RegisterErrorMiddleware(executionType runtime.ErrorMiddlewareExecutionType, errorMiddlewareFunc ErrorMiddlewareFunc) {
	switch executionType {
	case runtime.OnInputValidationError:
		onInputValidationMiddlewares = append(onInputValidationMiddlewares, errorMiddlewareFunc)
	case runtime.OnOutputValidationError:
		onOutputValidationMiddlewares = append(onOutputValidationMiddlewares, errorMiddlewareFunc)
	case runtime.OnOperationError:
		onErrorMiddlewares = append(onErrorMiddlewares, errorMiddlewareFunc)
	}
}
func RegisterCustomValidator(validateTagName string, validateFunc runtime.ValidationFunc) {
	validatorInstance.RegisterValidation(validateTagName, func(fl validator.FieldLevel) bool {
		return validateFunc(fl)
	})
}
func RegisterRoutes(engine *gin.Engine) {
	urlParamRegex = regexp.MustCompile(`\{([\w\d-_]+)\}`)
	// register routes extension placeholder
	// UnitsController
	engine.POST(toGinUrl("/units/post-unit"), func(ctx *gin.Context) {
		// route start routes extension placeholder
		authErr := authorize(
			ctx,
			[]SecurityCheckList{},
		)
		if authErr != nil {
			handleAuthorizationError(ctx, authErr, "TestUnit")
			return
		}
		controller := UnitsControllerImport.UnitsController{}
		controller.InitController(ctx)
		var conversionErr error
		var responseQuantityRawPtr *Param0responseQuantity.LengthUnits = nil
		responseQuantityRaw, isresponseQuantityExists := ctx.GetQuery("responseQuantity")
		if isresponseQuantityExists {
			responseQuantity := responseQuantityRaw
			switch responseQuantityRaw {
			case "Angstrom", "AstronomicalUnit", "Centimeter", "Chain", "DataMile", "Decameter", "Decimeter", "DtpPica", "DtpPoint", "Fathom", "Femtometer", "Foot", "Gigameter", "Hand", "Hectometer", "Inch", "Kilofoot", "KilolightYear", "Kilometer", "Kiloparsec", "Kiloyard", "LightYear", "MegalightYear", "Megameter", "Megaparsec", "Meter", "Microinch", "Micrometer", "Mil", "Mile", "Millimeter", "Nanometer", "NauticalMile", "Parsec", "Picometer", "PrinterPica", "PrinterPoint", "Shackle", "SolarRadius", "Twip", "UsSurveyFoot", "Yard":
				responseQuantityVar := Param0responseQuantity.LengthUnits(responseQuantity)
				responseQuantityRawPtr = &responseQuantityVar
			default:
				conversionErr := fmt.Errorf("responseQuantity must be one of \"Angstrom, AstronomicalUnit, Centimeter, Chain, DataMile, Decameter, Decimeter, DtpPica, DtpPoint, Fathom, Femtometer, Foot, Gigameter, Hand, Hectometer, Inch, Kilofoot, KilolightYear, Kilometer, Kiloparsec, Kiloyard, LightYear, MegalightYear, Megameter, Megaparsec, Meter, Microinch, Micrometer, Mil, Mile, Millimeter, Nanometer, NauticalMile, Parsec, Picometer, PrinterPica, PrinterPoint, Shackle, SolarRadius, Twip, UsSurveyFoot, Yard\" options only but got %s", responseQuantityRaw)
				// Middlewares onInputValidationMiddlewares section
				for _, middleware := range onInputValidationMiddlewares {
					if continueOperation := middleware(ctx, conversionErr); !continueOperation {
						return
					}
				}
				// End middlewares onInputValidationMiddlewares section
				validationError := runtime.Rfc7807Error{
					Type: http.StatusText(http.StatusUnprocessableEntity),
					Detail: fmt.Sprintf(
						"A request was made to operation 'TestUnit' but parameter '%s' was not properly sent - Expected %s but got %s",
						"responseQuantity",
						"LengthUnits",
						reflect.TypeOf(responseQuantityRaw).String(),
					),
					Status:     http.StatusUnprocessableEntity,
					Instance:   "/gleece/validation/error/TestUnit",
					Extensions: map[string]string{"error": conversionErr.Error()},
				}
				// params validation error response extension placeholder
				ctx.JSON(http.StatusUnprocessableEntity, validationError)
				return
			}
		}
		var dataRawPtr *Param1data.LengthDto = nil
		conversionErr = bindAndValidateBody(ctx, "application/json", "required", &dataRawPtr)
		if conversionErr != nil {
			// Middlewares onInputValidationMiddlewares section
			for _, middleware := range onInputValidationMiddlewares {
				if continueOperation := middleware(ctx, conversionErr); !continueOperation {
					return
				}
			}
			// End middlewares onInputValidationMiddlewares section
			validationError := runtime.Rfc7807Error{
				Type: http.StatusText(http.StatusUnprocessableEntity),
				Detail: fmt.Sprintf(
					"A request was made to operation 'TestUnit' but body parameter '%s' did not pass validation of '%s' - %s",
					"data",
					"LengthDto",
					extractValidationErrorMessage(conversionErr, nil),
				),
				Status:   http.StatusUnprocessableEntity,
				Instance: "/gleece/validation/error/TestUnit",
			}
			// json body validation error response extension placeholder
			ctx.JSON(http.StatusUnprocessableEntity, validationError)
			return
		}
		// Middlewares beforeOperationMiddlewares section
		for _, middleware := range beforeOperationMiddlewares {
			if continueOperation := middleware(ctx); !continueOperation {
				return
			}
		}
		// End middlewares beforeOperationMiddlewares section
		// before operation routes extension placeholder
		value, opError := controller.TestUnit(responseQuantityRawPtr, *dataRawPtr)
		for key, value := range controller.GetHeaders() {
			ctx.Header(key, value)
		}
		// response headers extension placeholder
		statusCode := getStatusCode(&controller, true, opError)
		if opError != nil {
			// Middlewares onErrorMiddlewares section
			for _, middleware := range onErrorMiddlewares {
				if continueOperation := middleware(ctx, opError); !continueOperation {
					return
				}
			}
			// End middlewares onErrorMiddlewares section
			stdError := runtime.Rfc7807Error{
				Type:       http.StatusText(statusCode),
				Detail:     "Encountered an error during operation 'TestUnit'",
				Status:     statusCode,
				Instance:   "/gleece/controller/error/TestUnit",
				Extensions: map[string]string{"error": opError.Error()},
			}
			// json error response extension placeholder
			ctx.JSON(statusCode, stdError)
			return
		}
		// json response extension placeholder
		// Middlewares afterOperationSuccessMiddlewares section
		for _, middleware := range afterOperationSuccessMiddlewares {
			if continueOperation := middleware(ctx); !continueOperation {
				return
			}
		}
		// End middlewares afterOperationSuccessMiddlewares section
		// after operation routes extension placeholder
		ctx.JSON(statusCode, value)
		// route end routes extension placeholder
	})
}
