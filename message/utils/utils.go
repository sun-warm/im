package utils

func GenerateConversationID(userName, receiver string) string {
	if userName < receiver {
		return userName + receiver
	}
	return receiver + userName
}
