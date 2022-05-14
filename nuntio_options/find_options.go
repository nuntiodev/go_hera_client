package nuntio_options

type FindOptions struct {
	Id       string
	Username string
	Email    string
}

func (fo *FindOptions) Validate() bool {
	if fo == nil || (fo.Username == "" && fo.Id == "" && fo.Email == "") {
		return false
	}
	return true
}
