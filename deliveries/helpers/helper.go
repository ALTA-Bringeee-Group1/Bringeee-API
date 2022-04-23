package helpers

func ResponseFailed(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "failed",
		"message": message,
	}
}
