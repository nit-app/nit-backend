package sessions

type OtpState struct {
	Code        string
	Attempt     int
	PhoneNumber string
}
