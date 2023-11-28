package requests

type PhoneNumberRequest struct {
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" binding:"required,min=10,max=12,ruPhoneNumber"`
}

type OtpCheckRequest struct {
	Code string `form:"code" json:"code" binding:"required,min=4,max=6"`
}

type FinishRegistrationRequest struct {
	FirstName string  `form:"firstName" json:"firstName" binding:"required"`
	LastName  *string `form:"lastName" json:"lastName"`
}
