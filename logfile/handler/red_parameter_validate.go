package handler

import (
    "github.com/go-playground/validator"
    "reflect"
    "strconv"
    "unicode/utf8"
    "fmt"
)

var Validate *validator.Validate

func init()  {
    Validate = validator.New()

    Validate.RegisterValidation("lte_len", LteLengthOf)
}

func LteLengthOf(fl validator.FieldLevel) bool {

    field := fl.Field()
    param := fl.Param()

    switch field.Kind() {

    case reflect.String:
        p, _ := strconv.ParseInt(param, 0, 64)

        return int64(utf8.RuneCountInString(field.String())) <= p
    }

    panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

