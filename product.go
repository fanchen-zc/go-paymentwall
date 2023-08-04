package paymentwall

const (
	/**
	 * Product types
	 */
	TYPE_SUBSCRIPTION = "subscription"
	TYPE_FIXED        = "fixed"

	/**
	 * Product period types
	 */
	PERIOD_TYPE_DAY   = "day"
	PERIOD_TYPE_WEEK  = "week"
	PERIOD_TYPE_MONTH = "month"
	PERIOD_TYPE_YEAR  = "year"
)

/**
 * Paymentwall_Product class's properties
 */

type Product struct {
	productId    string
	amount       float64
	currencyCode string
	name         string
	productType  string
	periodLength int
	periodType   string
	recurring    bool
	trialProduct interface{}
}

type Productfun func(*Product)

func WithAmount(amount float64) Productfun {
	return func(p *Product) {
		p.amount = amount
	}
}

func WithCurrencyCode(currencyCode string) Productfun {
	return func(p *Product) {
		p.currencyCode = currencyCode
	}
}

func WithProductType(productType string) Productfun {
	return func(p *Product) {
		p.productType = productType
	}
}

func WithPeriodLength(periodLength int) Productfun {
	return func(p *Product) {
		p.periodLength = periodLength
	}
}

func WithPeriodType(periodType string) Productfun {
	return func(p *Product) {
		p.productType = periodType
	}
}

func WithRecurring(recurring bool) Productfun {
	return func(p *Product) {
		p.recurring = recurring
	}
}

func WithTrialProduct(trialProduct Product) Productfun {
	return func(p *Product) {
		p.trialProduct = trialProduct
	}
}

func NewProduct(productId, name, currencyCode string, amount float64, params ...Productfun) *Product {
	p := &Product{
		productId: productId,
		name: name,
		currencyCode: currencyCode,
		amount: amount,
		productType: TYPE_FIXED,
		periodLength: 0,
		periodType: "",
		recurring: false,
	}

	for _, paramfun := range params {
		paramfun(p)
	}
	if p.periodType != TYPE_SUBSCRIPTION && !p.recurring {
		p.trialProduct = nil
	}
	return p
}
