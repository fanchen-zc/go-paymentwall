package paymentwall

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
)

const (
	BASE_URL = "https://api.paymentwall.com/api"
)

type Widget struct {
	userId      string
	widgetCode  string
	products    []Product
	extraParams map[string]string
	base        PaymentwallBase
}

type WidgetFun func(*Widget)

func WithExtraParams(extraParams map[string]string) WidgetFun {
	return func(w *Widget) {
		w.extraParams = extraParams
	}
}

func NewWidget(userId, widgetCode string, products []Product, extraParams map[string]string, base PaymentwallBase, params ...WidgetFun) *Widget {
	w := &Widget{
		userId:      userId,
		products:    products,
		widgetCode:  widgetCode,
		extraParams: extraParams,
		base:        base,
	}
	for _, p := range params {
		p(w)
	}
	return w
}

func (w *Widget) getDefaultSignatureVersion() int {
	if w.base.apiType == API_CART {
		return DEFAULT_SIGNATURE_VERSION
	}
	return SIGNATURE_VERSION_2

}

func (w *Widget) GetUrl() (string, error) {
	parameters := map[string]string{}
	parameters["key"] = w.base.appKey
	parameters["uid"] = w.userId
	parameters["widget"] = w.widgetCode

	productsNumber := len(w.products)

	if w.base.apiType == API_GOODS {
		if productsNumber > 0 {
			if productsNumber == 1 {
				product := w.products[0]
				postTrialProduct := &Product{}
				switch product.trialProduct.(type) {
				case Product:
					postTrialProduct = &product
					product = product.trialProduct.(Product)
				}
				parameters["amount"] = fmt.Sprintf("%.2f", product.amount)
				parameters["currencyCode"] = product.currencyCode
				parameters["ag_name"] = product.name
				parameters["ag_external_id"] = product.productId
				parameters["ag_type"] = product.productType
				if product.productType == TYPE_SUBSCRIPTION {
					parameters["ag_period_length"] = string(rune(product.periodLength))
					parameters["ag_period_type"] = product.periodType
					if product.recurring {
						parameters["ag_recurring"] = strconv.FormatBool(product.recurring)
						if postTrialProduct != nil {
							parameters["ag_trial"] = "1"
							parameters["ag_post_trial_external_id"] = postTrialProduct.productId
							parameters["ag_post_trial_period_length"] = string(rune(postTrialProduct.periodLength))
							parameters["ag_post_trial_period_type"] = postTrialProduct.productType
							parameters["ag_post_trial_name"] = postTrialProduct.name
							parameters["post_trial_amount"] = fmt.Sprintf("%.2f", postTrialProduct.amount)
							parameters["post_trial_currencyCode"] = postTrialProduct.currencyCode
						}
					}
				}
			}
		}
	}
	if w.base.apiType == API_CART {

		for i, product := range w.products {
			parameters["external_ids["+string(rune(i))+"]"] = product.productId
			if product.amount > 0 {
				parameters["prices["+string(rune(i))+"]"] = fmt.Sprintf("%.2f", product.amount)
			}
			if product.currencyCode != "" {
				parameters["currencies["+string(rune(i))+"]"] = product.currencyCode
			}
			if product.name != "" {
				parameters["names["+string(rune(i))+"]"] = product.name
			}
		}
	}

	signatureVersion := w.getDefaultSignatureVersion()
	parameters["sign_version"] = fmt.Sprintf("%d", signatureVersion)

	if _, ok := w.extraParams["sign_version"]; ok {
		parameters["sign_version"] = w.extraParams["sign_version"]
		signatureVersion, _ = strconv.Atoi(w.extraParams["sign_version"])
	}
	parameters = mergeMaps(parameters, w.extraParams)
	parameters["sign"] = w.calculateSignature(parameters, w.base.secretKey, signatureVersion)
	controller, err := w.buildController(w.widgetCode, false)
	if err != nil {
		return "", err
	}
	return BASE_URL + "/" + controller + "?" + w.buildQueryString(parameters, "&"), nil
}

func (w *Widget) buildController(widget string, flexibleCall bool) (string, error) {
	if w.base.apiType == API_VC {
		ok, err := regexp.Match(widget, []byte("^w|s|mw"))
		if err != nil {
			return "", err
		}
		if !ok {
			return CONTROLLER_PAYMENT_VIRTUAL_CURRENCY, nil
		}
		return "", nil
	}
	if w.base.apiType == API_GOODS {
		if !flexibleCall {
			ok, err := regexp.Match(widget, []byte("^w|s|mw"))
			if err != nil {
				return "", err
			}
			if !ok {
				return CONTROLELR_PAYMENT_DIGITAL_GOODS, nil
			}
			return "", nil
		}
		return CONTROLELR_PAYMENT_DIGITAL_GOODS, nil
	}
	return CONTROLLER_PAYMENT_CART, nil
}

func (w *Widget) buildQueryString(dict map[string]string, s string) string {
	var (
		queryString      = ""
		count       int  = 0
		end         bool = false
	)
	for k, v := range dict {
		if count == len(dict)-1 {
			end = true
		}
		escapedValue := url.QueryEscape(v)
		if end {
			queryString += fmt.Sprintf(`%s=%s`, k, escapedValue)
		} else {
			queryString += fmt.Sprintf(`%s=%s%s`, k, escapedValue, s)
		}
		count++
	}
	return queryString
}

func (w *Widget) calculateSignature(parameters map[string]string, secret string, version int) string {
	baseString := ""
	if version == SIGNATURE_VERSION_1 {
		if parameters["uid"] != "" {
			baseString += parameters["uid"]
		} else {
			baseString += secret
		}
		return gethash(baseString, "md5")

	}
	keys := make([]string, 0)
	for k := range parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		baseString += fmt.Sprintf(`%s=%s`, keys[i], parameters[keys[i]])
	}
	baseString += secret
	if version == SIGNATURE_VERSION_2 {
		return gethash(baseString, "md5")
	}
	return gethash(baseString, "sha256")
}
