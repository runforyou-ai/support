// Package number formats numbers for display and provides small numeric
// helpers, in the spirit of Laravel's Number class.
//
// Formatting rounds half away from zero on the shortest decimal
// representation of a float64, so 2.675 rounds to 2.68 and 0.125 to 0.13.
// NaN and infinities are rendered as "NaN", "+Inf" and "-Inf".
//
// Example:
//
//	number.Format(1234567.891, 2)  // "1,234,567.89"
//	number.FileSize(1536, 2)       // "1.50 KB"
//	number.Abbreviate(1234567, 1)  // "1.2M"
//	number.ForHumans(1500000, 1)   // "1.5 million"
//	number.Ordinal(22)             // "22nd"
package number
