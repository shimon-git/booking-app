package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/shimon-git/booking-app/internal/config"
)

var app *config.AppConfig

// NewHelpers sets the config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println("Server error", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

/*
StringToDateTime - Converting date from string to time.Time object
(e.g: "Y-M-D","2023-6-24")
*/
func StringToDateTime(dateFormat string, userDate string) (time.Time, error) {
	var year, day, month int
	var err error
	splitedDateFormat := strings.Split(dateFormat, "-")
	splitedUserDate := strings.Split(userDate, "-")
	if len(splitedDateFormat) != 3 && len(splitedUserDate) != 3 {
		return time.Now(), errors.New("invalid time format,\nvalid format needs to contain: {D,Y,M} spereated by '-' as string")
	}
	for idx, timePart := range splitedDateFormat {
		switch timePart {
		case "D":
			day, err = strconv.Atoi(splitedUserDate[idx])
			if err != nil {
				return time.Time{}, fmt.Errorf("cannot Parse the Day\n%v", err)
			}
		case "M":
			month, err = strconv.Atoi(splitedUserDate[idx])
			if err != nil {
				return time.Time{}, fmt.Errorf("cannot Parse the Month\n%v", err)
			}
		case "Y":
			year, err = strconv.Atoi(splitedUserDate[idx])
			if err != nil {
				return time.Time{}, fmt.Errorf("cannot Parse the Year\n%v", err)
			}
		default:
			return time.Time{}, errors.New("invalid time format,\nvalid format needs to contain: {D,Y,M} spereated by '-' as string")
		}
	}
	location, err := time.LoadLocation("Israel")
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, time.Now().Hour(), time.Now().Minute(), time.Now().Second(),
		time.Now().Nanosecond(), location), nil
}

func DateTimeToString(dateFormat time.Time, dateLayout string) (string, error) {
	var year, day, month int
	var reformatedDate []string
	splitedDateLayout := strings.Split(dateLayout, "-")
	if len(splitedDateLayout) != 3 {
		return "", errors.New("invalid time format,\nvalid format needs to contain: {D,Y,M} spereated by '-' as string")
	}
	year = dateFormat.Year()
	month = int(dateFormat.Month())
	day = dateFormat.Day()

	for _, datePart := range splitedDateLayout {
		switch datePart {
		case "D":
			reformatedDate = append(reformatedDate, strconv.Itoa(day))
		case "M":
			reformatedDate = append(reformatedDate, strconv.Itoa(month))
		case "Y":
			reformatedDate = append(reformatedDate, strconv.Itoa(year))
		default:
			return "", errors.New("invalid time format,\nvalid format needs to contain: {D,Y,M} spereated by '-' as string")
		}
	}
	return strings.Join(reformatedDate, "-"), nil
}
