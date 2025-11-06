package handlers

import (
	"net/http"
	"strconv"

	"flypro-assessment-ayo/internal/dto"
	"flypro-assessment-ayo/internal/services"
	"flypro-assessment-ayo/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService services.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService services.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		handleValidationError(c, err)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr)
			return
		}
		utils.ErrorResponse(c, utils.NewInternalError("failed to create user", err))
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError("invalid user ID"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr)
			return
		}
		utils.ErrorResponse(c, utils.NewInternalError("failed to get user", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

func handleValidationError(c *gin.Context, err error) {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			tag := fieldErr.Tag()
			var message string

			switch tag {
			case "required":
				message = field + " is required"
			case "email":
				message = field + " must be a valid email"
			case "min":
				message = field + " must be at least " + fieldErr.Param() + " characters"
			case "max":
				message = field + " must be at most " + fieldErr.Param() + " characters"
			case "gt":
				message = field + " must be greater than " + fieldErr.Param()
			case "len":
				message = field + " must be exactly " + fieldErr.Param() + " characters"
			case "currency":
				message = field + " must be a valid currency code"
			case "oneof":
				message = field + " must be one of: " + fieldErr.Param()
			default:
				message = field + " is invalid"
			}

			errors[field] = message
		}
		utils.ValidationErrorResponse(c, errors)
	} else {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
	}
}