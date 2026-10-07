package validate

import "time"

// Timezone reports whether name is an IANA time zone name, such as
// "Asia/Shanghai" or "UTC", that time.LoadLocation can load. It returns false
// for "" and "Local", which time.LoadLocation maps to UTC and the system zone.
//
// Zone data is read from the system at run time. Programs that must validate
// zones on systems without a zone database, such as minimal containers or
// Windows hosts without Go installed, should embed one by importing
// _ "time/tzdata" in their main package.
func Timezone(name string) bool {
	if name == "" || name == "Local" {
		return false
	}
	_, err := time.LoadLocation(name)
	return err == nil
}
