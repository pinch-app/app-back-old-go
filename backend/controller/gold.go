package controller

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/utils"
)

type GoldController struct {
	gold        interfaces.AugmontService
	augmontUser interfaces.AugmontUserRepo
}

func NewGoldController(router *gin.Engine, gold interfaces.AugmontService, au interfaces.AugmontUserRepo) {
	c := &GoldController{
		gold:        gold,
		augmontUser: au,
	}

	// Profile Endpoints
	{
		group := router.Group("/gold/profile")
		// Create Augmont Client User
		// access -> user/admin
		group.POST("", c.CreateProfile)

		// Get profile data of logged in user
		// access -> user
		group.GET("", c.GetProfile)

		// Update profile data of logged in user
		// access -> user
		group.PUT("", c.UpdateProfile)
	}

	// Banks Endpoints
	{
		group := router.Group("/gold/profile/bank")
		// Create banks
		group.POST("", c.CreateBank)
		// Get all user banks
		group.GET("", c.GetUserBank)
		// Delete user bank
		group.DELETE(":userBankID", c.DeleteBank)
		// Update user bank
		group.PUT("", c.UpdateBank)
	}
	// Address Endpoint
	{
		group := router.Group("/gold/profile/address")
		// Create address
		group.POST("", c.CreateAddress)
		// Get all user address
		group.GET("", c.GetUserAddress)
		// Delete user address
		group.DELETE(":userAddressID", c.DeleteuserAddress)
	}

	{
		group := router.Group("/gold/profile/kyc")
		group.POST("", c.CreateKYC)
		group.GET("", c.GetKycStatus)
	}

	{
		group := router.Group("/gold/buy")
		group.POST("", c.BuyOrder)
		group.GET("/order/:taxID", c.GetBuyInfo)
		group.GET("/order", c.GetBuyList)
	}

}

// CreateProfile handle create profile request
func (c *GoldController) CreateProfile(ctx *gin.Context) {
	// Get Pinch User from Contex
	var user *models.User
	{
		userVal, ok := ctx.Get("user")
		if !ok {
			return
		}
		user = userVal.(*models.User)
	}

	log.Println("Create profile")

	augUser := &models.AugmontUser{
		UserID: user.ID,
	}

	userinfo := &utils.AugmontUserInfo{}
	err := ctx.Bind(userinfo)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	err = c.gold.CreateUser(userinfo, augUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "user profile created",
	})
}

func (c *GoldController) GetProfile(ctx *gin.Context) {

	// Get Pinch User from Context
	var user *models.User
	{
		var ok bool
		user, ok = ctx.Keys["user"].(*models.User)
		if !ok {
			err := errors.New("user not found in context")
			domain.ErrLog(err)
			domain.ErrFailedGinReq(ctx, err)
			return
		}
	}

	// Fetch augmont user from database
	augUser, err := c.augmontUser.FindUser(
		&models.AugmontUser{
			UserID: user.ID,
		})

	// handle error
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	info, err := c.gold.GetUserInfo(*augUser.UID)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "user profile fetched",
		"profile": info,
	})

}

func (c *GoldController) UpdateProfile(ctx *gin.Context) {

	// Get Pinch User from Context
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	// Fetch augmont user from database
	augUser, err := c.augmontUser.FindUser(
		&models.AugmontUser{
			UserID: user.ID,
		})

	// handle error
	if err != nil {
		err := errors.Wrap(err, "failed to fetch augmont user")
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	userinfo := &utils.AugmontUserInfo{}
	err = ctx.Bind(userinfo)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	userinfo.UniqueID = *augUser.UID

	err = c.gold.UpdateUser(userinfo)
	if err != nil {
		err := errors.Wrap(err, "failed to update augmont user info")
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "user profile updated",
	})
}

func (c *GoldController) CreateBank(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	bankInfo := &utils.AugmontUserBankInfo{}
	if err := ctx.Bind(bankInfo); err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	err = c.gold.CreateUserBank(agUser, bankInfo)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":     "ok",
		"userBankID": bankInfo.UserBankID,
	})
}

func (c *GoldController) GetUserBank(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	banks, err := c.gold.GetUserBanks(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"banks":  banks,
	})
}

