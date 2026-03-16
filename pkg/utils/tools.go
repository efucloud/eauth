package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ArrayItemHasPrefix(key string, arrays []string) bool {
	for _, item := range arrays {
		if strings.HasPrefix(key, item) || key == item {
			return true
		}
	}
	return false

}
func RandNumber(length int) (result string) {
	rand.Int63n(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		result += fmt.Sprintf("%d", rand.Intn(10))
	}
	return result
}
func GetRemoteAddress(address string) string {
	return strings.TrimPrefix(strings.TrimSuffix(strings.Split(address, ":")[0], "]"), "[")
}
func ComparePassword(passwordHash string, password, salt string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password+salt))
	if err != nil {
		err = errors.New("username or password is not right")
	}
	return err

}
func GetBrowserDetail(useAgent string) (browser, platformDetail string) {
	re := regexp.MustCompile(`\((.*?)\)`)
	match := re.FindStringSubmatch(useAgent)
	if len(match) > 1 {
		platformDetail = match[1]
	}
	if (strings.Index(useAgent, "Firefox")) > -1 {
		browser = "Firefox"
	} else if (strings.Index(useAgent, "QQBrowser")) > -1 {
		browser = "QQBrowser"
	} else if (strings.Index(useAgent, "QQ")) > -1 {
		browser = "QQ"
	} else if (strings.Index(useAgent, "UCBrowser")) > -1 {
		browser = "UCBrowser"
	} else if (strings.Index(useAgent, "Opera")) > -1 || (strings.Index(useAgent, "OPR")) > -1 {
		browser = "Opera"
	} else if (strings.Index(useAgent, "Wechat")) > -1 {
		browser = "Wechat"
	} else if (strings.Index(useAgent, "Trident")) > -1 {
		browser = "InternetExplorer"
	} else if (strings.Index(useAgent, "Edge")) > -1 {
		browser = "Edge"
	} else if (strings.Index(useAgent, "Chrome")) > -1 {
		browser = "Chrome"
	} else if (strings.Index(useAgent, "Safari")) > -1 {
		browser = "Safari"
	} else {
		browser = "Unknown"
	}
	return

}
