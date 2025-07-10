package utils

func FormatPrintableMessage(message string, maxLength int) string {
	if len(message) <= maxLength {
		return message
	}

	return message[:maxLength] + "..."
}