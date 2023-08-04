package paymentwall

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

const (
	/**
	 * Pingback types
	 */
	PINGBACK_TYPE_REGULAR  = 0
	PINGBACK_TYPE_GOODWILL = 1
	PINGBACK_TYPE_NEGATIVE = 2

	PINGBACK_TYPE_RISK_UNDER_REVIEW      = 200
	PINGBACK_TYPE_RISK_REVIEWED_ACCEPTED = 201
	PINGBACK_TYPE_RISK_REVIEWED_DECLINED = 202

	PINGBACK_TYPE_RISK_AUTHORIZATION_VOIDED = 203

	PINGBACK_TYPE_SUBSCRIPTION_CANCELLATION   = 12
	PINGBACK_TYPE_SUBSCRIPTION_EXPIRED        = 13
	PINGBACK_TYPE_SUBSCRIPTION_PAYMENT_FAILED = 14
)

type pingback struct {
	parameters map[string]string
	ipAddress  string
	apiType    int
	secretKey  string
}

func Newpingpack(parameters map[string]string, ipAddress string,apiType int, secretKey string) *pingback {
	return &pingback{
		parameters: parameters,
		ipAddress:  ipAddress,
		apiType: apiType,
		secretKey: secretKey,
	}
}

func (p *pingback) Validate(skipIpWhiteListCheck bool) bool {
	validated := false
	if p.IsParametersValid() {

		if p.isIpAddressValid() || skipIpWhiteListCheck {
			validated = true
			if ok := p.isSignatureValid(); ok {
				validated = true
			} else {
				Errorsstring = append(Errorsstring, "Wrong signature")
			}
		} else {
			Errorsstring = append(Errorsstring, "IP address is not whitelisted")
		}
	} else {

		Errorsstring = append(Errorsstring, "Missing parameters")
	}
	return validated
}

/**
 * @description:
 * @return {*}
 */
func (p *pingback) IsParametersValid() bool {
	var (
		errorsNumber   = 0
		requiredParams = []string{}
	)
	if p.apiType == API_VC {
		requiredParams = []string{"uid", "currency", "type", "ref", "sig"}
	}
	if p.apiType == API_GOODS {
		requiredParams = []string{"uid", "goodsid", "type", "ref", "sig"}
	} else { //API_CART
		requiredParams = []string{"uid", "goodsid", "type", "ref", "sig"}
	}
	for _, v := range requiredParams {
		if val, ok := p.parameters[v]; !ok || len(val) <= 0 {
			Errorsstring = append(Errorsstring, "Parameter "+v+" is missing")
			errorsNumber++
		}
	}

	return errorsNumber == 0
}

