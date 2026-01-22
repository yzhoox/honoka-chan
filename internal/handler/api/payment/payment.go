package payment

import "fmt"

func PaymentApi(action string) (res any, err error) {
	switch action {
	case "productList":
		res, err = productList()
	default:
		err = fmt.Errorf("unimplemented action: payment: %s", action)
	}
	return res, err
}
