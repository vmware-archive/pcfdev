package usecases

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"fmt"
	"errors"
	"strings"
)

type UaaCredentialReplacement struct{}

func (u *UaaCredentialReplacement) ReplaceUaaConfigAdminCredentials(uaaConfig string, password string) (_ string, err error) {
	users, _ := findArray(uaaConfig, "scim.users")

	adminCredentialPattern := regexp.MustCompile(`admin\|admin\|`)
	for index, userSettingInterface := range users {
		userSetting := userSettingInterface.(string)
		if adminCredentialPattern.MatchString(userSetting) {
			securedUserSetting := adminCredentialPattern.ReplaceAllString(userSetting, fmt.Sprintf("admin|%s|", password))
			return setAt(uaaConfig, "scim.users", index, securedUserSetting)
		}
	}

	return "", errors.New("failed to parse UAA config file")
}




type interfaceMap map[interface{}]interface{}

func recoverFromInterfaceConversion(recovered interface{}) error {
	if recovered != nil {
		errorMessage := recovered.(error).Error()
		if strings.Contains(errorMessage, "interface conversion")  {
			return errors.New("failed to parse yaml")
		} else {
			panic(recovered)
		}
	}
	return nil
}

func setAt(contents string, path string, arrayIndex int, value string) (string, error) {
	yamlContents := make(interfaceMap)

	if err := yaml.Unmarshal([]byte(contents), &yamlContents); err != nil {
		return "", errors.New("failed to parse yaml")
	}

	var currentNode interface{} = yamlContents
	keys := strings.Split(path, ".")
	for _, key := range keys {
		currentNode = currentNode.(interfaceMap)[key]
	}
	currentNode.([]interface{})[arrayIndex] = value
	replacedContents, err := yaml.Marshal(yamlContents)
	return string(replacedContents), err
}

func findArray(contents string, path string) (_ []interface{}, err error) {
	defer func() {
		err = recoverFromInterfaceConversion(recover())
	}()

	yamlContents := make(interfaceMap)

	if err := yaml.Unmarshal([]byte(contents), &yamlContents); err != nil {
		return nil, errors.New("failed to parse yaml")
	}

	var currentNode interface{} = yamlContents
	keys := strings.Split(path, ".")
	for _, key := range keys {
		currentNode = currentNode.(interfaceMap)[key]
	}

	return currentNode.([]interface{}), nil
}