func (p *pingback) isIpAddressValid() bool {
	ipWhitelist := []string{
		"174.36.92.186",
		"174.36.96.66",
		"174.36.92.187",
		"174.36.92.192",
		"174.37.14.28",
		"216.127.71.0",
		"216.127.71.1",
		"216.127.71.2",
		"216.127.71.3",
		"216.127.71.4",
		"216.127.71.5",
		"216.127.71.6",
		"216.127.71.7",
		"216.127.71.8",
		"216.127.71.9",
		"216.127.71.10",
		"216.127.71.11",
		"216.127.71.12",
		"216.127.71.13",
		"216.127.71.14",
		"216.127.71.15",
		"216.127.71.16",
		"216.127.71.17",
		"216.127.71.18",
		"216.127.71.19",
		"216.127.71.20",
		"216.127.71.21",
		"216.127.71.22",
		"216.127.71.23",
		"216.127.71.24",
		"216.127.71.25",
		"216.127.71.26",
		"216.127.71.27",
		"216.127.71.28",
		"216.127.71.29",
		"216.127.71.30",
		"216.127.71.31",
		"216.127.71.32",
		"216.127.71.33",
		"216.127.71.34",
		"216.127.71.35",
		"216.127.71.36",
		"216.127.71.37",
		"216.127.71.38",
		"216.127.71.39",
		"216.127.71.40",
		"216.127.71.41",
		"216.127.71.42",
		"216.127.71.43",
		"216.127.71.44",
		"216.127.71.45",
		"216.127.71.46",
		"216.127.71.47",
		"216.127.71.48",
		"216.127.71.49",
		"216.127.71.50",
		"216.127.71.51",
		"216.127.71.52",
		"216.127.71.53",
		"216.127.71.54",
		"216.127.71.55",
		"216.127.71.56",
		"216.127.71.57",
		"216.127.71.58",
		"216.127.71.59",
		"216.127.71.60",
		"216.127.71.61",
		"216.127.71.62",
		"216.127.71.63",
		"216.127.71.64",
		"216.127.71.65",
		"216.127.71.66",
		"216.127.71.67",
		"216.127.71.68",
		"216.127.71.69",
		"216.127.71.70",
		"216.127.71.71",
		"216.127.71.72",
		"216.127.71.73",
		"216.127.71.74",
		"216.127.71.75",
		"216.127.71.76",
		"216.127.71.77",
		"216.127.71.78",
		"216.127.71.79",
		"216.127.71.80",
		"216.127.71.81",
		"216.127.71.82",
		"216.127.71.83",
		"216.127.71.84",
		"216.127.71.85",
		"216.127.71.86",
		"216.127.71.87",
		"216.127.71.88",
		"216.127.71.89",
		"216.127.71.90",
		"216.127.71.91",
		"216.127.71.92",
		"216.127.71.93",
		"216.127.71.94",
		"216.127.71.95",
		"216.127.71.96",
		"216.127.71.97",
		"216.127.71.98",
		"216.127.71.99",
		"216.127.71.100",
		"216.127.71.101",
		"216.127.71.102",
		"216.127.71.103",
		"216.127.71.104",
		"216.127.71.105",
		"216.127.71.106",
		"216.127.71.107",
		"216.127.71.108",
		"216.127.71.109",
		"216.127.71.110",
		"216.127.71.111",
		"216.127.71.112",
		"216.127.71.113",
		"216.127.71.114",
		"216.127.71.115",
		"216.127.71.116",
		"216.127.71.117",
		"216.127.71.118",
		"216.127.71.119",
		"216.127.71.120",
		"216.127.71.121",
		"216.127.71.122",
		"216.127.71.123",
		"216.127.71.124",
		"216.127.71.125",
		"216.127.71.126",
		"216.127.71.127",
		"216.127.71.128",
		"216.127.71.129",
		"216.127.71.130",
		"216.127.71.131",
		"216.127.71.132",
		"216.127.71.133",
		"216.127.71.134",
		"216.127.71.135",
		"216.127.71.136",
		"216.127.71.137",
		"216.127.71.138",
		"216.127.71.139",
		"216.127.71.140",
		"216.127.71.141",
		"216.127.71.142",
		"216.127.71.143",
		"216.127.71.144",
		"216.127.71.145",
		"216.127.71.146",
		"216.127.71.147",
		"216.127.71.148",
		"216.127.71.149",
		"216.127.71.150",
		"216.127.71.151",
		"216.127.71.152",
		"216.127.71.153",
		"216.127.71.154",
		"216.127.71.155",
		"216.127.71.156",
		"216.127.71.157",
		"216.127.71.158",
		"216.127.71.159",
		"216.127.71.160",
		"216.127.71.161",
		"216.127.71.162",
		"216.127.71.163",
		"216.127.71.164",
		"216.127.71.165",
		"216.127.71.166",
		"216.127.71.167",
		"216.127.71.168",
		"216.127.71.169",
		"216.127.71.170",
		"216.127.71.171",
		"216.127.71.172",
		"216.127.71.173",
		"216.127.71.174",
		"216.127.71.175",
		"216.127.71.176",
		"216.127.71.177",
		"216.127.71.178",
		"216.127.71.179",
		"216.127.71.180",
		"216.127.71.181",
		"216.127.71.182",
		"216.127.71.183",
		"216.127.71.184",
		"216.127.71.185",
		"216.127.71.186",
		"216.127.71.187",
		"216.127.71.188",
		"216.127.71.189",
		"216.127.71.190",
		"216.127.71.191",
		"216.127.71.192",
		"216.127.71.193",
		"216.127.71.194",
		"216.127.71.195",
		"216.127.71.196",
		"216.127.71.197",
		"216.127.71.198",
		"216.127.71.199",
		"216.127.71.200",
		"216.127.71.201",
		"216.127.71.202",
		"216.127.71.203",
		"216.127.71.204",
		"216.127.71.205",
		"216.127.71.206",
		"216.127.71.207",
		"216.127.71.208",
		"216.127.71.209",
		"216.127.71.210",
		"216.127.71.211",
		"216.127.71.212",
		"216.127.71.213",
		"216.127.71.214",
		"216.127.71.215",
		"216.127.71.216",
		"216.127.71.217",
		"216.127.71.218",
		"216.127.71.219",
		"216.127.71.220",
		"216.127.71.221",
		"216.127.71.222",
		"216.127.71.223",
		"216.127.71.224",
		"216.127.71.225",
		"216.127.71.226",
		"216.127.71.227",
		"216.127.71.228",
		"216.127.71.229",
		"216.127.71.230",
		"216.127.71.231",
		"216.127.71.232",
		"216.127.71.233",
		"216.127.71.234",
		"216.127.71.235",
		"216.127.71.236",
		"216.127.71.237",
		"216.127.71.238",
		"216.127.71.239",
		"216.127.71.240",
		"216.127.71.241",
		"216.127.71.242",
		"216.127.71.243",
		"216.127.71.244",
		"216.127.71.245",
		"216.127.71.246",
		"216.127.71.247",
		"216.127.71.248",
		"216.127.71.249",
		"216.127.71.250",
		"216.127.71.251",
		"216.127.71.252",
		"216.127.71.253",
		"216.127.71.254",
		"216.127.71.255"}
	return exists(p.ipAddress, ipWhitelist)
}

