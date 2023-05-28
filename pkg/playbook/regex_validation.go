package playbook

import (
	"errors"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
)

func CustomRegexValidate(value, pattern string) error {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("custom regex validation failed")
	}

	return nil
}

func RegexPatternValidate(value string, question Question) error {
	var flag bool

	if question.Validation == "domain_name" {
		flag = validateDomain(value)
	} else if question.Validation == "ip_address" {
		flag = validateIP(value)
	} else if question.Validation == "email" {
		flag = validateEmail(value)
	} else if question.Validation == "url" {
		var isAny, isHttps, isHttp bool
		var schema string
		for _, valildPattern := range question.ValidPatterns {
			if valildPattern == "any" {
				isAny = true
			} else if valildPattern == "https" {
				isHttps = true
			} else if valildPattern == "http" {
				isHttp = true
			} else {
				return errors.New("invalid value found in the place of validPatterns field")
			}
		}
		if isAny || (isHttp && isHttps) {
			schema = "any"
		} else if isHttp {
			schema = "http"
		} else if isHttps {
			schema = "https"
		} else {
			return errors.New("found no value of validPattern")
		}
		flag = validateURL(value, schema)
	} else if question.Validation == "integer_range" {
		if question.Range == nil {
			return errors.New("for integer_range the range field is necessary")
		}

		mn := question.Range.Min
		mx := question.Range.Max

		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		flag = validateIntegerRange(intValue, mn, mx)
	} else {
		return errors.New("invalid pattern name for validation, provide valid value")
	}

	if !flag {
		return errors.New("regex pattern validation failed")
	}

	return nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func validateDomain(domain string) bool {
	matched, _ := regexp.MatchString(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`, domain)
	return matched
}

func validateURL(rawurl string, schema string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	if schema == "https" && u.Scheme != "https" {
		return false
	}
	if schema == "http" && u.Scheme != "http" {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func validateIntegerRange(number, min, max int) bool {
	return number >= min && number <= max
}
