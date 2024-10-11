package user

// RandomUserResponse represents struct of randomuser.me API response.
type RandomUserResponse struct {
	Results []RandomUser `json:"results"`
}

// RandomUser represents singular generated user response data.
type RandomUser struct {
	RandomUserName RandomUserName `json:"name"`
	Email          string         `json:"email"`
	RandomUserDob  RandomUserDob  `json:"dob"`
}

// RandomUserName represents name and surname of generated user.
type RandomUserName struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// RandomUserDob represents date of birth information about user.
type RandomUserDob struct {
	Age int `json:"age"`
}