func (p *pingback) isSignatureValid() bool {
	signature := ""
	signatureParamsToSign := map[string]string{}
	if val, ok := p.parameters["sig"]; ok {
		signature = val
	}

	signatureParams := []string{}
	if p.apiType == API_VC {
		signatureParams = append(signatureParams, []string{"uid", "currency", "type", "ref"}...)
	} else if p.apiType == API_GOODS {
		signatureParams = append(signatureParams, []string{"uid", "goodsid", "slength", "speriod", "type", "ref"}...)
	} else {
		signatureParams = append(signatureParams, []string{"uid", "goodsid", "type", "ref"}...)
		p.parameters["sign_version"] = string(rune(SIGNATURE_VERSION_2))
	}

	if _, ok := p.parameters["sign_version"] ; !ok { //Check if signature version 1
		for i := range signatureParams {
			if v, o := p.parameters[signatureParams[i]]; o && len(v) > 0 {
				signatureParamsToSign[signatureParams[i]] = p.parameters[signatureParams[i]]
			} else {
				signatureParamsToSign[signatureParams[i]] = ""
			}
		}
		p.parameters["sign_version"] = string(rune(SIGNATURE_VERSION_1))
	} else {
		signatureParamsToSign = p.parameters
	}

	signVersion, err := strconv.Atoi(p.parameters["sign_version"])

	if err != nil {
		Errorsstring = append(Errorsstring, err.Error())
	}
	signatureCalculated := p.calculateSignature(signatureParamsToSign, p.secretKey, signVersion)
	fmt.Println(signatureCalculated,signature)
	return signatureCalculated == signature
}

func (p *pingback) getPingbackType() (int, error) { //changed to getPingbackType() to avoid duplicate name with C# method getType()
	if p.parameters["type"] != "" {
		t, err := strconv.Atoi(p.parameters["type"])
		if err != nil {
			return -1, err
		}
		return t, nil
	} else {
		return -1, nil
	}
}

