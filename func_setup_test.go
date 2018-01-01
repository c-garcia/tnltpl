package main_test

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

func waitUntilRegexp(is io.ReadCloser, re string, timeout time.Duration) (bool, error) {
	tryToFindREInOutput := func(ch chan interface{}) {
		compRe := regexp.MustCompile(re)
		found := false
		r := bufio.NewReader(is)
		s := bufio.NewScanner(r)

		for !found && s.Scan() {
			if compRe.MatchString(s.Text()) {
				found = true
			}
		}
		if found {
			ch <- nil
		}
		close(ch)
	}

	reFoundOrTimeout := func(ch chan interface{}) (bool, error) {
		select {
		case _, ok := <-ch:
			return ok, nil
		case <-time.After(timeout):
			return false, nil
		}

	}

	foundCh := make(chan interface{})
	go tryToFindREInOutput(foundCh)
	return reFoundOrTimeout(foundCh)
}

func setupNgrok() (*exec.Cmd, error) {

	const HTTPS_UP = `Forwarding\s*https`

	cmd := exec.Command("ngrok", "http", "8080")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	_, err = waitUntilRegexp(out, HTTPS_UP, 5*time.Second)

	return cmd, err
}

func TestMain(m *testing.M) {
	if cmd, err := setupNgrok(); err == nil {
		testStatus := m.Run()
		cmd.Process.Kill()
		os.Exit(testStatus)
	} else {
		os.Exit(1)
	}
}
