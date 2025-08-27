package validation

import (
	"github.com/go-playground/validator/v10"
)

// allowedAlgorithms
var allowedAlgorithms = map[string]struct{}{
	"RSA": {},
	"ECC": {},
}

// registerRules
func registerRules(v *validator.Validate) {
	// algorithm = "RSA" | "ECC"
	_ = v.RegisterValidation("algorithm", func(fl validator.FieldLevel) bool {
		_, ok := allowedAlgorithms[fl.Field().String()]
		return ok
	})
}
