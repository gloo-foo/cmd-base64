package base64_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-base64"
)

func ExampleBase64_encode() {
	// echo "Hello World" | base64
	output, _ := testable.Test(command.Base64(), "Hello World\n")
	fmt.Print(output)
	// Output:
	// SGVsbG8gV29ybGQ=
}
