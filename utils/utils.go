package utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Define structs to match the target XML structure
type Invoice struct {
	XMLName                       xml.Name                         `xml:"invoice"`
	InvoiceNumber                 *Identifier                      `xml:"Invoice_number,omitempty"`
	InvoiceIssueDate              *Date                            `xml:"Invoice_issue_date,omitempty"`
	InvoiceTypeCode               *Code                            `xml:"Invoice_type_code,omitempty"`
	InvoiceCurrencyCode           *Code                            `xml:"Invoice_currency_code,omitempty"`
	VATAccountingCurrencyCode     *Code                            `xml:"VAT_accounting_currency_code,omitempty"`
	ValueAddedTaxPointDate        *ValueAddedTaxPointDate          `xml:"Value_added_tax_point_date,omitempty"`
	ValueAddedTaxPointDateCode    *Code                            `xml:"Value_added_tax_point_date_code,omitempty"`
	PaymentDueDate                *PaymentDueDate                  `xml:"Payment_due_date,omitempty"`
	BuyerReference                *Text                            `xml:"Buyer_reference,omitempty"`
	ProjectReference              *DocumentReference               `xml:"Project_reference,omitempty"`
	ContractReference             *DocumentReference               `xml:"Contract_reference,omitempty"`
	PurchaseOrderReference        *DocumentReference               `xml:"Purchase_order_reference,omitempty"`
	SalesOrderReference           *DocumentReference               `xml:"Sales_order_reference,omitempty"`
	ReceivingAdviceReference      *DocumentReference               `xml:"Receiving_advice_reference,omitempty"`
	DespatchAdviceReference       *DocumentReference               `xml:"Despatch_advice_reference,omitempty"`
	TenderOrLotReference          *DocumentReference               `xml:"Tender_or_lot_reference,omitempty"`
	InvoicedObjectIdentifier      *IdentifierWithScheme            `xml:"Invoiced_object_identifier,omitempty"`
	BuyerAccountingReference      *Text                            `xml:"Buyer_accounting_reference,omitempty"`
	PaymentTerms                  *PaymentTerms                    `xml:"Payment_terms,omitempty"`
	InvoiceNote                   []*InvoiceNote                   `xml:"INVOICE_NOTE,omitempty"`
	ProcessControl                *ProcessControl                  `xml:"PROCESS_CONTROL,omitempty"`
	PrecedingInvoiceReference     *PrecedingInvoiceReference       `xml:"PRECEDING_INVOICE_REFERENCE,omitempty"`
	Seller                        *Party                           `xml:"SELLER,omitempty"`
	Buyer                         *Party                           `xml:"BUYER,omitempty"`
	Payee                         *Party                           `xml:"PAYEE,omitempty"`
	SellerTaxRepresentativeParty  *TaxRepresentativeParty          `xml:"SELLER_TAX_REPRESENTATIVE_PARTY,omitempty"`
	DeliveryInformation           *DeliveryInformation             `xml:"DELIVERY_INFORMATION,omitempty"`
	InvoicingPeriod               *InvoicingPeriod                 `xml:"INVOICING_PERIOD,omitempty"`
	PaymentInstructions           []*PaymentInstructions           `xml:"PAYMENT_INSTRUCTIONS,omitempty"`
	DocumentLevelAllowances       []*DocumentLevelAllowances       `xml:"DOCUMENT_LEVEL_ALLOWANCES,omitempty"`
	DocumentLevelCharges          []*DocumentLevelCharges          `xml:"DOCUMENT_LEVEL_CHARGES,omitempty"`
	DocumentTotals                *DocumentTotals                  `xml:"DOCUMENT_TOTALS,omitempty"`
	VATBreakdown                  []*VATBreakdown                  `xml:"VAT_BREAKDOWN,omitempty"`
	AdditionalSupportingDocuments []*AdditionalSupportingDocuments `xml:"ADDITIONAL_SUPPORTING_DOCUMENTS,omitempty"`
	InvoiceLine                   []*InvoiceLine                   `xml:"INVOICE_LINE,omitempty"`
}

// Common basic types
type Identifier struct {
	XMLName           xml.Name `xml:",omitempty"`
	Scheme_identifier string   `xml:"scheme_identifier,attr,omitempty"`
	Text              string   `xml:",chardata"`
}
type IdentifierWithScheme struct {
	XMLName                   xml.Name `xml:",omitempty"`
	Scheme_identifier         string   `xml:"scheme_identifier,attr,omitempty"`
	Scheme_version_identifier string   `xml:"scheme_version_identifier,attr,omitempty"`
	Text                      string   `xml:",chardata"`
}
type Code struct {
	XMLName xml.Name `xml:",omitempty"`
	Text    string   `xml:",chardata"`
}
type Date struct {
	XMLName xml.Name `xml:",omitempty"`
	Text    string   `xml:",chardata"`
}

type Text struct {
	XMLName xml.Name `xml:",omitempty"`
	Text    string   `xml:",chardata"`
}
type DocumentReference struct {
	XMLName xml.Name `xml:",omitempty"`
	Text    string   `xml:",chardata"`
}
type BinaryObject struct {
	XMLName   xml.Name `xml:",omitempty"`
	Mime_code string   `xml:"mime_code,attr,omitempty"`
	Filename  string   `xml:"filename,attr,omitempty"`
	Text      string   `xml:",chardata"`
}

