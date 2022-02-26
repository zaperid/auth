package user

func processUser(data *Data) {
	if data.Password != "" {
		data.Password = hash([]byte(data.Password))
	}
}
