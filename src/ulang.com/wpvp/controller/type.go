package controller

type Speaker struct {
	ID        int64  "ID"
	User      string "User"
	Device    string "Device"
	NickName  string "NickName"
	Email     string "Email"
	Gender    string "Gender"
	TrainType int    "TrainType"
	AppType   string "AppType"
	Lang      string "Lang"
	Online    bool   "Online"
	State     int    "State"
	Threshold int    "Threshold"
	Notes     string "Notes"
	Usage     int    "Usage"
	CreatedAt int64  "CreatedAt"
	UpdatedAt int64  "UpdatedAt"
}

type UserForm struct {
	Username    string `json:"Username"form:"Username"`
	PhoneNumber string `json:"PhoneNumber"form:"PhoneNumber"`
	Email       string `json:"Email"form:"Email"`
	Avatar      string `json:"Avatar"form:"Avatar"`
	Password    string `json:"Password"form:"Password"`
}

type ForwardRequest struct {
	Cond     map[string]interface{}
	MaxCount int64
}

type ForwardResult struct {
	code  string "code"
	msg   interface{}
	count int
}