type PaymentDueDate struct {
	XMLName xml.Name `xml:",omitempty"`
	Id      string   `xml:"xr:id,attr"`
	Src     string   `xml:"xr:src,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type ValueAddedTaxPointDate struct {
	XMLName xml.Name `xml:",omitempty"`
	Id      string   `xml:"xr:id,attr"`
	Src     string   `xml:"xr:src,attr,omitempty"`
	Text    string   `xml:",chardata"`
}
type PaymentTerms struct {
	XMLName xml.Name `xml:",omitempty"`
	Id      string   `xml:"xr:id,attr"`
	Src     string   `xml:"xr:src,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

// BG-1: INVOICE_NOTE
type InvoiceNote struct {
	XMLName                xml.Name `xml:"INVOICE_NOTE"`
	Id                     string   `xml:"xr:id,attr"`
	Src                    string   `xml:"xr:src,attr,omitempty"`
	InvoiceNoteSubjectCode *Code    `xml:"Invoice_note_subject_code,omitempty"`
	InvoiceNote            *Text    `xml:"Invoice_note,omitempty"`
}

// BG-2: PROCESS_CONTROL
type ProcessControl struct {
	XMLName                 xml.Name    `xml:"PROCESS_CONTROL"`
	Id                      string      `xml:"xr:id,attr"`
	Src                     string      `xml:"xr:src,attr,omitempty"`
	BusinessProcessType     *Text       `xml:"Business_process_type,omitempty"`
	SpecificationIdentifier *Identifier `xml:"Specification_identifier,omitempty"`
}

// BG-3: PRECEDING_INVOICE_REFERENCE
type PrecedingInvoiceReference struct {
	XMLName                   xml.Name           `xml:"PRECEDING_INVOICE_REFERENCE"`
	Id                        string             `xml:"xr:id,attr"`
	Src                       string             `xml:"xr:src,attr,omitempty"`
	PrecedingInvoiceReference *DocumentReference `xml:"Preceding_Invoice_reference,omitempty"`
	PrecedingInvoiceIssueDate *Date              `xml:"Preceding_Invoice_issue_date,omitempty"`
}

// BG-4: SELLER
type Party struct {
	XMLName                           xml.Name              `xml:",omitempty"`
	SellerTradingName                 *Text                 `xml:"Seller_trading_name,omitempty"`
	SellerIdentifier                  *IdentifierWithScheme `xml:"Seller_identifier,omitempty"`
	SellerLegalRegistrationIdentifier *IdentifierWithScheme `xml:"Seller_legal_registration_identifier,omitempty"`
	SellerVATIdentifier               *Identifier           `xml:"Seller_VAT_identifier,omitempty"`
	SellerTaxRegistrationIdentifier   *Identifier           `xml:"Seller_tax_registration_identifier,omitempty"`
	SellerAdditionalLegalInformation  *Text                 `xml:"Seller_additional_legal_information,omitempty"`
	SellerElectronicAddress           *IdentifierWithScheme `xml:"Seller_electronic_address,omitempty"`
	SellerPostalAddress               *PostalAddress        `xml:"SELLER_POSTAL_ADDRESS,omitempty"`
	SellerContact                     *Contact              `xml:"SELLER_CONTACT,omitempty"`
	BuyerName                         *Text                 `xml:"Buyer_name,omitempty"`
	BuyerTradingName                  *Text                 `xml:"Buyer_trading_name,omitempty"`
	BuyerIdentifier                   *IdentifierWithScheme `xml:"Buyer_identifier,omitempty"`
	BuyerLegalRegistrationIdentifier  *IdentifierWithScheme `xml:"Buyer_legal_registration_identifier,omitempty"`
	BuyerVATIdentifier                *Identifier           `xml:"Buyer_VAT_identifier,omitempty"`
	BuyerElectronicAddress            *IdentifierWithScheme `xml:"Buyer_electronic_address,omitempty"`
	BuyerPostalAddress                *PostalAddress        `xml:"BUYER_POSTAL_ADDRESS,omitempty"`
	BuyerContact                      *Contact              `xml:"BUYER_CONTACT,omitempty"`
	PayeeName                         *Text                 `xml:"Payee_name,omitempty"`
	PayeeIdentifier                   *IdentifierWithScheme `xml:"Payee_identifier,omitempty"`
	PayeeLegalRegistrationIdentifier  *IdentifierWithScheme `xml:"Payee_legal_registration_identifier,omitempty"`
}
type TaxRepresentativeParty struct {
	XMLName                              xml.Name       `xml:",omitempty"`
	Id                                   string         `xml:"xr:id,attr"`
	Src                                  string         `xml:"xr:src,attr,omitempty"`
	SellerTaxRepresentativeName          *Text          `xml:"Seller_tax_representative_name,omitempty"`
	SellerTaxRepresentativeVATIdentifier *Identifier    `xml:"Seller_tax_representative_VAT_identifier,omitempty"`
	SellerTaxRepresentativePostalAddress *PostalAddress `xml:"SELLER_TAX_REPRESENTATIVE_POSTAL_ADDRESS,omitempty"`
}

// BG-5, BG-8, BG-12, BG-15: POSTAL_ADDRESS
type PostalAddress struct {
	XMLName                             xml.Name `xml:",omitempty"`
	Id                                  string   `xml:"xr:id,attr"`
	Src                                 string   `xml:"xr:src,attr,omitempty"`
	SellerAddressLine1                  *Text    `xml:"Seller_address_line_1,omitempty"`
	SellerAddressLine2                  *Text    `xml:"Seller_address_line_2,omitempty"`
	SellerAddressLine3                  *Text    `xml:"Seller_address_line_3,omitempty"`
	SellerCity                          *Text    `xml:"Seller_city,omitempty"`
	SellerPostCode                      *Text    `xml:"Seller_post_code,omitempty"`
	SellerCountrySubdivision            *Text    `xml:"Seller_country_subdivision,omitempty"`
	SellerCountryCode                   *Code    `xml:"Seller_country_code,omitempty"`
	BuyerAddressLine1                   *Text    `xml:"Buyer_address_line_1,omitempty"`
	BuyerAddressLine2                   *Text    `xml:"Buyer_address_line_2,omitempty"`
	BuyerAddressLine3                   *Text    `xml:"Buyer_address_line_3,omitempty"`
	BuyerCity                           *Text    `xml:"Buyer_city,omitempty"`
	BuyerPostCode                       *Text    `xml:"Buyer_post_code,omitempty"`
	BuyerCountrySubdivision             *Text    `xml:"Buyer_country_subdivision,omitempty"`
	BuyerCountryCode                    *Code    `xml:"Buyer_country_code,omitempty"`
	TaxRepresentativeAddressLine1       *Text    `xml:"Tax_representative_address_line_1,omitempty"`
	TaxRepresentativeAddressLine2       *Text    `xml:"Tax_representative_address_line_2,omitempty"`
	TaxRepresentativeAddressLine3       *Text    `xml:"Tax_representative_address_line_3,omitempty"`
	TaxRepresentativeCity               *Text    `xml:"Tax_representative_city,omitempty"`
	TaxRepresentativePostCode           *Text    `xml:"Tax_representative_post_code,omitempty"`
	TaxRepresentativeCountrySubdivision *Text    `xml:"Tax_representative_country_subdivision,omitempty"`
	TaxRepresentativeCountryCode        *Code    `xml:"Tax_representative_country_code,omitempty"`
	DeliverToAddressLine1               *Text    `xml:"Deliver_to_address_line_1,omitempty"`
	DeliverToAddressLine2               *Text    `xml:"Deliver_to_address_line_2,omitempty"`
	DeliverToAddressLine3               *Text    `xml:"Deliver_to_address_line_3,omitempty"`
	DeliverToCity                       *Text    `xml:"Deliver_to_city,omitempty"`
	DeliverToPostCode                   *Text    `xml:"Deliver_to_post_code,omitempty"`
	DeliverToCountrySubdivision         *Text    `xml:"Deliver_to_country_subdivision,omitempty"`
	DeliverToCountryCode                *Code    `xml:"Deliver_to_country_code,omitempty"`
}

// BG-6, BG-9: CONTACT
type Contact struct {
	XMLName                      xml.Name `xml:",omitempty"`
	Id                           string   `xml:"xr:id,attr"`
	Src                          string   `xml:"xr:src,attr,omitempty"`
	SellerContactPoint           *Text    `xml:"Seller_contact_point,omitempty"`
	SellerContactTelephoneNumber *Text    `xml:"Seller_contact_telephone_number,omitempty"`
	SellerContactEmailAddress    *Text    `xml:"Seller_contact_email_address,omitempty"`
	BuyerContactPoint            *Text    `xml:"Buyer_contact_point,omitempty"`
	BuyerContactTelephoneNumber  *Text    `xml:"Buyer_contact_telephone_number,omitempty"`
	BuyerContactEmailAddress     *Text    `xml:"Buyer_contact_email_address,omitempty"`
}

// BG-13: DELIVERY_INFORMATION
type DeliveryInformation struct {
	XMLName                     xml.Name              `xml:"DELIVERY_INFORMATION"`
	Id                          string                `xml:"xr:id,attr"`
	Src                         string                `xml:"xr:src,attr,omitempty"`
	DeliverToPartyName          *Text                 `xml:"Deliver_to_party_name,omitempty"`
	DeliverToLocationIdentifier *IdentifierWithScheme `xml:"Deliver_to_location_identifier,omitempty"`
	ActualDeliveryDate          *Date                 `xml:"Actual_delivery_date,omitempty"`
	DeliverToAddress            *PostalAddress        `xml:"DELIVER_TO_ADDRESS,omitempty"`
}

// BG-14: INVOICING_PERIOD
type InvoicingPeriod struct {
	XMLName                  xml.Name `xml:"INVOICING_PERIOD"`
	Id                       string   `xml:"xr:id,attr"`
	Src                      string   `xml:"xr:src,attr,omitempty"`
	InvoicingPeriodStartDate *Date    `xml:"Invoicing_period_start_date,omitempty"`
	InvoicingPeriodEndDate   *Date    `xml:"Invoicing_period_end_date,omitempty"`
}

// BG-16: PAYMENT_INSTRUCTIONS
type PaymentInstructions struct {
	XMLName                xml.Name                `xml:"PAYMENT_INSTRUCTIONS"`
	Id                     string                  `xml:"xr:id,attr"`
	Src                    string                  `xml:"xr:src,attr,omitempty"`
	PaymentMeansTypeCode   *Code                   `xml:"Payment_means_type_code,omitempty"`
	PaymentMeansText       *Text                   `xml:"Payment_means_text,omitempty"`
	RemittanceInformation  *Text                   `xml:"Remittance_information,omitempty"`
	CreditTransfer         *CreditTransfer         `xml:"CREDIT_TRANSFER,omitempty"`
	PaymentCardInformation *PaymentCardInformation `xml:"PAYMENT_CARD_INFORMATION,omitempty"`
	DirectDebit            *DirectDebit            `xml:"DIRECT_DEBIT,omitempty"`
}

// BG-17: CREDIT_TRANSFER
type CreditTransfer struct {
	XMLName                          xml.Name    `xml:"CREDIT_TRANSFER"`
	Id                               string      `xml:"xr:id,attr,omitempty"`
	Src                              string      `xml:"xr:src,attr,omitempty"`
	PaymentAccountIdentifier         *Identifier `xml:"Payment_account_identifier,omitempty"`
	PaymentAccountName               *Text       `xml:"Payment_account_name,omitempty"`
	PaymentServiceProviderIdentifier *Text       `xml:"Payment_service_provider_identifier,omitempty"`
}

// BG-18: PAYMENT_CARD_INFORMATION
type PaymentCardInformation struct {
	XMLName                         xml.Name `xml:"PAYMENT_CARD_INFORMATION"`
	Id                              string   `xml:"xr:id,attr"`
	Src                             string   `xml:"xr:src,attr,omitempty"`
	PaymentCardPrimaryAccountNumber *Text    `xml:"Payment_card_primary_account_number,omitempty"`
	PaymentCardHolderName           *Text    `xml:"Payment_card_holder_name,omitempty"`
}

// BG-19: DIRECT_DEBIT
type DirectDebit struct {
	XMLName                        xml.Name    `xml:"DIRECT_DEBIT"`
	Id                             string      `xml:"xr:id,attr"`
	Src                            string      `xml:"xr:src,attr,omitempty"`
	MandateReferenceIdentifier     *Identifier `xml:"Mandate_reference_identifier,omitempty"`
	BankAssignedCreditorIdentifier *Identifier `xml:"Bank_assigned_creditor_identifier,omitempty"`
	DebitedAccountIdentifier       *Identifier `xml:"Debited_account_identifier,omitempty"`
}

// BG-20: DOCUMENT_LEVEL_ALLOWANCES
type DocumentLevelAllowances struct {
	XMLName                          xml.Name `xml:"DOCUMENT_LEVEL_ALLOWANCES"`
	Id                               string   `xml:"xr:id,attr"`
	Src                              string   `xml:"xr:src,attr,omitempty"`
	DocumentLevelAllowanceAmount     *Text    `xml:"Document_level_allowance_amount,omitempty"`
	DocumentLevelAllowanceBaseAmount *Text    `xml:"Document_level_allowance_base_amount,omitempty"`
	DocumentLevelAllowancePercentage *Text    `xml:"Document_level_allowance_percentage,omitempty"`
	DocumentLevelVATCategoryCode     *Code    `xml:"Document_level_allowance_VAT_category_code,omitempty"`
	DocumentLevelVATRate             *Text    `xml:"Document_level_allowance_VAT_rate,omitempty"`
	DocumentLevelAllowanceReason     *Text    `xml:"Document_level_allowance_reason,omitempty"`
	DocumentLevelAllowanceReasonCode *Code    `xml:"Document_level_allowance_reason_code,omitempty"`
}

// BG-21: DOCUMENT_LEVEL_CHARGES
type DocumentLevelCharges struct {
	XMLName                       xml.Name `xml:"DOCUMENT_LEVEL_CHARGES"`
	Id                            string   `xml:"xr:id,attr"`
	Src                           string   `xml:"xr:src,attr,omitempty"`
	DocumentLevelChargeAmount     *Text    `xml:"Document_level_charge_amount,omitempty"`
	DocumentLevelChargeBaseAmount *Text    `xml:"Document_level_charge_base_amount,omitempty"`
	DocumentLevelChargePercentage *Text    `xml:"Document_level_charge_percentage,omitempty"`
	DocumentLevelVATCategoryCode  *Code    `xml:"Document_level_charge_VAT_category_code,omitempty"`
	DocumentLevelVATRate          *Text    `xml:"Document_level_charge_VAT_rate,omitempty"`
	DocumentLevelChargeReason     *Text    `xml:"Document_level_charge_reason,omitempty"`
	DocumentLevelChargeReasonCode *Code    `xml:"Document_level_charge_reason_code,omitempty"`
}

// BG-22: DOCUMENT_TOTALS
type DocumentTotals struct {
	XMLName xml.Name `xml:"DOCUMENT_TOTALS"`
	Id      string   `xml:"xr:id,attr"`
	Src     string   `xml:"xr:src,attr,omitempty"`
}

type InvoiceLineNetAmount struct {
	SumOfAllowancesOnDocumentLevel            *Text `xml:"Sum_of_allowances_on_document_level,omitempty"`
	SumOfChargesOnDocumentLevel               *Text `xml:"Sum_of_charges_on_document_level,omitempty"`
	InvoiceTotalAmountWithoutVAT              *Text `xml:"Invoice_total_amount_without_VAT,omitempty"`
	InvoiceTotalVATAmount                     *Text `xml:"Invoice_total_VAT_amount,omitempty"`
	InvoiceTotalVATAmountInAccountingCurrency *Text `xml:"Invoice_total_VAT_amount_in_accounting_currency,omitempty"`
	InvoiceTotalAmountWithVAT                 *Text `xml:"Invoice_total_amount_with_VAT,omitempty"`
	PaidAmount                                *Text `xml:"Paid_amount,omitempty"`
	RoundingAmount                            *Text `xml:"Rounding_amount,omitempty"`
	AmountDueForPayment                       *Text `xml:"Amount_due_for_payment,omitempty"`
}

// BG-23: VAT_BREAKDOWN
type VATBreakdown struct {
	XMLName                  xml.Name `xml:"VAT_BREAKDOWN"`
	Id                       string   `xml:"xr:id,attr"`
	Src                      string   `xml:"xr:src,attr,omitempty"`
	VATCategoryTaxableAmount *Text    `xml:"VAT_category_taxable_amount,omitempty"`
	VATCategoryTaxAmount     *Text    `xml:"VAT_category_tax_amount,omitempty"`
	VATCategoryCode          *Code    `xml:"VAT_category_code,omitempty"`
	VATCategoryRate          *Text    `xml:"VAT_category_rate,omitempty"`
	VATExemptionReasonText   *Text    `xml:"VAT_exemption_reason_text,omitempty"`
	VATExemptionReasonCode   *Code    `xml:"VAT_exemption_reason_code,omitempty"`
}

// BG-24: ADDITIONAL_SUPPORTING_DOCUMENTS
type AdditionalSupportingDocuments struct {
	XMLName                       xml.Name           `xml:"ADDITIONAL_SUPPORTING_DOCUMENTS"`
	Id                            string             `xml:"xr:id,attr"`
	Src                           string             `xml:"xr:src,attr,omitempty"`
	SupportingDocumentReference   *DocumentReference `xml:"Supporting_document_reference,omitempty"`
	SupportingDocumentDescription *Text              `xml:"Supporting_document_description,omitempty"`
	ExternalDocumentLocation      *Text              `xml:"External_document_location,omitempty"`
	AttachedDocument              *BinaryObject      `xml:"Attached_document,omitempty"`
}

// BG-25: INVOICE_LINE
type InvoiceLine struct {
	XMLName                              xml.Name                 `xml:"INVOICE_LINE"`
	Id                                   string                   `xml:"xr:id,attr"`
	Src                                  string                   `xml:"xr:src,attr,omitempty"`
	InvoiceLineIdentifier                *Identifier              `xml:"Invoice_line_identifier,omitempty"`
	InvoiceLineNote                      *Text                    `xml:"Invoice_line_note,omitempty"`
	InvoiceLineObjectIdentifier          *IdentifierWithScheme    `xml:"Invoice_line_object_identifier,omitempty"`
	InvoicedQuantity                     *Text                    `xml:"Invoiced_quantity,omitempty"`
	InvoicedQuantityUnitOfMeasureCode    *Code                    `xml:"Invoiced_quantity_unit_of_measure_code,omitempty"`
	InvoiceLineNetAmount                 *Text                    `xml:"Invoice_line_net_amount,omitempty"`
	ReferencedPurchaseOrderLineReference *DocumentReference       `xml:"Referenced_purchase_order_line_reference,omitempty"`
	InvoiceLineBuyerAccountingReference  *Text                    `xml:"Invoice_line_Buyer_accounting_reference,omitempty"`
	InvoiceLinePeriod                    *InvoiceLinePeriod       `xml:"INVOICE_LINE_PERIOD,omitempty"`
	InvoiceLineAllowances                []*InvoiceLineAllowances `xml:"INVOICE_LINE_ALLOWANCES,omitempty"`
	InvoiceLineCharges                   []*InvoiceLineCharges    `xml:"INVOICE_LINE_CHARGES,omitempty"`
	PriceDetails                         *PriceDetails            `xml:"PRICE_DETAILS,omitempty"`
	LineVATInformation                   *LineVATInformation      `xml:"LINE_VAT_INFORMATION,omitempty"`
	ItemInformation                      *ItemInformation         `xml:"ITEM_INFORMATION,omitempty"`
}

// BG-26: INVOICE_LINE_PERIOD
type InvoiceLinePeriod struct {
	XMLName                    xml.Name `xml:"INVOICE_LINE_PERIOD"`
	Id                         string   `xml:"xr:id,attr"`
	Src                        string   `xml:"xr:src,attr,omitempty"`
	InvoiceLinePeriodStartDate *Date    `xml:"Invoice_line_period_start_date,omitempty"`
	InvoiceLinePeriodEndDate   *Date    `xml:"Invoice_line_period_end_date,omitempty"`
}

// BG-27: INVOICE_LINE_ALLOWANCES
type InvoiceLineAllowances struct {
	XMLName                        xml.Name `xml:"INVOICE_LINE_ALLOWANCES"`
	Id                             string   `xml:"xr:id,attr"`
	Src                            string   `xml:"xr:src,attr,omitempty"`
	InvoiceLineAllowanceAmount     *Text    `xml:"Invoice_line_allowance_amount,omitempty"`
	InvoiceLineAllowanceBaseAmount *Text    `xml:"Invoice_line_allowance_base_amount,omitempty"`
	InvoiceLineAllowancePercentage *Text    `xml:"Invoice_line_allowance_percentage,omitempty"`
	InvoiceLineAllowanceReason     *Text    `xml:"Invoice_line_allowance_reason,omitempty"`
	InvoiceLineAllowanceReasonCode *Code    `xml:"Invoice_line_allowance_reason_code,omitempty"`
}

// BG-28: INVOICE_LINE_CHARGES
type InvoiceLineCharges struct {
	XMLName                     xml.Name `xml:"INVOICE_LINE_CHARGES"`
	Id                          string   `xml:"xr:id,attr"`
	Src                         string   `xml:"xr:src,attr,omitempty"`
	InvoiceLineChargeAmount     *Text    `xml:"Invoice_line_charge_amount,omitempty"`
	InvoiceLineChargeBaseAmount *Text    `xml:"Invoice_line_charge_base_amount,omitempty"`
	InvoiceLineChargePercentage *Text    `xml:"Invoice_line_charge_percentage,omitempty"`
	InvoiceLineChargeReason     *Text    `xml:"Invoice_line_charge_reason,omitempty"`
	InvoiceLineChargeReasonCode *Code    `xml:"Invoice_line_charge_reason_code,omitempty"`
}

// BG-29: PRICE_DETAILS
type PriceDetails struct {
	XMLName                            xml.Name `xml:"PRICE_DETAILS"`
	Id                                 string   `xml:"xr:id,attr"`
	Src                                string   `xml:"xr:src,attr,omitempty"`
	ItemNetPrice                       *Text    `xml:"Item_net_price,omitempty"`
	ItemPriceDiscount                  *Text    `xml:"Item_price_discount,omitempty"`
	ItemGrossPrice                     *Text    `xml:"Item_gross_price,omitempty"`
	ItemPriceBaseQuantity              *Text    `xml:"Item_price_base_quantity,omitempty"`
	ItemPriceBaseQuantityUnitOfMeasure *Code    `xml:"Item_price_base_quantity_unit_of_measure,omitempty"`
}

// BG-30: LINE_VAT_INFORMATION
type LineVATInformation struct {
	XMLName                     xml.Name `xml:"LINE_VAT_INFORMATION"`
	Id                          string   `xml:"xr:id,attr"`
	Src                         string   `xml:"xr:src,attr,omitempty"`
	InvoicedItemVATCategoryCode *Code    `xml:"Invoiced_item_VAT_category_code,omitempty"`
	InvoicedItemVATRate         *Text    `xml:"Invoiced_item_VAT_rate,omitempty"`
}

// BG-31: ITEM_INFORMATION
type ItemInformation struct {
	XMLName                      xml.Name              `xml:"ITEM_INFORMATION"`
	Id                           string                `xml:"xr:id,attr"`
	Src                          string                `xml:"xr:src,attr,omitempty"`
	ItemName                     *Text                 `xml:"Item_name,omitempty"`
	ItemDescription              *Text                 `xml:"Item_description,omitempty"`
	ItemSellersIdentifier        *Identifier           `xml:"Item_Sellers_identifier,omitempty"`
	ItemBuyersIdentifier         *Identifier           `xml:"Item_Buyers_identifier,omitempty"`
	ItemStandardIdentifier       *IdentifierWithScheme `xml:"Item_standard_identifier,omitempty"`
	ItemClassificationIdentifier *IdentifierWithScheme `xml:"Item_classification_identifier,omitempty"`
	ItemCountryOfOrigin          *Code                 `xml:"Item_country_of_origin,omitempty"`
	ItemAttributes               []*ItemAttributes     `xml:"ITEM_ATTRIBUTES,omitempty"`
}

// BG-32: ITEM_ATTRIBUTES
type ItemAttributes struct {
	XMLName            xml.Name `xml:"ITEM_ATTRIBUTES"`
	Id                 string   `xml:"xr:id,attr"`
	Src                string   `xml:"xr:src,attr,omitempty"`
	ItemAttributeName  *Text    `xml:"Item_attribute_name,omitempty"`
	ItemAttributeValue *Text    `xml:"Item_attribute_value,omitempty"`
}

func TransformXML(r io.Reader) (string, error) {
	// 1. Unmarshal the source XML
	var sourceXML struct {
		CrossIndustryInvoice struct {
			ExchangedDocument struct {
				ID            string `xml:"ID"`
				IssueDateTime struct {
					DateTimeString string `xml:"DateTimeString"`
				} `xml:"IssueDateTime"`
				TypeCode     string `xml:"TypeCode"`
				IncludedNote []struct {
					SubjectCode string `xml:"SubjectCode"`
					Content     string `xml:"Content"`
				} `xml:"IncludedNote"`
			} `xml:"ExchangedDocument"`
			ExchangedDocumentContext struct {
				BusinessProcessSpecifiedDocumentContextParameter struct {
					ID string `xml:"ID"`
				} `xml:"BusinessProcessSpecifiedDocumentContextParameter"`
				GuidelineSpecifiedDocumentContextParameter struct {
					ID string `xml:"ID"`
				} `xml:"GuidelineSpecifiedDocumentContextParameter"`
			} `xml:"ExchangedDocumentContext"`
			SupplyChainTradeTransaction struct {
				ApplicableHeaderTradeSettlement struct {
					InvoiceCurrencyCode string `xml:"InvoiceCurrencyCode"`
					TaxCurrencyCode     string `xml:"TaxCurrencyCode"`
					ApplicableTradeTax  []struct {
						TaxPointDate struct {
							DateString string `xml:"DateString"`
						} `xml:"TaxPointDate"`
						DueDateTypeCode       string `xml:"DueDateTypeCode"`
						BasisAmount           string `xml:"BasisAmount"`
						CalculatedAmount      string `xml:"CalculatedAmount"`
						CategoryCode          string `xml:"CategoryCode"`
						RateApplicablePercent string `xml:"RateApplicablePercent"`
						ExemptionReason       string `xml:"ExemptionReason"`
						ExemptionReasonCode   string `xml:"ExemptionReasonCode"`
					} `xml:"ApplicableTradeTax"`
					SpecifiedTradePaymentTerms []struct {
						DueDateDateTime struct {
							DateTimeString string `xml:"DateTimeString"`
						} `xml:"DueDateDateTime"`
						Description                         string `xml:"Description"`
						DirectDebitMandateID                string `xml:"DirectDebitMandateID"`
						ApplicableTradePaymentDiscountTerms struct {
							BasisPeriodMeasure struct {
								Text     string `xml:",chardata"`
								UnitCode string `xml:"unitCode,attr"`
							} `xml:"BasisPeriodMeasure"`
							CalculationPercent string `xml:"CalculationPercent"`
						} `xml:"ApplicableTradePaymentDiscountTerms"`
					} `xml:"SpecifiedTradePaymentTerms"`
					InvoiceReferencedDocument struct {
						IssuerAssignedID       string `xml:"IssuerAssignedID"`
						FormattedIssueDateTime struct {
							DateTimeString string `xml:"DateTimeString"`
						} `xml:"FormattedIssueDateTime"`
					} `xml:"InvoiceReferencedDocument"`
					ReceivableSpecifiedTradeAccountingAccount struct {
						ID string `xml:"ID"`
					} `xml:"ReceivableSpecifiedTradeAccountingAccount"`
					PaymentReference                     string `xml:"PaymentReference"`
					SpecifiedTradeSettlementPaymentMeans []struct {
						TypeCode                           string `xml:"TypeCode"`
						Information                        string `xml:"Information"`
						PayeePartyCreditorFinancialAccount struct {
							ProprietaryID string `xml:"ProprietaryID"`
							IBANID        string `xml:"IBANID"`
							AccountName   string `xml:"AccountName"`
						} `xml:"PayeePartyCreditorFinancialAccount"`
						ApplicableTradeSettlementFinancialCard struct {
							ID             string `xml:"ID"`
							CardholderName string `xml:"CardholderName"`
						} `xml:"ApplicableTradeSettlementFinancialCard"`
						PayeeSpecifiedCreditorFinancialInstitution struct {
							BICID string `xml:"BICID"`
						} `xml:"PayeeSpecifiedCreditorFinancialInstitution"`
						PayerPartyDebtorFinancialAccount struct {
							IBANID string `xml:"IBANID"`
						} `xml:"PayerPartyDebtorFinancialAccount"`
					} `xml:"SpecifiedTradeSettlementPaymentMeans"`
					SpecifiedTradeAllowanceCharge []struct {
						ChargeIndicator struct {
							Indicator string `xml:"Indicator"`
						} `xml:"ChargeIndicator"`
						ActualAmount       string `xml:"ActualAmount"`
						BasisAmount        string `xml:"BasisAmount"`
						CalculationPercent string `xml:"CalculationPercent"`
						Reason             string `xml:"Reason"`
						ReasonCode         string `xml:"ReasonCode"`
						CategoryTradeTax   struct {
							CategoryCode          string `xml:"CategoryCode"`
							RateApplicablePercent string `xml:"RateApplicablePercent"`
						} `xml:"CategoryTradeTax"`
					} `xml:"SpecifiedTradeAllowanceCharge"`
					SpecifiedTradeSettlementHeaderMonetarySummation struct {
						LineTotalAmount      string `xml:"LineTotalAmount"`
						AllowanceTotalAmount string `xml:"AllowanceTotalAmount"`
						ChargeTotalAmount    string `xml:"ChargeTotalAmount"`
						TaxBasisTotalAmount  string `xml:"TaxBasisTotalAmount"`
						TaxTotalAmount       []struct {
							Text       string `xml:",chardata"`
							CurrencyID string `xml:"currencyID,attr"`
						} `xml:"TaxTotalAmount"`
						GrandTotalAmount   string `xml:"GrandTotalAmount"`
						TotalPrepaidAmount string `xml:"TotalPrepaidAmount"`
						RoundingAmount     string `xml:"RoundingAmount"`
						DuePayableAmount   string `xml:"DuePayableAmount"`
					} `xml:"SpecifiedTradeSettlementHeaderMonetarySummation"`
					CreditorReferenceID    string `xml:"CreditorReferenceID"`
					BillingSpecifiedPeriod []struct {
						StartDateTime struct {
							DateTimeString string `xml:"DateTimeString"`
						} `xml:"StartDateTime"`
						EndDateTime struct {
							DateTimeString string `xml:"DateTimeString"`
						} `xml:"EndDateTime"`
					} `xml:"BillingSpecifiedPeriod"`
				} `xml:"ApplicableHeaderTradeSettlement"`
				ApplicableHeaderTradeAgreement struct {
					BuyerReference            string `xml:"BuyerReference"`
					SpecifiedProcuringProject struct {
						ID string `xml:"ID"`
					} `xml:"SpecifiedProcuringProject"`
					ContractReferencedDocument struct {
						IssuerAssignedID string `xml:"IssuerAssignedID"`
					} `xml:"ContractReferencedDocument"`
					BuyerOrderReferencedDocument struct {
						IssuerAssignedID string `xml:"IssuerAssignedID"`
					} `xml:"BuyerOrderReferencedDocument"`
					SellerOrderReferencedDocument struct {
						IssuerAssignedID string `xml:"IssuerAssignedID"`
					} `xml:"SellerOrderReferencedDocument"`
					AdditionalReferencedDocument []struct {
						IssuerAssignedID       string `xml:"IssuerAssignedID"`
						TypeCode               string `xml:"TypeCode"`
						ReferenceTypeCode      string `xml:"ReferenceTypeCode"`
						Name                   string `xml:"Name"`
						URIID                  string `xml:"URIID"`
						AttachmentBinaryObject struct {
							MimeCode string `xml:"mimeCode,attr"`
							Filename string `xml:"filename,attr"`
							Text     string `xml:",chardata"`
						} `xml:"AttachmentBinaryObject"`
					} `xml:"AdditionalReferencedDocument"`
					SellerTradeParty struct {
						Name                       string `xml:"Name"`
						SpecifiedLegalOrganization struct {
							TradingBusinessName string `xml:"TradingBusinessName"`
							ID                  string `xml:"ID"`
						} `xml:"SpecifiedLegalOrganization"`
						ID       []string `xml:"ID"`
						GlobalID []struct {
							Text     string `xml:",chardata"`
							SchemeID string `xml:"schemeID,attr"`
						} `xml:"GlobalID"`
						SpecifiedTaxRegistration []struct {
							ID struct {
								Text     string `xml:",chardata"`
								SchemeID string `xml:"schemeID,attr"`
							} `xml:"ID"`
						} `xml:"SpecifiedTaxRegistration"`
						Description               string `xml:"Description"`
						URIUniversalCommunication struct {
							URIID string `xml:"URIID"`
						} `xml:"URIUniversalCommunication"`
						PostalTradeAddress struct {
							LineOne                string `xml:"LineOne"`
							LineTwo                string `xml:"LineTwo"`
							LineThree              string `xml:"LineThree"`
							CityName               string `xml:"CityName"`
							PostcodeCode           string `xml:"PostcodeCode"`
							CountrySubDivisionName string `xml:"CountrySubDivisionName"`
							CountryID              string `xml:"CountryID"`
						} `xml:"PostalTradeAddress"`
						DefinedTradeContact struct {
							DepartmentName                  string `xml:"DepartmentName"`
							PersonName                      string `xml:"PersonName"`
							TelephoneUniversalCommunication struct {
								CompleteNumber string `xml:"CompleteNumber"`
							} `xml:"TelephoneUniversalCommunication"`
							EmailURIUniversalCommunication struct {
								URIID string `xml:"URIID"`
							} `xml:"EmailURIUniversalCommunication"`
						} `xml:"DefinedTradeContact"`
					} `xml:"SellerTradeParty"`
					BuyerTradeParty struct {
						Name                       string `xml:"Name"`
						SpecifiedLegalOrganization struct {
							TradingBusinessName string `xml:"TradingBusinessName"`
							ID                  string `xml:"ID"`
						} `xml:"SpecifiedLegalOrganization"`
						ID       []string `xml:"ID"`
						GlobalID []struct {
							Text     string `xml:",chardata"`
							SchemeID string `xml:"schemeID,attr"`
						} `xml:"GlobalID"`
						SpecifiedTaxRegistration []struct {
							ID struct {
								Text     string `xml:",chardata"`
								SchemeID string `xml:"schemeID,attr"`
							} `xml:"ID"`
						} `xml:"SpecifiedTaxRegistration"`
						URIUniversalCommunication struct {
							URIID string `xml:"URIID"`
						} `xml:"URIUniversalCommunication"`
						PostalTradeAddress struct {
							LineOne                string `xml:"LineOne"`
							LineTwo                string `xml:"LineTwo"`
							LineThree              string `xml:"LineThree"`
							CityName               string `xml:"CityName"`
							PostcodeCode           string `xml:"PostcodeCode"`
							CountrySubDivisionName string `xml:"CountrySubDivisionName"`
							CountryID              string `xml:"CountryID"`
						} `xml:"PostalTradeAddress"`
						DefinedTradeContact struct {
							DepartmentName                  string `xml:"DepartmentName"`
							PersonName                      string `xml:"PersonName"`
							TelephoneUniversalCommunication struct {
								CompleteNumber string `xml:"CompleteNumber"`
							} `xml:"TelephoneUniversalCommunication"`
							EmailURIUniversalCommunication struct {
								URIID string `xml:"URIID"`
							} `xml:"EmailURIUniversalCommunication"`
						} `xml:"DefinedTradeContact"`
					} `xml:"BuyerTradeParty"`
					SellerTaxRepresentativeTradeParty struct {
						Name                     string `xml:"Name"`
						SpecifiedTaxRegistration struct {
							ID string `xml:"ID"`
						} `xml:"SpecifiedTaxRegistration"`
						PostalTradeAddress struct {
							LineOne                string `xml:"LineOne"`
							LineTwo                string `xml:"LineTwo"`
							LineThree              string `xml:"LineThree"`
							CityName               string `xml:"CityName"`
							PostcodeCode           string `xml:"PostcodeCode"`
							CountrySubDivisionName string `xml:"CountrySubDivisionName"`
							CountryID              string `xml:"CountryID"`
						} `xml:"PostalTradeAddress"`
					} `xml:"SellerTaxRepresentativeTradeParty"`
					ApplicableHeaderTradeDelivery struct {
						ReceivingAdviceReferencedDocument struct {
							IssuerAssignedID string `xml:"IssuerAssignedID"`
						} `xml:"ReceivingAdviceReferencedDocument"`
						DespatchAdviceReferencedDocument struct {
							IssuerAssignedID string `xml:"IssuerAssignedID"`
						} `xml:"DespatchAdviceReferencedDocument"`
						ShipToTradeParty struct {
							Name     string   `xml:"Name"`
							ID       []string `xml:"ID"`
							GlobalID []struct {
								Text     string `xml:",chardata"`
								SchemeID string `xml:"schemeID,attr"`
							} `xml:"GlobalID"`
							PostalTradeAddress struct {
								LineOne                string `xml:"LineOne"`
								LineTwo                string `xml:"LineTwo"`
								LineThree              string `xml:"LineThree"`
								CityName               string `xml:"CityName"`
								PostcodeCode           string `xml:"PostcodeCode"`
								CountrySubDivisionName string `xml:"CountrySubDivisionName"`
								CountryID              string `xml:"CountryID"`
							} `xml:"PostalTradeAddress"`
						} `xml:"ShipToTradeParty"`
						ActualDeliverySupplyChainEvent struct {
							OccurrenceDateTime struct {
								DateTimeString string `xml:"DateTimeString"`
							} `xml:"OccurrenceDateTime"`
						} `xml:"ActualDeliverySupplyChainEvent"`
					} `xml:"ApplicableHeaderTradeDelivery"`
				} `xml:"ApplicableHeaderTradeAgreement"`
				IncludedSupplyChainTradeLineItem []struct {
					AssociatedDocumentLineDocument struct {
						LineID       string `xml:"LineID"`
						IncludedNote []struct {
							Content string `xml:"Content"`
						} `xml:"IncludedNote"`
					} `xml:"AssociatedDocumentLineDocument"`
					SpecifiedLineTradeSettlement struct {
						AdditionalReferencedDocument struct {
							IssuerAssignedID  string `xml:"IssuerAssignedID"`
							TypeCode          string `xml:"TypeCode"`
							ReferenceTypeCode string `xml:"ReferenceTypeCode"`
						} `xml:"AdditionalReferencedDocument"`
						SpecifiedTradeSettlementLineMonetarySummation struct {
							LineTotalAmount string `xml:"LineTotalAmount"`
						} `xml:"SpecifiedTradeSettlementLineMonetarySummation"`
						ReceivableSpecifiedTradeAccountingAccount struct {
							ID string `xml:"ID"`
						} `xml:"ReceivableSpecifiedTradeAccountingAccount"`
						BillingSpecifiedPeriod []struct {
							StartDateTime struct {
								DateTimeString string `xml:"DateTimeString"`
							} `xml:"StartDateTime"`
							EndDateTime struct {
								DateTimeString string `xml:"DateTimeString"`
							} `xml:"EndDateTime"`
						} `xml:"BillingSpecifiedPeriod"`
						SpecifiedTradeAllowanceCharge []struct {
							ChargeIndicator struct {
								Indicator string `xml:"Indicator"`
							} `xml:"ChargeIndicator"`
							ActualAmount       string `xml:"ActualAmount"`
							BasisAmount        string `xml:"BasisAmount"`
							CalculationPercent string `xml:"CalculationPercent"`
							Reason             string `xml:"Reason"`
							ReasonCode         string `xml:"ReasonCode"`
						} `xml:"SpecifiedTradeAllowanceCharge"`
						ApplicableTradeTax struct {
							CategoryCode          string `xml:"CategoryCode"`
							RateApplicablePercent string `xml:"RateApplicablePercent"`
						} `xml:"ApplicableTradeTax"`
					} `xml:"SpecifiedLineTradeSettlement"`
					SpecifiedLineTradeDelivery struct {
						BilledQuantity struct {
							Text     string `xml:",chardata"`
							UnitCode string `xml:"unitCode,attr"`
						} `xml:"BilledQuantity"`
					} `xml:"SpecifiedLineTradeDelivery"`
					SpecifiedLineTradeAgreement struct {
						NetPriceProductTradePrice struct {
							ChargeAmount  string `xml:"ChargeAmount"`
							BasisQuantity struct {
								Text     string `xml:",chardata"`
								UnitCode string `xml:"unitCode,attr"`
							} `xml:"BasisQuantity"`
						} `xml:"NetPriceProductTradePrice"`
						GrossPriceProductTradePrice struct {
							AppliedTradeAllowanceCharge struct {
								ActualAmount string `xml:"ActualAmount"`
							} `xml:"AppliedTradeAllowanceCharge"`
							ChargeAmount  string `xml:"ChargeAmount"`
							BasisQuantity struct {
								Text     string `xml:",chardata"`
								UnitCode string `xml:"unitCode,attr"`
							} `xml:"BasisQuantity"`
						} `xml:"GrossPriceProductTradePrice"`
						BuyerOrderReferencedDocument struct {
							LineID string `xml:"LineID"`
						} `xml:"BuyerOrderReferencedDocument"`
					} `xml:"SpecifiedLineTradeAgreement"`
					SpecifiedTradeProduct struct {
						Name             string `xml:"Name"`
						Description      string `xml:"Description"`
						SellerAssignedID string `xml:"SellerAssignedID"`
						BuyerAssignedID  string `xml:"BuyerAssignedID"`
						GlobalID         struct {
							Text     string `xml:",chardata"`
							SchemeID string `xml:"schemeID,attr"`
						} `xml:"GlobalID"`
						DesignatedProductClassification struct {
							ClassCode struct {
								Text          string `xml:",chardata"`
								ListID        string `xml:"listID,attr"`
								ListVersionID string `xml:"listVersionID,attr"`
							} `xml:"ClassCode"`
						} `xml:"DesignatedProductClassification"`
						OriginTradeCountry struct {
							ID string `xml:"ID"`
						} `xml:"OriginTradeCountry"`
						ApplicableProductCharacteristic []struct {
							Description string `xml:"Description"`
							Value       string `xml:"Value"`
						} `xml:"ApplicableProductCharacteristic"`
					} `xml:"SpecifiedTradeProduct"`
				} `xml:"IncludedSupplyChainTradeLineItem"`
			} `xml:"SupplyChainTradeTransaction"`
		} `xml:"CrossIndustryInvoice"`
	}

	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&sourceXML); err != nil {
		return "", fmt.Errorf("error decoding XML: %w", err)
	}

	// 2. Transform into the target XML structure
	targetXML := Invoice{
		XMLName: xml.Name{Local: "xr:invoice"},
		InvoiceNumber: &Identifier{
			XMLName: xml.Name{Local: "xr:Invoice_number"},
			Text:    sourceXML.CrossIndustryInvoice.ExchangedDocument.ID,
		},
		InvoiceIssueDate: &Date{
			XMLName: xml.Name{Local: "xr:Invoice_issue_date"},
			Text:    formatDate(sourceXML.CrossIndustryInvoice.ExchangedDocument.IssueDateTime.DateTimeString),
		},
		InvoiceTypeCode: &Code{
			XMLName: xml.Name{Local: "xr:Invoice_type_code"},
			// ... (Logic for InvoiceTypeCode) ...
		},
		// ... (Other fields) ...
	}

	// ... (Processing for nested structures like INVOICE_NOTE, etc.) ...

	// 3. Marshal the target XML
	output, err := xml.MarshalIndent(targetXML, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling XML: %w", err)
	}

	return string(output), nil
}

// Helper functions for formatting, mapping, etc. (e.g., formatDate)
func formatDate(dateStr string) string {
	// Implement logic to format the date (e.g., remove hyphens)
	return strings.ReplaceAll(dateStr, "-", "")
}

// ... (Other helper functions as needed) ...