func (p *pingback) calculateSignature(signatureParamsToSign map[string]string, secret string, version int) string {

	baseString := ""
	delete(signatureParamsToSign, "sig")
	if version != SIGNATURE_VERSION_1 {
		keys := make([]string, 0)
		for k, _ := range signatureParamsToSign {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, v := range keys {
			baseString += fmt.Sprintf(`%s=%s`, v, signatureParamsToSign[v])
		}
	} else {
		for k, v := range signatureParamsToSign {
			baseString += fmt.Sprintf(`%s=%s`, k, v)
		}
	}
	baseString += secret

	if version == SIGNATURE_VERSION_3 {
		return gethash(baseString, "sha256")
	}
	return gethash(baseString, "md5")

}

func (p *pingback) IsDeliverable() (bool, error) {
	t, err := p.getPingbackType()
	if err != nil {
		return false, err
	}
	return t == PINGBACK_TYPE_REGULAR || t == PINGBACK_TYPE_GOODWILL || t == PINGBACK_TYPE_RISK_REVIEWED_ACCEPTED, nil
}

func (p *pingback) IsCancelable() (bool, error) {

	t, err := p.getPingbackType()
	if err != nil {
		return false, err
	}
	return t == PINGBACK_TYPE_REGULAR || t == PINGBACK_TYPE_NEGATIVE || t == PINGBACK_TYPE_RISK_REVIEWED_DECLINED, nil
}

func (p *pingback) IsUnderReview() (bool, error) {
	t, err := p.getPingbackType()
	if err != nil {
		return false, err
	}
	return t == PINGBACK_TYPE_RISK_UNDER_REVIEW, nil
}

func (p *pingback) GetParameter(param string) string {

	if len(p.parameters[param]) > 0 {
		return p.parameters[param]
	}

	return ""
}

/**
 * Get pingback parameter "type"
 *
 * @return int
 */
func (p *pingback) GetPingbackType() (int, error) {
	val, ok := p.parameters["type"]
	if ok {
		t, err := strconv.Atoi(val)
		if err != nil {
			return -1, err
		}
		return t, nil
	}
	return -1, errors.New("get error")
}

/**
 * Get verbal explanation of the informational pingback
 *
 * @return string
 */
func (p *pingback) getTypeVerbal() string {
	pingbackTypes := map[string]string{}
	pingbackTypes[fmt.Sprintf("%d", PINGBACK_TYPE_SUBSCRIPTION_CANCELLATION)] = "user_subscription_cancellation"
	pingbackTypes[fmt.Sprintf("%d", PINGBACK_TYPE_SUBSCRIPTION_EXPIRED)] = "user_subscription_expired"
	pingbackTypes[fmt.Sprintf("%d", PINGBACK_TYPE_SUBSCRIPTION_PAYMENT_FAILED)] = "user_subscription_payment_failed"

	if len(p.parameters["type"]) > 0 {
		if val, ok := pingbackTypes[p.parameters["type"]]; ok {
			return pingbackTypes[val]
		}
		return ""

	}
	return ""
}

/**
 * Get pingback parameter "uid"
 *
 * @return string
 */
func (p *pingback) GetUserId() string {

	return p.GetParameter("uid")
}

/**
 * Get pingback parameter "currency"
 *
 * @return string
 */
func (p *pingback) GetVirtualCurrencyAmount() string {
	return p.GetParameter("currency")
}

/**
 * Get product id
 *
 * @return string
 */
func (p *pingback) GetProductId() string {
	return p.GetParameter("goodsid")
}

/**
 * @return int
 */
func (p *pingback) GetProductPeriodLength() (int, error) {

	val := p.GetParameter("slength")
	t, err := strconv.Atoi(val)
	return t, err
}

/*
 * @return string
 */
func (p *pingback) GetProductPeriodType() string {
	return p.GetParameter("speriod")
}

/*
 *  @return Paymentwall_Product
 */
func (p *pingback) GetProduct() (Product, error) {
	productType := ""
	productPeriodLength, err := p.GetProductPeriodLength()
	if err != nil {
		return Product{}, err
	}
	if productPeriodLength > 0 {
		productType = TYPE_SUBSCRIPTION
	} else {
		productType = TYPE_SUBSCRIPTION
	}

	return Product{
		productId:    p.GetProductId(),
		amount:       0,
		currencyCode: "",
		name:         "",
		productType:  productType,
		periodLength: productPeriodLength,
		periodType:   p.GetProductPeriodType(),
	}, nil

}

/*
 * @return List<Paymentwall_Product>
 */
func (p *pingback) GetProducts() []Product {
	products := []Product{}
	productIds := []string{}

	for _, v := range p.parameters["goodsid"] {
		productIds = append(productIds, fmt.Sprintf("%d", v))
	}

	if len(productIds) > 0 {
		for _, v := range productIds {
			products = append(products, Product{productId: v})
		}
	}

	return products
}

/*
 * Get pingback parameter "ref"
 *
 * @return string
 */
func (p *pingback) GetReferenceId() string {

	return p.GetParameter("ref")
}

/*
 * Returns unique identifier of the pingback that can be used for checking
 * If the same pingback was already processed by your servers
 * Two pingbacks with the same unique ID should not be processed more than once
 *
 * @return string
 */
func (p *pingback) GetPingbackUniqueId() (string, error) {

	pingBackType, err := p.GetPingbackType()
	if err != nil {
		return "", err
	}
	return p.GetReferenceId() + "_" + fmt.Sprintf("%d", pingBackType), nil
}
