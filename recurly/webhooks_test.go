package recurly

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestWebhooks(t *testing.T) {
	t.SkipNow()
	suite := []map[string]interface{}{
		map[string]interface{}{
			"name": "new_account_notification",
			"body": `<?xml version="1.0" encoding="UTF-8"?>
                <new_account_notification>
                    <account>
                        <account_code>1</account_code>
                        <username nil="true"></username>
                        <email>verena@example.com</email>
                        <first_name>Verena</first_name>
                        <last_name>Example</last_name>
                        <company_name nil="true"></company_name>
                        </account>
                 </new_account_notification>`,
			"expected": AccountNotification{
				webhook: webhook{
					XMLName: xml.Name{Local: "new_account_notification"},
				},
				Account: Account{
					Code:      "1",
					Email:     "verena@example.com",
					FirstName: "Verena",
					LastName:  "Example",
				},
			},
		},
		map[string]interface{}{
			"name": "new_invoice_notification",
			"body": `<?xml version="1.0" encoding="UTF-8"?>
                <new_invoice_notification>
                  <account>
                    <account_code>1</account_code>
                    <username nil="true"></username>
                    <email>verena@example.com</email>
                    <first_name>Verana</first_name>
                    <last_name>Example</last_name>
                    <company_name nil="true"></company_name>
                  </account>
                  <invoice>
                    <uuid>ffc64d71d4b5404e93f13aac9c63b007</uuid>
                    <subscription_id nil="true"></subscription_id>
                    <state>open</state>
                    <invoice_number_prefix></invoice_number_prefix>
                    <invoice_number type="integer">1000</invoice_number>
                    <po_number></po_number>
                    <vat_number></vat_number>
                    <total_in_cents type="integer">1000</total_in_cents>
                    <currency>USD</currency>
                    <date type="datetime">2014-01-01T20:21:44Z</date>
                    <closed_at type="datetime" nil="true"></closed_at>
                    <net_terms type="integer">0</net_terms>
                    <collection_method>manual</collection_method>
                  </invoice>
                </new_invoice_notification>`,
			"expected": InvoiceNotification{},
		},
	}

	for i, s := range suite {
		body := bytes.NewBufferString(s["body"].(string))
		name := s["name"].(string)

		notification, err := HandleWebhook(body)
		if err != nil {
			t.Errorf("TestWebhooks Error (%d): Error from HandleWebhook. Err: %s", i, err)
		}

		if name != notification.Type() {
			t.Errorf("TestWebhooks Error (%d): Expected %s, given %s", i, name, notification.Type())
		}

		if !reflect.DeepEqual(s["expected"], notification) {
			t.Errorf("TestWebhooks Error (%d): Expected %s to unmarshal to struct of %+v, given %+v", i, name, s["expected"], notification)
		}
	}
}
