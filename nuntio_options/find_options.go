package nuntio_options

type FindOptions struct {
	Id         string
	OptionalId string
	Email      string
}

func (fo *FindOptions) Validate() bool {
	if fo == nil || (fo.OptionalId == "" && fo.Id == "" && fo.Email == "") {
		return false
	}
	return true
}
