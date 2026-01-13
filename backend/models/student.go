package models

type Student struct {
	ID           string `json:"id"`
	StudentName  string `json:"studentName"`
	Address      string `json:"address"`
	State        string `json:"state"`
	District     string `json:"district"`
	Taluka       string `json:"taluka"`
	Gender       string `json:"gender"`
	DOB          string `json:"dob"`
	Photo        string `json:"photo"`
	Handicapped  bool   `json:"handicapped"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
	BloodGroup   string `json:"bloodGroup"`
}