func (c *GoldController) UpdateBank(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	bank := &utils.AugmontUserBankInfo{}
	if err := ctx.Bind(bank); err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	err = c.gold.UpdateUserBank(agUser, bank)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "ok",
		"message": "Banks data updated successfully",
	})
}

func (c *GoldController) DeleteBank(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	userBankID := ctx.Param("userBankID")
	if userBankID == "" {
		err = errors.New("userBankID not vaild")
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	bank := &utils.AugmontUserBankInfo{
		UserBankID: userBankID,
	}
	err = c.gold.DeleteUserBank(agUser, bank)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":     "ok",
		"userBankID": userBankID,
	})
}

func (c *GoldController) CreateAddress(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	addr := &utils.AugmontUserAddressInfo{}
	if err := ctx.Bind(addr); err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	err = c.gold.CreateUserAddress(agUser, addr)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":        "ok",
		"userAddressID": addr.UserAddressID,
	})
}

func (c *GoldController) GetUserAddress(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	addr, err := c.gold.GetUserAddresses(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"ststus":  "ok",
		"address": addr,
	})
}

func (c *GoldController) DeleteuserAddress(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	userAddressID := ctx.Param("userAddressID")
	if userAddressID == "" {
		err = errors.New("userAddressID not vaild")
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	addr := &utils.AugmontUserAddressInfo{
		UserAddressID: userAddressID,
	}
	err = c.gold.DeleteUserAddress(agUser, addr)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"ststus":  "ok",
		"message": "user address deleted",
	})
}

func (c *GoldController) CreateKYC(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	kyc := &struct {
		Name  string `form:"nameAsPerPan"`
		PanNo string `form:"panNumber"`
		DOB   string `form:"dateOfBirth"`
	}{}
	if err := ctx.Bind(kyc); err != nil {
		return
	}
	file, _ := ctx.FormFile("panAttachment")
	format := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	localFile := &utils.File{
		Key:         uuid.NewString(),
		ContentType: file.Header.Get("Content-Type"),
		Format:      format,
	}
	ctx.SaveUploadedFile(file, localFile.Path())
	defer localFile.Close()

	data, err := c.gold.PostUserKyc(
		kyc.Name, kyc.PanNo, kyc.DOB,
		agUser,
		localFile,
	)

	if err != nil {
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "ok",
		"message": data,
	})
}

func (c *GoldController) GetKycStatus(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}
	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	if agUser.KYCStatus != nil && *agUser.KYCStatus != "pending" {
		ctx.JSON(200, gin.H{
			"status":    "ok",
			"kycStatus": agUser.KYCStatus,
		})
	}

	err = c.gold.UpdateUserKycStatus(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status":    "ok",
		"kycStatus": agUser.KYCStatus,
	})
}

func (c *GoldController) BuyOrder(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	info := &utils.AugmontBugInfo{}
	if err := ctx.Bind(info); err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	data, err := c.gold.Buy(agUser, info)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"order":  data,
	})
}

func (c *GoldController) GetBuyInfo(ctx *gin.Context) {

	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	txnID := ctx.Param("txnID")

	data, err := c.gold.BuyInfo(*agUser.UID, txnID)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"order":  data,
	})
}

func (c *GoldController) GetBuyList(ctx *gin.Context) {

	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	data, err := c.gold.BuyList(*agUser.UID)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"order":  data,
	})
}

func (c *GoldController) SellOrder(ctx *gin.Context) {
	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	info := &utils.AugmontSellInfo{}
	if err := ctx.Bind(info); err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	data, err := c.gold.Sell(agUser, info)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"order":  data,
	})
}

func (c *GoldController) GetSellInfo(ctx *gin.Context) {

	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	txnID := ctx.Param("txnID")

	data, err := c.gold.SellInfo(*agUser.UID, txnID)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"order":  data,
	})
}

func (c *GoldController) GetSellList(ctx *gin.Context) {

	user, err := getPinchUserFromContext(ctx)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	agUser := &models.AugmontUser{
		UserID: user.ID,
	}

	agUser, err = c.augmontUser.FindUser(agUser)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	data, err := c.gold.SellList(*agUser.UID)
	if err != nil {
		domain.ErrLog(err)
		domain.ErrFailedGinReq(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
		"order":  data,
	})
}
