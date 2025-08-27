package toolkit

import "github.com/google/uuid"

// CreateUUID 创建UUID
func CreateUUID() (string, error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuidStr := uuidObj.String()
	return uuidStr, nil

}
