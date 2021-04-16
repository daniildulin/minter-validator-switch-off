package locale

func initErrors() map[string]map[string]string {

	errors := make(map[string]map[string]string)

	errors["INTERNAL_ERROR"] = make(map[string]string)
	errors["INTERNAL_ERROR"]["ru"] = "Ой, что то пошло не так"
	errors["INTERNAL_ERROR"]["en"] = "Oops, something wrong"

	return errors
}
