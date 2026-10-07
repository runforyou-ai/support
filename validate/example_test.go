package validate_test

import (
	"fmt"

	"github.com/runforyou-ai/support/validate"
)

func ExampleE164() {
	fmt.Println(validate.E164("+1 (555) 010-9999"))
	fmt.Println(validate.E164("555-0199"))
	// Output:
	// +15550109999 true
	//  false
}

func ExampleTimezone() {
	fmt.Println(validate.Timezone("UTC"))
	fmt.Println(validate.Timezone("Local"))
	fmt.Println(validate.Timezone("Mars/Olympus_Mons"))
	// Output:
	// true
	// false
	// false
}
