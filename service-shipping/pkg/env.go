package pkg

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func LoadEnv() (map[string]string, error) {
	//Production mode user local environment
	if _, ok := os.LookupEnv("MODE"); ok && os.Getenv("MODE") == "production" {
		if _, ok := os.LookupEnv("PORT"); !ok {
			return nil, errors.New("PORT not set! ")
		}

		if _, ok := os.LookupEnv("DB_DSN"); !ok {
			return nil, errors.New("DB_DSN not set! ")
		}

		if _, ok := os.LookupEnv("JWT_SIGNING_KEY"); !ok {
			return nil, errors.New("JWT_SIGNING_KEY not set! ")
		}

		if _, ok := os.LookupEnv("PROVIDER_HTTP_API_PROXY_URI"); !ok {
			return nil, errors.New("PROVIDER_HTTP_API_PROXY_URI not set! ")
		}
		//Shippop
		if _, ok := os.LookupEnv("PROVIDER_SHIPPOP_API_URL"); !ok {
			return nil, errors.New("PROVIDER_SHIPPOP_API_URL not set! ")
		}
		if _, ok := os.LookupEnv("PROVIDER_SHIPPOP_API_KEY"); !ok {
			return nil, errors.New("PROVIDER_SHIPPOP_API_KEY not set! ")
		}
		//Flash
		if _, ok := os.LookupEnv("PROVIDER_FLASH_API_URL"); !ok {
			return nil, errors.New("PROVIDER_FLASH_API_URL not set! ")
		}
		if _, ok := os.LookupEnv("PROVIDER_FLASH_API_KEY"); !ok {
			return nil, errors.New("PROVIDER_FLASH_API_KEY not set! ")
		}
		if _, ok := os.LookupEnv("PROVIDER_FLASH_MCH_ID"); !ok {
			return nil, errors.New("PROVIDER_FLASH_MCH_ID not set! ")
		}

		if _, ok := os.LookupEnv("JNA_NAME"); !ok {
			return nil, errors.New("JNA_NAME not set! ")
		}
		if _, ok := os.LookupEnv("JNA_PHONE_NUMBER"); !ok {
			return nil, errors.New("JNA_PHONE_NUMBER not set! ")
		}
		if _, ok := os.LookupEnv("JNA_ADDRESS"); !ok {
			return nil, errors.New("JNA_ADDRESS not set! ")
		}
		if _, ok := os.LookupEnv("JNA_DISTRICT"); !ok {
			return nil, errors.New("JNA_DISTRICT not set! ")
		}
		if _, ok := os.LookupEnv("JNA_STATE"); !ok {
			return nil, errors.New("JNA_STATE not set! ")
		}
		if _, ok := os.LookupEnv("JNA_PROVINCE"); !ok {
			return nil, errors.New("JNA_PROVINCE not set! ")
		}
		if _, ok := os.LookupEnv("JNA_POSTCODE"); !ok {
			return nil, errors.New("JNA_POSTCODE not set! ")
		}
		if _, ok := os.LookupEnv("JNA_ADDRESS_AS_ORIGIN"); !ok {
			return nil, errors.New("JNA_ADDRESS_AS_ORIGIN not set! ")
		}

		if _, ok := os.LookupEnv("AN_URL"); !ok {
			return nil, errors.New("AN_URL not set! ")
		}

		env := map[string]string{
			"PORT":                        os.Getenv("PORT"),
			"DB_DSN":                      os.Getenv("DB_DSN"),
			"JWT_SIGNING_KEY":             os.Getenv("JWT_SIGNING_KEY"),
			"PROVIDER_HTTP_API_PROXY_URI": os.Getenv("PROVIDER_HTTP_API_PROXY_URI"),
			"PROVIDER_SHIPPOP_API_URL":    os.Getenv("PROVIDER_SHIPPOP_API_URL"),
			"PROVIDER_SHIPPOP_API_KEY":    os.Getenv("PROVIDER_SHIPPOP_API_KEY"),
			"PROVIDER_FLASH_API_URL":      os.Getenv("PROVIDER_FLASH_API_URL"),
			"PROVIDER_FLASH_API_KEY":      os.Getenv("PROVIDER_FLASH_API_KEY"),
			"PROVIDER_FLASH_MCH_ID":       os.Getenv("PROVIDER_FLASH_MCH_ID"),
			"JNA_NAME":                    os.Getenv("JNA_NAME"),
			"JNA_PHONE_NUMBER":            os.Getenv("JNA_PHONE_NUMBER"),
			"JNA_ADDRESS":                 os.Getenv("JNA_ADDRESS"),
			"JNA_DISTRICT":                os.Getenv("JNA_DISTRICT"),
			"JNA_STATE":                   os.Getenv("JNA_STATE"),
			"JNA_PROVINCE":                os.Getenv("JNA_PROVINCE"),
			"JNA_POSTCODE":                os.Getenv("JNA_POSTCODE"),
			"JNA_ADDRESS_AS_ORIGIN":       os.Getenv("JNA_ADDRESS_AS_ORIGIN"),
			"AN_URL":                      os.Getenv("AN_URL"),
		}
		return env, nil
	}

	//Develop mode use environment from file
	file, _ := filepath.Abs(".env")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}
	env, _ := godotenv.Read(file)
	return env, nil
}
