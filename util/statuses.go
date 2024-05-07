package util

const (
	OK                          = "OK"
	Unauthorized                = "UnauthorizedUser"
	NoAuthorization             = "NoAuthorization"
	NotFound                    = "NotFound"
	InternalError               = "InternalError"
	InvalidID                   = "InvalidID"
	UnsupportedAuthorization    = "UnsupportedAuthorization"
	InvalidAuthorizationFormat  = "InvalidAuthorizationFormat"
	ExpiredToken                = "ExpiredToken"
	InvalidCurrency             = "InvalidCurrency"
	InvalidPageID               = "InvalidPageID"
	InvalidPageSize             = "InvalidPageSize"
	DuplicateUsername           = "DuplicateUsername"
	InvalidUsername             = "InvalidUsername"
	InvalidEmail                = "InvalidEmail"
	TooShortPassword            = "TooShortPassword"
	UserNotFound                = "UserNotFound"
	IncorrectPassword           = "IncorrectPassword"
	UnauthorizedUser            = "UnauthorizedUser"
	FromAccountNotFound         = "FromAccountNotFound"
	ToAccountNotFound           = "ToAccountNotFound"
	FromAccountCurrencyMismatch = "FromAccountCurrencyMismatch"
	ToAccountCurrencyMismatch   = "ToAccountCurrencyMismatch"
	NegativeAmount              = "NegativeAmount"
	GetAccountError             = "GetAccountError"
	TransferTxError             = "TransferTxError"

	OtherDepositorCannotUpdateThisUserInfo = "OtherDepositorCannotUpdateThisUserInfo"
	BankerCanUpdateUserInfo                = "BankerCanUpdateUserInfo"
)