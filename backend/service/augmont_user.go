package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/nleeper/goment"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/utils"
)

// AumontService provides augmont merchant api functionality
type augmontService struct {
	user  interfaces.AugmontUserRepo
	order interfaces.AugmontOrderRepo
	inMem interfaces.AugmontInMemRepo
}

type AugmontLogInResp struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`

	Result struct {
		Data struct {
			MerchantID int    `json:"merchantId"`
			Token      string `json:"accessToken"`
			ExpireAt   string `json:"expireAt"`
		} `json:"data"`
	} `json:"result"`
}

// Create New Augmond Service
func NewAugmondService(
	user interfaces.AugmontUserRepo,
	order interfaces.AugmontOrderRepo,
	inMem interfaces.AugmontInMemRepo,
) interfaces.AugmontService {
	return &augmontService{
		user:  user,
		order: order,
		inMem: inMem,
	}
}

func (augmontService) newUniqueID() string {
	uniqueIDMaxLen := 30
	return utils.NewUniqueString(uniqueIDMaxLen)
}

func (augmontService) newTnxID() string {
	TnxIDMaxLen := 30
	return utils.NewUniqueString(TnxIDMaxLen)
}

// Login authenticates user with augmont
func (s *augmontService) logIn() (*AugmontLogInResp, error) {
	// get config
	email := domain.Config().Augmont.Email
	password := domain.Config().Augmont.Password
	host := domain.Config().Augmont.Host
	// parse augmont login augUrl

	augUrl := host + "/merchant/v1/auth/login"

	resp, err := http.PostForm(augUrl, url.Values{
		"email":    {email},
		"password": {password},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respData := AugmontLogInResp{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	if respData.StatusCode != 200 {
		err = fmt.Errorf("%v", respData.Message)
		return nil, domain.NewError(err, domain.ErrInvalidArgument)
	}
	return &respData, nil
}

// AuthToken returns authenication token for augmont
func (s *augmontService) AuhtToken() (string, error) {
	// Retive the stored token
	token, err := s.inMem.GetToken()
	if err != nil {
		return "", err
	}

	// If token is not expired, return it
	if token != "" {
		return token, nil
	}

	// If token is expired, generate new token
	data, err := s.logIn()
	if err != nil {
		return "", err
	}
	token = data.Result.Data.Token
	gom, _ := goment.New(data.Result.Data.ExpireAt, "YYYY-MM-DD HH-mm-ss")
	expireAt := gom.ToTime()
	s.inMem.SetToken(token, expireAt)

	return token, nil
}

// CreateUser creates a new customer account
// with augmont and create augmont user in db
func (s *augmontService) CreateUser(
	userInfo *utils.AugmontUserInfo,
	user *models.AugmontUser,
) error {
	// Verify if AppUser has Augmont Account

	// Generate uniqueID
	{
		uniqueID := s.newUniqueID()
		userInfo.UniqueID = uniqueID
		user.UID = &uniqueID
	}

	// Make Post Req to Augmont
	host := domain.Config().Augmont.Host
	augUrl := host + "/merchant/v1/users"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	{
		_ = writer.WriteField("uniqueId", userInfo.UniqueID)
		_ = writer.WriteField("userName", userInfo.Name)
		if (userInfo.MobileNo) != "" {
			_ = writer.WriteField("mobileNumber", userInfo.MobileNo)
		}
		if (userInfo.EmailID) != "" {
			_ = writer.WriteField("emailId", userInfo.EmailID)
		}
		if userInfo.City != "" {
			_ = writer.WriteField("userCity", userInfo.City)
		}
		if userInfo.State != "" {
			_ = writer.WriteField("userState", userInfo.State)
		}
		if userInfo.Pincode != "" {
			_ = writer.WriteField("userPincode", userInfo.Pincode)
		}
		if userInfo.DOB != "" {
			_ = writer.WriteField("dateOfBirth", userInfo.DOB)
		}
		if userInfo.NomineeName != "" {
			_ = writer.WriteField("nomineeName", userInfo.NomineeName)
		}
		if userInfo.NomineeDOB != "" {
			_ = writer.WriteField("nomineeDateOfBirth", userInfo.NomineeDOB)
		}
		if userInfo.NomineeRelation != "" {
			_ = writer.WriteField("nomineeRelation", userInfo.NomineeRelation)
		}
		if userInfo.UtmSource != "" {
			_ = writer.WriteField("utmSource", userInfo.UtmSource)
		}
		if userInfo.UtmMedium != "" {
			_ = writer.WriteField("utmMedium", userInfo.UtmMedium)
		}
		if userInfo.UtmCampaign != "" {
			_ = writer.WriteField("utmCampaign", userInfo.UtmCampaign)
		}
	}
	// Create New Request with payload
	req, err := http.NewRequest(method, augUrl, payload)
	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	// Make Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If success, create user in db
	err = s.user.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates user info in augmont with uniqueID
func (s *augmontService) UpdateUser(
	userInfo *utils.AugmontUserInfo,
) error {

	// get user data from augmont, replace empty fields with userInfo
	{
		augUser, err := s.GetUserInfo(userInfo.UniqueID)
		if err != nil {
			return errors.WithMessage(err, "failed to get user data from augmont")
		}
		augUser.City = ""
		augUser.State = ""
		utils.CopyNonEmptyFiled(augUser, userInfo)
	}
	uniqueID := userInfo.UniqueID
	host := domain.Config().Augmont.Host

	augUrl := host + "/merchant/v1/users/" + uniqueID
	method := "PUT"

	// Omit UniqueID from payload, not need for update user
	userInfo.UniqueID = ""
	infoDict := utils.GetNonEmptyFields(userInfo)
	payload := strings.NewReader(infoDict.ToString())

	// Create New Request with payload
	req, err := http.NewRequest(method, augUrl, payload)
	if err != nil {
		return err
	}
	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")
	}
	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	// If not success, return error messsage
	if data["statusCode"] != float64(200) {
		dict := utils.Dict(data["errors"].(map[string]interface{}))
		return domain.NewError(err, domain.ErrInvalidArgument, dict.ToString())
	}
	return nil
}

// GetUserInfo returns user info from augmont
func (s *augmontService) GetUserInfo(
	uniqueID string,
) (
	*utils.AugmontUserInfo,
	error,
) {
	host := domain.Config().Augmont.Host
	url := host + "/merchant/v1/users/" + uniqueID

	// Create Get Request
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
	}
	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, domain.NewError(err, domain.ErrInternalError)
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	// Check if the  request is success

	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["message"])
		return nil, domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Decode userInfo from response result
	result, ok := data["result"].(map[string]interface{})
	if !ok {
		err = errors.New("type conversion failed for result")
		return nil, domain.NewError(err, domain.ErrInternalError)
	}
	info := utils.AugmontUserInfo{}
	mapstructure.Decode(result["data"], &info)
	return &(info), nil
}

func (s *augmontService) PostUserKyc(
	name, pan, dob string,
	user *models.AugmontUser,
	file *utils.File,
) (utils.Any, error) {

	host := domain.Config().Augmont.Host
	url := fmt.Sprintf("%v/merchant/v1/users/%v/kyc", host, *user.UID)
	method := "POST"

	// Create Request Payload
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	{
		_ = writer.WriteField("panNumber", pan)
		_ = writer.WriteField("dateOfBirth", dob)
		_ = writer.WriteField("nameAsPerPan", name)

		// Copy file to payload
		part, _ := writer.CreateFormFile("panAttachment", file.Path())
		fileio, err := file.Open()
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, fileio)
		if err != nil {
			return nil, err
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}
	}

	// Create New Request with payload
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := &utils.AugmontResponse{}

	json.NewDecoder(resp.Body).Decode(&data)
	if data.IsError() {
		return nil, data.Error()
	}

	status := "pending"
	s.user.UpdateUser(&models.AugmontUser{
		ID:        user.UserID,
		KYCStatus: &status,
	})

	return data.Result, nil
}

func (s *augmontService) UpdateUserKycStatus(
	user *models.AugmontUser,
) error {

	// Update only if user KYC is Pending, else return
	if user.KYCStatus != nil && *user.KYCStatus != "pending" {
		return nil
	}

	host := domain.Config().Augmont.Host
	url := fmt.Sprintf("%v/merchant/v1/users/%v/kyc", host, *user.UID)
	method := "GET"

	// Create Get Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}
	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["message"])
		return domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Update user kyc status
	status := "approved"
	err = s.user.UpdateUser(&models.AugmontUser{
		ID:        user.UserID,
		KYCStatus: &status,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *augmontService) CreateUserBank(
	user *models.AugmontUser,
	bankInfo *utils.AugmontUserBankInfo,
) error {

	if len(user.Banks) >= 10 {
		err := fmt.Errorf("user can have maximum 10 banks")
		return domain.NewError(err, domain.ErrBadRequest)
	}

	host := domain.Config().Augmont.Host
	uniqueID := *user.UID

	// Prepare request body
	agUrl := host + "/merchant/v1/users/" + uniqueID + "/banks"

	info := utils.GetNonEmptyFields(bankInfo)
	log.Println(info.ToUrlString())
	payload := strings.NewReader(info.ToUrlString())

	// Create New Request
	req, err := http.NewRequest(http.MethodPost, agUrl, payload)

	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(info.ToUrlString())))
	}

	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["errors"])
		return domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Get user Bank Info
	result := data["result"].(map[string]interface{})
	bank := utils.AugmontUserBankInfo{}
	mapstructure.Decode(result["data"], &bank)

	// Update user bank info
	userBank := models.AugmontUserBank{
		UserBankID:    &bank.UserBankID,
		AugmontUserID: user.ID,
	}
	s.user.CreateBank(&userBank)

	return nil
}

func (s *augmontService) UpdateUserBank(
	user *models.AugmontUser,
	bankInfo *utils.AugmontUserBankInfo,
) error {

	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/users/%v/banks/%v",
		domain.Config().Augmont.Host,
		*user.UID,
		bankInfo.UserBankID,
	)
	method := "PUT"
	info := utils.GetNonEmptyFields(bankInfo)
	payload := strings.NewReader(info.ToUrlString())

	// Create New Request
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(info.ToUrlString())))
	}

	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	if data["statusCode"] != float64(200) {
		errs := utils.Dict(data["errors"].(map[string]interface{}))
		err = fmt.Errorf("%v", errs.ToString())
		return domain.NewError(err, domain.ErrInvalidArgument)
	}

	return nil
}

func (s *augmontService) DeleteUserBank(
	user *models.AugmontUser,
	bankInfo *utils.AugmontUserBankInfo,
) error {

	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/users/%v/banks/%v",
		domain.Config().Augmont.Host,
		*user.UID,
		bankInfo.UserBankID,
	)
	method := "DELETE"

	// Create New Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
	}

	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["message"])
		return domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Delete user bank info
	s.user.DeleteBank(&models.AugmontUserBank{
		AugmontUserID: user.ID,
		UserBankID:    &bankInfo.UserBankID,
	})

	return nil
}

func (s *augmontService) GetUserBanks(
	user *models.AugmontUser,
) (
	[]*utils.AugmontUserBankInfo,
	error,
) {

	// Prepare request
	url := fmt.Sprintf("%v/merchant/v1/users/%v/banks",
		domain.Config().Augmont.Host,
		*user.UID,
	)
	method := "GET"

	// Create New Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}

	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	// handle error in response
	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["message"])
		return nil, domain.NewError(err, domain.ErrInvalidArgument)
	}

	// parse response bank info
	var userBanks []*utils.AugmontUserBankInfo
	result, ok := data["result"].([]interface{})
	// if Not ok, the user has no bank data
	if !ok {
		return userBanks, nil
	}

	for _, bank := range result {
		bankInfo := utils.AugmontUserBankInfo{}
		mapstructure.Decode(bank, &bankInfo)
		userBanks = append(userBanks, &bankInfo)
	}

	return userBanks, nil
}

func (s *augmontService) CreateUserAddress(
	user *models.AugmontUser,
	addressInfo *utils.AugmontUserAddressInfo,
) error {

	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/users/%v/address",
		domain.Config().Augmont.Host,
		*user.UID,
	)
	method := "POST"
	info := utils.GetNonEmptyFields(addressInfo)
	payload := strings.NewReader(info.ToUrlString())

	// Create New Request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&data)
	// handle error in response
	if data["statusCode"] != float64(200) {
		errs := utils.Dict(data["errors"].(map[string]interface{}))
		err = fmt.Errorf("%v", errs.ToString())
		return domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Save user address info
	result := data["result"].(map[string]interface{})
	address := utils.AugmontUserAddressInfo{}
	mapstructure.Decode(result, &address)
	err = s.user.CreateAddress(&models.AugmontUserAddress{
		AugmontUserID: user.ID,
		UserAddressID: &address.UserAddressID,
	})

	return err
}

func (s *augmontService) DeleteUserAddress(
	user *models.AugmontUser,
	addressInfo *utils.AugmontUserAddressInfo,
) error {

	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/users/%v/address/%v",
		domain.Config().Augmont.Host,
		*user.UID,
		addressInfo.UserAddressID,
	)
	method := "DELETE"

	// Create New Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}

	// Make Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["message"])
		return domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Delete user address info
	err = s.user.DeleteAddress(&models.AugmontUserAddress{
		AugmontUserID: user.ID,
		UserAddressID: &addressInfo.UserAddressID,
	})

	return err
}

func (s *augmontService) GetUserAddresses(
	user *models.AugmontUser,
) (
	[]*utils.AugmontUserAddressInfo,
	error,
) {

	// Prepare request
	url := fmt.Sprintf("%v/merchant/v1/users/%v/address",
		domain.Config().Augmont.Host,
		*user.UID,
	)
	method := "GET"

	// Create New Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode Response body
	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)
	// handle error in response
	if data["statusCode"] != float64(200) {
		errs := utils.Dict(data["errors"].(map[string]interface{}))
		err = fmt.Errorf("%v", errs.ToString())
		return nil, domain.NewError(err, domain.ErrInvalidArgument)
	}

	// parse response bank info
	var userAddresses []*utils.AugmontUserAddressInfo
	result, ok := data["result"].([]interface{})
	// not ok if user have no bank info
	if !ok {
		return userAddresses, nil
	}
	for _, address := range result {
		addressInfo := utils.AugmontUserAddressInfo{}
		mapstructure.Decode(address, &addressInfo)
		userAddresses = append(userAddresses, &addressInfo)
	}

	return userAddresses, nil
}

func (s *augmontService) Buy(
	user *models.AugmontUser,
	buyInfo *utils.AugmontBugInfo,
) (utils.Any, error) {
	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/buy",
		domain.Config().Augmont.Host,
	)
	method := "POST"

	// Create New Merchant Transaction Id, Should be unique
	{
		tnxID := s.newTnxID()
		buyInfo.MerchantTnxID = tnxID

		buyInfo.UserUID = *user.UID
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	{
		order := utils.GetNonEmptyFields(buyInfo)
		for key, value := range order {
			if value == "" {
				continue
			}
			if err := writer.WriteField(key, value.(string)); err != nil {
				return nil, err
			}
		}
		err := writer.Close()
		if err != nil {
			return nil, err
		}
	}

	// Create New Request
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	// Make Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode Response body
	data := utils.AugmontResponse{}
	json.NewDecoder(res.Body).Decode(&data)

	// handle error in response
	if data.IsError() {
		return nil, data.Error()
	}

	// Update buy orders table
	err = s.order.CreateBuy(&models.AugmontBuyOrder{
		AugmontUserID: user.ID,
		MerchantTxnID: &buyInfo.MerchantTnxID,
	})

	return data.Result, err
}

func (s *augmontService) getOrder(
	url string,
) (utils.Any, error) {
	// Prepare request
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", " Bearer "+token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode Response body
	data := utils.Dict{}
	json.NewDecoder(res.Body).Decode(&data)

	if data["statusCode"] != float64(200) {
		err = fmt.Errorf("%v", data["message"])
		return nil, domain.NewError(err, domain.ErrInvalidArgument)
	}

	// Get embedder type
	result := data.Get("result", "data")

	return result, err
}

func (s *augmontService) BuyInfo(
	userUniqueID,
	tnxID string,
) (utils.Any, error) {
	url := fmt.Sprintf("%v/merchant/v1/buy/%v/%v",
		domain.Config().Augmont.Host,
		userUniqueID,
		tnxID,
	)

	result, err := s.getOrder(url)
	return result, err
}

func (s *augmontService) BuyList(userUniqueID string) (utils.Any, error) {
	url := fmt.Sprintf("%v/merchant/v1/%v/buy",
		domain.Config().Augmont.Host,
		userUniqueID,
	)

	result, err := s.getOrder(url)
	return result, err
}

func (s *augmontService) Sell(
	user *models.AugmontUser,
	sellInfo *utils.AugmontSellInfo,
) (utils.Any, error) {

	// Generate New Transaction ID
	{
		tnxID := s.newTnxID()
		sellInfo.MerchantTnxID = tnxID
	}

	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/sell",
		domain.Config().Augmont.Host,
	)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	{
		order := utils.GetNonEmptyFields(sellInfo)
		for key, value := range order {
			if value == "" {
				continue
			}
			if err := writer.WriteField(key, value.(string)); err != nil {
				return nil, err
			}
		}
		err := writer.Close()
		if err != nil {
			return nil, err
		}
	}

	// Create New Request
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("token", token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	// Make Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode Response body
	data := utils.AugmontResponse{}
	json.NewDecoder(res.Body).Decode(&data)

	// handle error in response
	if data.IsError() {
		return nil, data.Error()
	}

	// Update sell orders table
	err = s.order.CreateSell(&models.AugmontSellOrder{
		AugmontUserID: user.ID,
		MerchantTxnID: &sellInfo.MerchantTnxID,
	})

	return data.Result, err
}

func (s *augmontService) SellInfo(
	userUniqueID,
	tnxID string,
) (utils.Any, error) {
	url := fmt.Sprintf("%v/merchant/v1/sell/%v/%v",
		domain.Config().Augmont.Host,
		tnxID,
		userUniqueID,
	)
	resp, err := s.getOrder(url)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (s *augmontService) SellList(
	userUniqueID string,
) (utils.Any, error) {

	url := fmt.Sprintf("%v/merchant/v1/%v/sell",
		domain.Config().Augmont.Host,
		userUniqueID,
	)
	resp, err := s.getOrder(url)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (s *augmontService) Redeem(
	user *models.AugmontUser,
	redeemInfo *utils.AugmontRedeemInfo,
) (utils.Any, error) {

	// Generate New Transaction ID
	{
		tnxID := s.newTnxID()
		redeemInfo.MerchantTnxID = tnxID
	}

	// Prepare request body
	url := fmt.Sprintf("%v/merchant/v1/order",
		domain.Config().Augmont.Host,
	)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	{
		order := utils.GetNonEmptyFields(redeemInfo)
		for key, value := range order {
			if value == "" {
				continue
			}
			if err := writer.WriteField(key, value.(string)); err != nil {
				return nil, err
			}
		}
		err := writer.Close()
		if err != nil {
			return nil, err
		}
	}
	// Create New Request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	// Add Request Headers
	{
		token, err := s.AuhtToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	// Make Request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := utils.AugmontResponse{}
	json.NewDecoder(res.Body).Decode(&data)

	// handle error in response
	if data.IsError() {
		return nil, data.Error()
	}

	// Update sell orders table
	err = s.order.CreateRedeem(&models.AugmontRedeemOrder{
		AugmontUserID: user.ID,
		MerchantTxnID: &redeemInfo.MerchantTnxID,
	})

	return data.Result, err
}

func (s *augmontService) RedeemInfo(
	userUniqueID,
	tnxID string,
) (utils.Any, error) {
	url := fmt.Sprintf("%v/merchant/v1/order/%v/%v",
		domain.Config().Augmont.Host,
		tnxID,
		userUniqueID,
	)
	resp, err := s.getOrder(url)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (s *augmontService) RedeemList(userUniqueID string) (utils.Any, error) {

	url := fmt.Sprintf("%v/merchant/v1/%v/order",
		domain.Config().Augmont.Host,
		userUniqueID,
	)
	resp, err := s.getOrder(url)
	if err != nil {
		return nil, err
	}

	return resp, err
}
