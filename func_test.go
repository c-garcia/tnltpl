package main_test

import (
	"net/http"
	"os/exec"
	"regexp"
	"testing"
)

const (
	NGROK_API_EP = "http://localhost:4040"
)

func TestNgrokIsThere(t *testing.T) {
	resp, err := http.Get(NGROK_API_EP)
	if err != nil {
		t.Fatalf("ngrok is not there")
	}
	defer resp.Body.Close()
}

func TestGenerateDoc(t *testing.T) {

	givenATemplateFileIncludingTheURLPattern := func() string {
		return "testdata/test.tpl"
	}

	whenIExecuteTheToolWithADash := func(f string) (string, error) {
		cmd := exec.Command(
			"./tnltpl", "-in", "testdata/test.tpl", "-out", "-",
		)
		var out string
		var err error
		if outBytes, err := cmd.CombinedOutput(); err == nil {
			out = string(outBytes)
		} else {
			t.Logf("Failed to execute command: %v", err)
		}
		return out, err
	}

	thenIGetTheExecutedTemplateOnStdout := func(o string) {
		re := regexp.MustCompile(`URL:\shttps://[^\.]+.ngrok.io`)
		if !re.MatchString(o) {
			t.Errorf("Template does not contain the https ngrok endpoint")
		}

	}

	andStatusCodeIsZero := func(err error) {
		if err != nil {
			t.Errorf("Exit status expected to be 0 but was %v", err)
		}
	}

	tpl := givenATemplateFileIncludingTheURLPattern()
	out, status := whenIExecuteTheToolWithADash(tpl)
	thenIGetTheExecutedTemplateOnStdout(out)
	andStatusCodeIsZero(status)
}
