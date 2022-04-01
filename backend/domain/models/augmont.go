package models

import "time"

// Augment User Model
type AugmontUser struct {
	ID        *uint64    `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Unique for each Augmont User
	UID *string `json:"uid" gorm:"type:varchar(30); not null; unique; <-:create"`

	// Augmont User KYC status
	// NULL -> KYC not done,  pending -> KYC pending
	// approved -> KYC approved, rejected -> KYC rejected
	KYCStatus *string `json:"status" gorm:"type:augmont_kyc_status"`

	// User Table Relation
	// User Can Have only one Augmont User
	UserID *uint64 `json:"userId" gorm:"not null; unique"`
	User   *User   `json:"user" gorm:"foreignkey:UserID"`

	// Fileds won't be created while migration, used only for reading relations
	Banks   []*AugmontUserBank    `json:"banks" gorm:"->"`
	Address []*AugmontUserAddress `json:"address" gorm:"->"`
}

// Augment User Bank Model
// atmost 10 banks per Augmont User
type AugmontUserBank struct {
	ID        *uint64    `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Unique for each bank entry
	UserBankID    *string `json:"userBankID" gotm:"not null"`
	AugmontUserID *uint64 `json:"goldUserID" gorm:"not null"`

	// Relations
	AugmontUser *AugmontUser `json:"goldUser" gorm:"foreignkey:AugmontUserID"`
}

// Augmont User Address Model
type AugmontUserAddress struct {
	ID        *uint64    `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Unique  for each bank entry
	UserAddressID *string `json:"userAddressID" gorm:"not null"`

	AugmontUserID *uint64 `json:"goldUserID" gorm:"not null"`

	// Relations
	AugmontUser *AugmontUser `json:"goldUser" gorm:"foreignkey:AugmontUserID"`
}

type AugmontRedeemOrder struct {
	ID        *uint64    `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Unique for each redeeme order entry
	MerchantTxnID *string `json:"merchantTxnID" gorm:"not null; unique"`
	AugmontUserID *uint64 `json:"goldUserID" gorm:"not null"`

	// Relations
	AugmontUser *AugmontUser `json:"goldUser" gorm:"foreignkey:AugmontUserID"`
}

type AugmontBuyOrder struct {
	ID        *uint64    `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Unique for each redeeme order entry
	MerchantTxnID *string `json:"merchantTxnID" gorm:"not null; unique"`
	AugmontUserID *uint64 `json:"goldUserID" gorm:"not null"`

	// Relations
	AugmontUser *AugmontUser `json:"goldUser" gorm:"foreignkey:AugmontUserID"`
}

type AugmontSellOrder struct {
	ID        *uint64    `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Unique for each redeeme order entry
	MerchantTxnID *string `json:"merchantTxnID" gorm:"not null; unique"`
	AugmontUserID *uint64 `json:"goldUserID" gorm:"not null"`

	// Relations
	AugmontUser *AugmontUser `json:"goldUser" gorm:"foreignkey:AugmontUserID"`
}
