package iox_test

import (
	"fmt"

	"github.com/runforyou-ai/support/iox"
)

func ExampleNewHeadTailBuffer() {
	buf := iox.NewHeadTailBuffer(5, 5, " ... ")
	_, _ = fmt.Fprint(buf, "first line\n")
	_, _ = fmt.Fprint(buf, "many lines in the middle\n")
	_, _ = fmt.Fprint(buf, "last")
	fmt.Printf("%q\n", buf.String())
	fmt.Println(buf.Truncated())
	// Output:
	// "first ... \nlast"
	// true
}
