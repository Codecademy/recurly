package recurly

import (
	"encoding/xml"
	"fmt"
)

type (
	// RedemptionsService handles communication with the coupon redemption
	// related methods of the recurly API.
	RedemptionsService struct {
		client *Client
	}

	// Redemption holds redeemed coupons for an account or invoice.
	Redemption struct {
		CouponCode             string
		AccountCode            string
		SingleUse              NullBool
		TotalDiscountedInCents int
		Currency               string
		State                  string
		CreatedAt              NullTime
	}
)

// UnmarshalXML unmarshal a coupon redemption object. Minaly converts href links
// for coupons and accounts to CouponCode and AccountCodes.
func (r *Redemption) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName                xml.Name   `xml:"redemption"`
		CouponCode             hrefString `xml:"coupon,omitempty"`
		AccountCode            hrefString `xml:"account,omitempty"`
		SingleUse              NullBool   `xml:"single_use,omitempty"`
		TotalDiscountedInCents int        `xml:"total_discounted_in_cents,omitempty"`
		Currency               string     `xml:"currency,omitempty"`
		State                  string     `xml:"state,omitempty"`
		CreatedAt              NullTime   `xml:"created_at,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*r = Redemption{
		CouponCode:             string(v.CouponCode),
		AccountCode:            string(v.AccountCode),
		SingleUse:              v.SingleUse,
		TotalDiscountedInCents: v.TotalDiscountedInCents,
		Currency:               v.Currency,
		State:                  v.State,
		CreatedAt:              v.CreatedAt,
	}

	return nil
}

// GetForAccount looks up information about the 'active' coupon redemption on
// an account
// https://dev.recurly.com/docs/lookup-a-coupon-redemption-on-an-account
func (s *RedemptionsService) GetForAccount(accountCode string) (*Response, *Redemption, error) {
	action := fmt.Sprintf("accounts/%s/redemption", accountCode)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Redemption
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// GetForInvoice looks up information about a coupon redemption applied
// to an invoice.
// https://dev.recurly.com/docs/lookup-a-coupon-redemption-on-an-invoice
func (s *RedemptionsService) GetForInvoice(invoiceNumber string) (*Response, *Redemption, error) {
	action := fmt.Sprintf("invoices/%s/redemption", invoiceNumber)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Redemption
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Redeem will redeem a coupon before or after a subscription. Most coupons are
// redeemed during a new subscription. This endpoint allows you to redeem a
// coupon for a customer after their initial subscription, or in anticipation
// of a future subscription. When you redeem a coupon on an account, the coupon
// will be applied to the next subscription creation (new subscription),
// modification (e.g. upgrade or downgrade), or renewal.
// https://dev.recurly.com/docs/redeem-a-coupon-before-or-after-a-subscription
func (s *RedemptionsService) Redeem(code string, accountCode string, currency string) (*Response, *Redemption, error) {
	action := fmt.Sprintf("coupons/%s/redeem", code)
	data := struct {
		XMLName     xml.Name `xml:"redemption"`
		AccountCode string   `xml:"account_code"`
		Currency    string   `xml:"currency"`
	}{
		AccountCode: accountCode,
		Currency:    currency,
	}
	req, err := s.client.newRequest("POST", action, nil, data)
	if err != nil {
		return nil, nil, err
	}

	var dst Redemption
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Delete removes a coupon from an account. Recurly will automatically remove
// coupons after they expire or are otherwise no longer valid for an account.
// If you want to remove a coupon from an account before it expires, use this
// function. Please note: the coupon will still count towards the
// "maximum redemption total" of a coupon.
// https://dev.recurly.com/docs/remove-a-coupon-from-an-account
func (s *RedemptionsService) Delete(accountCode string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s/redemption", accountCode)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}
