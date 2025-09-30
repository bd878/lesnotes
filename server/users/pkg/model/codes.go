package model

var (
	CodePasswordTooShort   int = 101
	CodePasswordUpperLower int = 102
	CodePasswordOneNumber  int = 103
	CodePasswordOneSpecial int = 104
	CodeLoginTooShort      int = 105

	CodeRegisterFailed     int = 111
	// 112 - free to use 
	CodeLogoutFailed       int = 113
	CodeDeleteFailed       int = 114
	CodeUpdateFailed       int = 115

	CodeNoLogin            int = 121
	CodeNoPassword         int = 122
	CodeBadCookie          int = 124
)
