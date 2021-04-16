package locale

type Locale struct {
	messages map[string]map[string]string
	labels   map[string]map[string]string
	errors   map[string]map[string]string
}

func NewBotLocale() *Locale {
	return &Locale{
		messages: initMessages(),
		errors:   initErrors(),
		labels:   initLabels(),
	}
}

func (l Locale) GetMessage(msg, lang string, data ...interface{}) string {

	if l.messages[msg] == nil || l.messages[msg][lang] == "" {
		return msg
	}

	return l.messages[msg][lang]
}
func (l Locale) GetLabel(label, lang string) string {

	if l.labels[label] == nil || l.labels[label][lang] == "" {
		return label
	}

	return l.labels[label][lang]
}

func (l Locale) GetError(label, lang string) string {

	if l.errors[label] == nil || l.errors[label][lang] == "" {
		return label
	}

	return l.errors[label][lang]
}
