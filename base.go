package paymentwall

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

const (
	/*
	 * Paymentwall library version
	 */
	LIB_VERSION = "1.0.0"
	/*
	 * URLs for Paymentwall Pro
	 */
	CHARGE_URL = "https://api.paymentwall.com/api/pro/v1/charge"
	SUBS_URL   = "https://api.paymentwall.com/api/pro/v1/subscription"

	/*
	 * API types
	 */
	API_VC    = 1
	API_GOODS = 2
	API_CART  = 3

	/*
	 * Controllers for APIs
	 */

	CONTROLLER_PAYMENT_VIRTUAL_CURRENCY = "ps"
	CONTROLELR_PAYMENT_DIGITAL_GOODS    = "subscription"
	CONTROLLER_PAYMENT_CART             = "cart"

	/**
	 * Signature versions
	 */
	DEFAULT_SIGNATURE_VERSION = 3
	SIGNATURE_VERSION_1       = 1
	SIGNATURE_VERSION_2       = 2
	SIGNATURE_VERSION_3       = 3
)

var (
	Errorsstring = []string{}
)

type PaymentwallBase struct {
	apiType   int
	appKey    string
	secretKey string
	proApiKey string
}

func NewPaymentwall(apiType int, appKey, secretKey, proApiKey string) *PaymentwallBase {
	return &PaymentwallBase{
		appKey: appKey,
		apiType:   apiType,
		secretKey: secretKey,
		proApiKey: proApiKey,  //没用
	}
}

func gethash(inputString, algorithm string) string {
	var h hash.Hash

	if algorithm == "md5" {
		h = md5.New()
	} else if algorithm == "sha256" {
		h = sha256.New()
	}
	h.Write([]byte(inputString))
	return hex.EncodeToString(h.Sum(nil))

}

func exists(value string, arr []string) bool {
	for i := range arr {
		if arr[i] == value {
			return true
		}
	}
	return false
}

func mergeMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
