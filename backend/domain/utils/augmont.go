package utils

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/mitchellh/mapstructure"
)

// Model for get, validate api data, and presenting
type AugmontUserInfo struct {
	MobileNo string `json:"mobileNumber,omitempty" mapstructure:"mobileNumber"`
	EmailID  string `json:"emailId,omitempty" mapstructure:"emailId"`
	UniqueID string `json:"uniqueId,omitempty" mapstructure:"uniqueId"`
	Name     string `json:"userName,omitempty" mapstructure:"userName"`
	City     string `json:"userCity,omitempty" mapstructure:"userCity"`
	State    string `json:"userState,omitempty" mapstructure:"userState"`
	Pincode  string `json:"userPincode,omitempty" mapstructure:"userPincode"`
	DOB      string `json:"dateOfBirth,omitempty" mapstructure:"dateOfBirth"`

	NomineeName     string `json:"nomineeName,omitempty" mapstructure:"nomineeName"`
	NomineeDOB      string `json:"nomineeDateOfBirth,omitempty" mapstructure:"nomineeDateOfBirth"`
	NomineeRelation string `json:"nomineeRelation,omitempty" mapstructure:"nomineeRelation"`

	UtmSource   string `json:"utmSource,omitempty" mapstructure:"utmSource"`
	UtmMedium   string `json:"utmMedium,omitempty" mapstructure:"utmMedium"`
	UtmCampaign string `json:"utmCampaign,omitempty" mapstructure:"utmCampaign"`
}

type AugmontUserBankInfo struct {
	UserBankID string `json:"userBankId,omitempty" mapstructure:"userBankId"`
	AccNo      string `json:"accountNumber,omitempty" mapstructure:"accountNumber"`
	AccName    string `json:"accountName,omitempty" mapstructure:"accountName"`
	Ifsc       string `json:"ifscCode,omitempty" mapstructure:"ifscCode"`
}

type AugmontUserAddressInfo struct {
	UserAddressID string `json:"userAddressId,omitempty" mapstructure:"userAddressId"`
	Name          string `json:"name,omitempty" mapstructure:"name"`
	MobileNo      string `json:"mobileNumber,omitempty" mapstructure:"mobileNumber"`
	Email         string `json:"email,omitempty" mapstructure:"email"`
	Address       string `json:"address,omitempty" mapstructure:"address"`
	Pincode       string `json:"pincode,omitempty" mapstructure:"pincode"`
}

type AugmontBugInfo struct {
	LockPrice     string `json:"lockPrice" `
	MetalType     string `json:"metalType" `
	Quantity      string `json:"quantity" `
	Amount        string `json:"amount" `
	MerchantTnxID string `json:"merchantTransactionId" `
	BlockID       string `json:"blockId"`

	PaymentMode string `json:"modeOfPayment"`
	UserUID     string `json:"uniqueId"`
	RefType     string `json:"referenceType"`
	RefID       string `json:"referenceId"`

	UtmSource   string `json:"utmSource"`
	UtmMedium   string `json:"utmMedium"`
	UtmCampaign string `json:"utmCampaign"`
}

type AugmontSellInfo struct {
	LockPrice     string `json:"lockPrice" `
	MetalType     string `json:"metalType" `
	Quantity      string `json:"quantity" `
	Amount        string `json:"amount" `
	MerchantTnxID string `json:"merchantTransactionId" `
	BlockID       string `json:"blockId"`

	UserBankID string `json:"userBankId"`
	AccNo      string `json:"accountNumber"`
	AccName    string `json:"accountName"`
	Ifsc       string `json:"ifscCode"`
}

type AugmontRedeemInfo struct {
	MobileNo      string               `json:"mobileNumber"`
	UserAddressID string               `json:"userAddressId"`
	Product       []AugmontProductInfo `json:"product" binding:"required"`

	PaymentMode   string `json:"modeOfPayment"`
	MerchantTnxID string `json:"merchantTransactionId" `
}

func (r *AugmontRedeemInfo) Write(writer *multipart.Writer) (err error) {
	for index, p := range r.Product {
		p.Write(writer, index)
	}
	return
}

type AugmontProductInfo struct {
	SKU      string `json:"sku" `
	Quantity string `json:"quantity" `
}

func (p *AugmontProductInfo) Write(writer *multipart.Writer, index int) (err error) {

	writer.WriteField(
		fmt.Sprintf("product[%v][sku]", index),
		p.SKU,
	)
	writer.WriteField(
		fmt.Sprintf("product[%v][quantity]", index),
		p.Quantity,
	)
	return
}

func CopyNonEmptyFiled(src interface{}, dest interface{}) {
	src1, _ := json.Marshal(src)
	dect1, _ := json.Marshal(dest)
	objSrc := Dict{}
	objDest := Dict{}
	json.Unmarshal(src1, &objSrc)
	json.Unmarshal(dect1, &objDest)
	for k, v := range objSrc {
		if _, ok := objDest[k]; v == "" || ok {
			continue
		}
		objDest[k] = v
	}
	mapstructure.Decode(objDest, dest)
}

func GetNonEmptyFields(obj interface{}) Dict {
	jsb, _ := json.Marshal(obj)
	var dict Dict
	json.Unmarshal(jsb, &dict)
	return dict
}
