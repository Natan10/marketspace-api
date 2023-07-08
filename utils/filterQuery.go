package utils

import (
	"fmt"
	"strings"
)

type Params map[string]interface{}

func mountPaymentMethodsQuery(params []string) string {
	args := []string{}
	for _, v := range params {
		args = append(args, fmt.Sprintf("%v='true'", v))
	}
	return strings.Join(args, " and ")
}

func MountFilterQuery(params Params) string {
	args := []string{}

	for k, v := range params {
		if k == "paymentMethods" {
			paymentMethods := params["paymentMethods"].([]string)
			args = append(args, mountPaymentMethodsQuery(paymentMethods))
		} else {
			args = append(args, fmt.Sprintf("%v='%v'", k, v))
		}
	}

	return strings.Join(args, " and ")
}
