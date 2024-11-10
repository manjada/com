package dto

const (
	DateFormatMediumDash     = "02-01-2006"
	OtpChar                  = "1234567890"
	STATUS_USER_VERIFICATION = "Verification"
	STATUS_USER_ACTIVE       = "Active"

	OTP_REGISTRATION_SUCCESS = "Registrasi Sukses silahkan login kembali"
	OTP_VERIFICATION_MAX     = 3

	ROLES_CODE_CLIENT = "CLT"
	ROLES_CODE_ADMIN  = "ADM"

	DataSuccessDeleted = "Data successfully deleted"

	// AUDIT LOG ACTION
	AUDIT_ACTION_CREATE  = "CREATE"
	AUDIT_ACTION_UPDATE  = "UPDATE"
	AUDIT_ACTION_DELETE  = "DELETE"
	AUDIT_ACTION_APPROVE = "APPROVED"
	AUDIT_ACTION_REJECT  = "REJECTED"
	AUDIT_ACTION_READ    = "READ"
)
