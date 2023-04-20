package model

import (
	"time"

	"gorm.io/gorm"
)

type SimpleOption struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ConfirmRequest struct {
	Password string `json:"password"`
	UserId   int64  `json:"-"`
}

type Bank struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	IconUrl   string    `json:"iconUrl"`
	TypeFlag  string    `json:"typeFlag"`
	CreatedAt time.Time `json:"createdAt"`
}
type BankResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IconUrl  string `json:"iconUrl"`
	TypeFlag string `json:"typeFlag"`
}
type BankListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}
type BankListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type AccountType struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TypeFlag  string    `json:"typeFlag"`
	CreatedAt time.Time `json:"createdAt"`
}
type AccountTypeResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
type AccountTypeListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}
type AccountTypeListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type BankAccount struct {
	Id                    int64          `json:"id"`
	BankId                int64          `json:"bankId"`
	BankCode              string         `json:"bankCode"`
	BankName              string         `json:"bankName"`
	BankIconUrl           string         `json:"bankIconUrl"`
	AccountTypeId         int64          `json:"accountTypeId"`
	AccountTypeName       string         `json:"accountTypeName"`
	AccountName           string         `json:"accountName"`
	AccountNumber         string         `json:"accountNumber"`
	AccountBalance        float64        `json:"accountBalance" sql:"type:decimal(14,2);"`
	AccountPriority       string         `json:"accountPriority"`
	AccountStatus         string         `json:"accountStatus"`
	DeviceUid             string         `json:"deviceUid"`
	PinCode               string         `json:"pinCode"`
	ConnectionStatus      string         `json:"connectionStatus"`
	ExternalId            string         `json:"-"`
	LastConnUpdateAt      *time.Time     `json:"lastConnUpdateAt"`
	AutoCreditFlag        string         `json:"autoCreditFlag"`
	AutoWithdrawFlag      string         `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount string         `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount string         `json:"autoTransferMaxAmount"`
	QrWalletStatus        string         `json:"qrWalletStatus"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             *time.Time     `json:"updatedAt"`
	DeletedAt             gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountListRequest struct {
	AccountNumber string `form:"accountNumber"`
	AccountType   string `form:"accountType"`
	Page          int    `form:"page" default:"1" min:"1"`
	Limit         int    `form:"limit" default:"10" min:"1" max:"100"`
	Search        string `form:"search"`
	SortCol       string `form:"sortCol"`
	SortAsc       string `form:"sortAsc"`
}

type BankAccountCreateBody struct {
	BankId                int64   `json:"bankId" validate:"required"`
	AccountTypeId         int64   `json:"accounTypeId" validate:"required"`
	AccountName           string  `json:"accountName" validate:"required"`
	AccountNumber         string  `json:"accountNumber" validate:"required"`
	AccountBalance        float64 `json:"-"`
	DeviceUid             string  `json:"deviceUid"`
	PinCode               string  `json:"pinCode"`
	AutoCreditFlag        string  `json:"autoCreditFlag"`
	AutoWithdrawFlag      string  `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount string  `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount string  `json:"autoTransferMaxAmount"`
	AccountPriority       string  `json:"accountPriority"`
	QrWalletStatus        string  `json:"qrWalletStatus"`
	AccountStatus         string  `json:"accountStatus"`
	ConnectionStatus      string  `json:"-"`
}

type BankAccountUpdateRequest struct {
	BankId                *int64  `json:"-"`
	AccountTypeId         *int64  `json:"accounTypeId"`
	AccountName           *string `json:"-"`
	AccountNumber         *string `json:"-"`
	DeviceUid             *string `json:"deviceUid"`
	PinCode               *string `json:"pinCode"`
	AutoCreditFlag        *string `json:"autoCreditFlag"`
	AutoWithdrawFlag      *string `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount *string `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount *string `json:"autoTransferMaxAmount"`
	AccountPriority       *string `json:"accountPriority"`
	QrWalletStatus        *string `json:"qrWalletStatus"`
	AccountStatus         *string `json:"accountStatus"`
}

type BankAccountUpdateBody struct {
	BankId                *int64     `json:"-"`
	AccountTypeId         *int64     `json:"accounTypeId"`
	AccountName           *string    `json:"-"`
	AccountNumber         *string    `json:"-"`
	DeviceUid             *string    `json:"deviceUid"`
	PinCode               *string    `json:"pinCode"`
	ExternalId            *int64     `json:"-"`
	AutoCreditFlag        *string    `json:"autoCreditFlag"`
	AutoWithdrawFlag      *string    `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount *string    `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount *string    `json:"autoTransferMaxAmount"`
	AccountPriority       *string    `json:"accountPriority"`
	QrWalletStatus        *string    `json:"qrWalletStatus"`
	AccountStatus         *string    `json:"accountStatus"`
	LastConnUpdateAt      *time.Time `json:"-"`
	ConnectionStatus      *string    `json:"-"`
	AccountBalance        *float64   `json:"-"`
}

type BankAccountDeleteBody struct {
	AccountNumber string    `json:"-"`
	DeletedAt     time.Time `json:"-"`
}

type BankAccountResponse struct {
	Id               int64          `json:"id"`
	BankId           int64          `json:"bankId"`
	BankCode         string         `json:"bankCode"`
	BankName         string         `json:"bankName"`
	BankIconUrl      string         `json:"bankIconUrl"`
	AccountTypeId    int64          `json:"accountTypeId"`
	AccountTypeName  string         `json:"accountTypeName"`
	AccountName      string         `json:"accountName"`
	AccountNumber    string         `json:"accountNumber"`
	AccountBalance   float64        `json:"accountBalance"`
	AccountPriority  string         `json:"accountPriority"`
	AccountStatus    string         `json:"accountStatus"`
	ConnectionStatus string         `json:"connectionStatus"`
	LastConnUpdateAt *time.Time     `json:"lastConnUpdateAt"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        *time.Time     `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransaction struct {
	Id                int64          `json:"id"`
	AccountId         int64          `json:"accountId"`
	Description       string         `json:"description"`
	TransferType      string         `json:"transferType"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransactionParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountTransactionListRequest struct {
	AccountId       int    `form:"accountId"`
	FromCreatedDate string `form:"fromCreatedDate"`
	ToCreatedDate   string `form:"toCreatedDate"`
	TransferType    string `form:"transferType"`
	Search          string `form:"search"`
	Page            int    `form:"page" default:"1" min:"1"`
	Limit           int    `form:"limit" default:"10" min:"1" max:"100"`
	SortCol         string `form:"sortCol"`
	SortAsc         string `form:"sortAsc"`
}

type BankAccountTransactionBody struct {
	AccountId         int64     `json:"accountId" validate:"required"`
	Description       string    `json:"description"`
	TransferType      string    `json:"transferType" validate:"required"`
	Amount            float32   `json:"amount" validate:"required"`
	TransferAt        time.Time `json:"transferAt" validate:"required"`
	CreatedByUsername string    `json:"-"`
}

type BankAccountTransactionResponse struct {
	Id                int64          `json:"id"`
	AccountId         int64          `json:"accountId"`
	BankName          string         `json:"bankName"`
	AccountName       string         `json:"accountName"`
	AccountNumber     string         `json:"accountNumber"`
	Description       string         `json:"description"`
	TransferType      string         `json:"transferType"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransfer struct {
	Id                int64          `json:"id"`
	FromAccountId     int64          `json:"fromAccountId"`
	FromBankId        int64          `json:"fromBankId"`
	FromBankName      string         `json:"fromBankName"`
	FromAccountName   string         `json:"fromAccountName"`
	FromAccountNumber string         `json:"fromAccountNumber"`
	ToAccountId       int64          `json:"toAccountId"`
	ToBankId          int64          `json:"toBankId"`
	ToBankName        string         `json:"toBankName"`
	ToAccountName     string         `json:"toAccountName"`
	ToAccountNumber   string         `json:"toAccountNumber"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	Status            string         `json:"status"`
	ConfirmedAt       time.Time      `json:"confirmedAt"`
	ConfirmedByUserId int64          `json:"confirmedByUserId"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransferParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountTransferListRequest struct {
	AccountId       int    `form:"accountId"`
	FromCreatedDate string `form:"fromCreatedDate"`
	ToCreatedDate   string `form:"toCreatedDate"`
	ToAccountId     int    `form:"toAccountId"`
	Search          string `form:"search"`
	Page            int    `form:"page" default:"1" min:"1"`
	Limit           int    `form:"limit" default:"10" min:"1" max:"100"`
	SortCol         string `form:"sortCol"`
	SortAsc         string `form:"sortAsc"`
}

type BankAccountTransferBody struct {
	Status            string    `json:"-"`
	FromAccountId     int64     `json:"fromAccountId" validate:"required"`
	FromBankId        int64     `json:"-"`
	FromAccountName   string    `json:"-"`
	FromAccountNumber string    `json:"-"`
	ToAccountId       int64     `json:"toAccountId" validate:"required"`
	ToBankId          int64     `json:"-"`
	ToAccountName     string    `json:"-"`
	ToAccountNumber   string    `json:"-"`
	Amount            float32   `json:"amount" validate:"required"`
	TransferAt        time.Time `json:"transferAt" validate:"required"`
	CreatedByUsername string    `json:"-"`
}

type BankAccountTransferConfirmBody struct {
	Status            string    `json:"status" validate:"required"`
	ConfirmedByUserId int64     `json:"confirmedByUserId" validate:"required"`
	ConfirmedAt       time.Time `json:"confirmedAt" validate:"required"`
}

type BankAccountTransferResponse struct {
	Id                int64          `json:"id"`
	FromAccountId     int64          `json:"fromAccountId"`
	FromBankId        int64          `json:"fromBankId"`
	FromBankName      string         `json:"fromBankName"`
	FromAccountName   string         `json:"fromAccountName"`
	FomAccountNumber  string         `json:"fromAccountNumber"`
	ToAccountId       int64          `json:"toAccountId"`
	ToBankId          int64          `json:"toBankId"`
	ToBankName        string         `json:"toBankName"`
	ToAccountName     string         `json:"toAccountName"`
	ToAccountNumber   string         `json:"toAccountNumber"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	Status            string         `json:"status"`
	ConfirmedAt       time.Time      `json:"confirmedAt"`
	ConfirmedByUserId int64          `json:"confirmedByUserId"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type ExternalAccount struct {
	BankId           int64   `json:"bankId"`
	BankCode         string  `json:"bankCode"`
	ClientName       string  `json:"clientName"`
	LastConnected    *int64  `json:"lastConnected"`
	CustomerId       int64   `json:"customerId"`
	DeviceId         string  `json:"deviceId"`
	WebhookUrl       *string `json:"webhookUrl"`
	WalletId         *int64  `json:"walletId"`
	Enable           bool    `json:"enable"`
	AccountNo        string  `json:"accountNo"`
	BankAccountId    *int64  `json:"bankAccountId"`
	VerifyLogin      bool    `json:"verifyLogin"`
	WebhookNotifyUrl *string `json:"webhookNotifyUrl"`
	Username         *string `json:"username"`
}

type ExternalAccountStatusRequest struct {
	AccountNumber string `json:"accountNumber"`
}
type ExternalAccountEnableRequest struct {
	AccountNo string `json:"accountNo"`
	Enable    bool   `json:"enable"`
}

type ExternalAccountBalance struct {
	LimitUsed            float32 `json:"limitUsed"`
	BranchId             string  `json:"branchId"`
	AccountName          string  `json:"accountName"`
	DailyLimitOtherBanks float32 `json:"dailyLimitOtherBanks"`
	DailyLimitPromptPay  float32 `json:"dailyLimitPromptPay"`
	AccruedInterest      float32 `json:"accruedInterest"`
	OverdraftLimit       float32 `json:"overdraftLimit"`
	DailyLimitSCBOther   float32 `json:"dailyLimitSCBOther"`
	DailyLimitSCBOwn     float32 `json:"dailyLimitSCBOwn"`
	AvailableBalance     string  `json:"availableBalance"`
	AccountNo            string  `json:"accountNo"`
	Currency             string  `json:"currency"`
	AccountBalance       string  `json:"accountBalance"`
	Status               struct {
		Code        int    `json:"code"`
		Header      string `json:"header"`
		Description string `json:"description"`
	} `json:"status"`
}
type ExternalAccountStatus struct {
	Success bool   `json:"success"`
	Enable  bool   `json:"enable"`
	Status  string `json:"status"`
}

type ExternalAccountCreateBody struct {
	AccountNo        string `json:"accountNo"`
	BankCode         string `json:"bankCode"`
	DeviceId         string `json:"deviceId"`
	Password         string `json:"password"`
	Pin              string `json:"pin"`
	Username         string `json:"username"`
	WebhookNotifyUrl string `json:"webhookNotifyUrl"`
	WebhookUrl       string `json:"webhookUrl"`
}

type ExternalAccountUpdateBody struct {
	AccountNo        string  `json:"accountNo"`
	BankCode         string  `json:"bankCode"`
	DeviceId         *string `json:"deviceId"`
	Password         string  `json:"password"`
	Pin              *string `json:"pin"`
	Username         string  `json:"username"`
	WebhookNotifyUrl string  `json:"webhookNotifyUrl"`
	WebhookUrl       string  `json:"webhookUrl"`
}

type ExternalAccountCreateResponse struct {
	Id               int64  `json:"id"`
	CustomerId       int64  `json:"customerId"`
	ApiKey           string `json:"apiKey"`
	BankId           int64  `json:"bankId"`
	BankCode         string `json:"bankCode"`
	DeviceId         string `json:"deviceId"`
	AccountNo        string `json:"accountNo"`
	Pin              string `json:"pin"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	WebhookUrl       string `json:"webhookUrl"`
	WebhookNotifyUrl string `json:"webhookNotifyUrl"`
	WalletId         int64  `json:"walletId"`
	Enable           bool   `json:"enable"`
	VerifyLogin      bool   `json:"verifyLogin"`
	Deleted          bool   `json:"deleted"`
}

type ExternalListWithPagination struct {
	Content       interface{} `json:"content"`
	TotalElements int64       `json:"totalElements"`
}

type ExternalAccountLog struct {
	Id                 int64          `json:"id"`
	ExternalId         int64          `json:"externalId"`
	ClientName         string         `json:"clientName"`
	LogType            string         `json:"logType"`
	Message            string         `json:"message"`
	ExternalCreateDate string         `json:"externalCreateDate"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          *time.Time     `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt"`
}

type ExternalAccountLogCreateBody struct {
	ExternalId         int64  `json:"externalId"`
	LogType            string `json:"logType"`
	Message            string `json:"message"`
	ExternalCreateDate string `json:"externalCreateDate"`
}

type ExternalAccountStatement struct {
	Id                 int64          `json:"id"`
	ExternalId         int64          `json:"externalId"`
	BankAccountId      int64          `json:"bankAccountId"`
	BankCode           string         `json:"bankCode"`
	Amount             float32        `json:"amount"`
	DateTime           time.Time      `json:"dateTime"`
	RawDateTime        time.Time      `json:"rawDateTime"`
	Info               string         `json:"info"`
	ChannelCode        string         `json:"channelCode"`
	ChannelDescription string         `json:"channelDescription"`
	TxnCode            string         `json:"txnCode"`
	TxnDescription     string         `json:"txnDescription"`
	Checksum           string         `json:"checksum"`
	IsRead             bool           `json:"isRead"`
	ExternalCreateDate string         `json:"externalCreateDate"`
	ExternalUpdateDate string         `json:"externalUpdateDate"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          *time.Time     `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt"`
}

type ExternalAccountStatementCreateBody struct {
	ExternalId         int64   `json:"externalId"`
	BankAccountId      int64   `json:"bankAccountId"`
	BankCode           string  `json:"bankCode"`
	Amount             float32 `json:"amount"`
	DateTime           string  `json:"dateTime"`
	RawDateTime        string  `json:"rawDateTime"`
	Info               string  `json:"info"`
	ChannelCode        string  `json:"channelCode"`
	ChannelDescription string  `json:"channelDescription"`
	TxnCode            string  `json:"txnCode"`
	TxnDescription     string  `json:"txnDescription"`
	Checksum           string  `json:"checksum"`
	IsRead             bool    `json:"isRead"`
	ExternalCreateDate string  `json:"externalCreateDate"`
	ExternalUpdateDate string  `json:"externalUpdateDate"`
}

type ExternalAccountTransferRequest struct {
	SystemAccountId int64  `json:"systemAccountId" validate:"required"`
	AccountNumber   string `json:"accountNumber" validate:"required"`
	BankCode        string `json:"bankCode" validate:"required"`
	Amount          string `json:"amount" validate:"required"`
}

type ExternalAccountTransferBody struct {
	AccountForm string `json:"accountFrom"`
	AccountTo   string `json:"accountTo"`
	Amount      string `json:"amount"`
	BankCode    string `json:"bankCode"`
	Pin         string `json:"pin"`
}

type ExternalAccountError struct {
	Timestamp int64  `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Path      string `json:"path"`
}

type WebhookLog struct {
	Id          int64          `json:"id"`
	JsonRequest string         `json:"jsonRequest"`
	JsonPayload string         `json:"jsonPayload"`
	LogType     string         `json:"logType"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}

type WebhookLogCreateBody struct {
	JsonRequest string `json:"jsonRequest"`
	JsonPayload string `json:"jsonPayload"`
	LogType     string `json:"logType"`
	Status      string `json:"status"`
}

type WebhookStatementResponse struct {
	NewStatementList []WebhookStatement `json:"newStatementList"`
}

type WebhookStatement struct {
	Id                 int64     `json:"id"`
	CustomerId         int64     `json:"customerId"`
	ClientId           int64     `json:"clientId"`
	ClientName         string    `json:"clientName"`
	BankAccountId      int64     `json:"bankAccountId"`
	BankCode           string    `json:"bankCode"`
	Amount             float32   `json:"amount"`
	DateTime           time.Time `json:"dateTime"`
	RawDateTime        time.Time `json:"rawDateTime"`
	Info               string    `json:"info"`
	ChannelCode        string    `json:"channelCode"`
	ChannelDescription string    `json:"channelDescription"`
	TxnCode            string    `json:"txnCode"`
	TxnDescription     string    `json:"txnDescription"`
	Checksum           string    `json:"checksum"`
	IsRead             bool      `json:"isRead"`
	CreatedDate        string    `json:"createdDate"`
	UpdatedDate        string    `json:"updatedDate"`
}

type BotAccountConfig struct {
	Id        int64  `json:"id"`
	ConfigKey string `json:"configKey"`
	ConfigVal string `json:"configVal"`
}

type BotAccountConfigListRequest struct {
	SearchKey   *string `form:"searchKey"`
	SearchValue *string `form:"searchValue"`
	Page        int     `form:"page" default:"1" min:"1"`
	Limit       int     `form:"limit" default:"10" min:"1" max:"100"`
	SortCol     string  `form:"sortCol"`
	SortAsc     string  `form:"sortAsc"`
}

type BotAccountConfigCreateBody struct {
	ConfigKey string `json:"configKey"`
	ConfigVal string `json:"configVal"`
}
