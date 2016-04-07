package Webserver

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var validateLogger *log.Logger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

type Validator interface {
	Validate(key, value string) bool
}

type RequiredFieldValidator struct {
}

func (this RequiredFieldValidator) Validate(key, value string) bool {
	retVal := len(value) > 0
	if !retVal {
		validateLogger.Println("The value in " + key + " is required!")
	}
	return retVal
}

type NumberFieldValidator struct {
}

func (this NumberFieldValidator) Validate(key, value string) bool {
	retVal, _ := regexp.MatchString("^[0-9]+$", value)
	if !retVal {
		validateLogger.Println("Only numbers are allowed for field \"" + key + "\" = \"" + value + "\".")
	}
	return retVal
}

type ChineseFieldValidator struct {
}

func (this ChineseFieldValidator) Validate(key, value string) bool {
	retVal, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+$", value)
	if !retVal {
		validateLogger.Println("Only Chinese characters are allowed for field \"" + key + "\" = \"" + value + "\".")
	}
	return retVal
}

type AsciiFieldValidator struct {
}

func (this AsciiFieldValidator) Validate(key, value string) bool {
	retVal, _ := regexp.MatchString("^[a-zA-Z]+$", value)
	if !retVal {
		validateLogger.Println("Only English characters are allowed for field \"" + key + "\" = \"" + value + "\".")
	}
	return retVal
}

type EmailAddressFieldValidator struct {
}

func (this EmailAddressFieldValidator) Validate(key, value string) bool {
	retVal, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, value)
	if !retVal {
		validateLogger.Println("No valid email address in field \"" + key + "\" = \"" + value + "\".")
	}
	return retVal
}

type DateFieldValidator struct {
}

func (this DateFieldValidator) Validate(key, value string) bool {
	_, err := time.Parse("02.01.2006", value)
	if err != nil {
		validateLogger.Println("No valid date format in field \"" + key + "\" = \"" + value + "\". Valid: \"DD.MM.YYYY\"")
		return false
	}
	return true
}

// Initialize the validator via this interface for Radio buttons, Checkboxes and drop down lists
type ListValidatorInitializer interface {
	Initialize(listEntries []string)
}

type ListFieldValidator struct {
	listEntries map[string]bool
}

func (this *ListFieldValidator) Initialize(listEntries []string) {
	if this.listEntries == nil {
		this.listEntries = make(map[string]bool, len(listEntries))
	}
	for _, entry := range listEntries {
		this.listEntries[strings.ToLower(entry)] = true
	}
}

func (this ListFieldValidator) Validate(key, value string) bool {
	retVal := this.listEntries[strings.ToLower(value)]
	if !retVal {
		validateLogger.Println("No valid entry in list \"" + key + "\" = \"" + value + "\".")
	}
	return retVal
}
