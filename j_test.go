package main

import (
	"strings"
	"testing"
)

func TestJo(t *testing.T) {
	var tests = []struct {
		arrayflag   bool
		boolflag    bool
		verflag     bool
		jsonverflag bool
		prettyflag  bool
		input       string
		output      string
	}{
		{true, false, false, false, false, "1\n2\n3", "[1,2,3]"},
		{false, false, false, false, false, "name=clark\nage=28\nsuperman=true", `{"age":28,"name":"clark","superman":true}`},
		{false, false, false, false, false, "name=x\nvalue=", `{"name":"x","value":null}`},
	}

	for _, test := range tests {
		arrayflag = test.arrayflag
		boolflag = test.boolflag
		verflag = test.verflag
		jsonverflag = test.jsonverflag
		prettyflag = test.prettyflag

		outputbuf, _ := tojson(strings.Split(test.input, "\n"))

		if output := strings.Trim(string(outputbuf), " \n"); output != test.output {
			t.Errorf("Got %v, expected: %v\n", output, test.output)
		}
	}

}
