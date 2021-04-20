package locale

func initMessages() map[string]map[string]string {

	messages := make(map[string]map[string]string)

	//Unknown error
	messages["unknown_error"] = make(map[string]string)
	messages["unknown_error"]["ru"] = "Неизветсная ошибка"
	messages["unknown_error"]["en"] = "Unknown error"

	//Unknown command
	messages["unknown_command"] = make(map[string]string)
	messages["unknown_command"]["ru"] = "Неизветсная команда"
	messages["unknown_command"]["en"] = "Unknown command"

	return messages
}
