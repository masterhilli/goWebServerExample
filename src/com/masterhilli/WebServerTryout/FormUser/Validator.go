package FormUser

import (
	"regexp"
	"log"
	"os"
)

var validateLogger *log.Logger = log.New(os.Stdout, "DEBUG: ",  log.Ldate|log.Ltime|log.Lshortfile)

type Validator interface {
	Validate(key, value string) bool
}

type RequiredFieldValidator struct {
}

func (this RequiredFieldValidator) Validate(key, value string) bool {
	retVal := len(value) > 0
	if  !retVal {
		validateLogger.Println("The value in " + key + " is required!")
	}
	return retVal
}

type NumberFieldValidator struct {
}

func (this NumberFieldValidator) Validate(key, value string) bool {
	retVal, _ := regexp.MatchString("^[0-9]+$", value)
	if !retVal {
		validateLogger.Println("Only numbers are allowed for field \"" + key + "\"")
	}
	return retVal
}

type ChineseFieldValidator struct {
}

func (this ChineseFieldValidator) Validate(key, value string) bool {
	retVal, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+$", value)
	if !retVal {
		validateLogger.Println("Only Chinese characters are allowed for field \"" + key + "\".")
	}
	return retVal
}

