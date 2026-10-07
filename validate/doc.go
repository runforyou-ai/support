// Package validate provides validators for values whose format comes from an
// external standard, such as international phone numbers and IANA time zone
// names. Validators that also normalize their input return (value, ok).
//
// Example:
//
//	validate.E164("+1 (555) 010-9999") // "+15550109999", true
//	validate.Timezone("Asia/Shanghai") // true
//	validate.Timezone("Local")         // false
package validate
