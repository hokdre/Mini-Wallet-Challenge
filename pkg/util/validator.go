package util

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hokdre/mini-ewallet/internal/model"
)

type Validator interface {
	Validate(i interface{}) error
}

func NewValidator() *validatorImpl {
	impl := &validatorImpl{}
	v := validator.New()
	_ = v.RegisterValidation("enumWalletStatus", impl.validateWalletStatus)
	_ = v.RegisterValidation("enumTransactionType", impl.validateEnumTransactionType)
	_ = v.RegisterValidation("enumTransactionStatus", impl.validateEnumTransactionStatus)
	_ = v.RegisterValidation("gteNow", impl.validateDateGTENow)
	impl.validate = v
	return impl
}

type validatorImpl struct {
	validate *validator.Validate
}

func (v *validatorImpl) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *validatorImpl) validateWalletStatus(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	return value == model.WalletStatus.Enabled || value == model.WalletStatus.Disabled
}

func (v *validatorImpl) validateEnumTransactionType(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	return value == model.TransactionType.Withdrawal ||
		value == model.TransactionType.Deposit
}

func (v *validatorImpl) validateEnumTransactionStatus(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	return value == model.TransactionStatus.Pending ||
		value == model.TransactionStatus.Success ||
		value == model.TransactionStatus.Failed
}

func (v *validatorImpl) validateDateGTENow(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Interface()
	if fieldValue == nil {
		return true
	}

	value, ok := fieldValue.(*time.Time)
	if !ok {
		return false
	}

	if value == nil {
		return true
	}

	now := time.Now()
	return now.Before(*value) || now.Equal(*value)
}
