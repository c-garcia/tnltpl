package main_test

import (
	"net/http"
	"os/exec"
	"regexp"
	"testing"
)

const (
	ngrokAPIEndPoint = "http://localhost:4040"
)

func TestNgrokIsThere(t *testing.T) {
	resp, err := http.Get(ngrokAPIEndPoint)
	if err != nil {
		t.Fatalf("ngrok is not there")
	}
	defer resp.Body.Close()
}

func TestRunsTemplate(t *testing.T) {

	givenATemplateFileIncludingTheURLPattern := func() string {
		return "testdata/test.tpl"
	}

	whenIExecuteTheToolWithADash := func(f string) (string, error) {
		inst := exec.Command("go", "install",".")
		if _, err := inst.CombinedOutput(); err != nil  {
			t.Errorf("Can't install current package %v", err)
			return "", err
		}
		cmd := exec.Command(
			"tnltpl", "-in", "testdata/test.tpl", "-out", "-",
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

	andStatusCodeIsNotError := func(err error) {
		if err != nil {
			t.Errorf("Exit status expected to be 0 but was %v", err)
		}
	}

	tpl := givenATemplateFileIncludingTheURLPattern()
	out, status := whenIExecuteTheToolWithADash(tpl)
	thenIGetTheExecutedTemplateOnStdout(out)
	andStatusCodeIsNotError(status)
}
