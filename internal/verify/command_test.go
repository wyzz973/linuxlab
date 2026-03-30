package verify

import (
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestCommandOutputVerifier_Pass(t *testing.T) {
	v := &CommandOutputVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "command_output", Command: "echo hello", Expect: "hello"})
	if !r.Passed {
		t.Fatalf("expected pass, got fail: %s", r.Message)
	}
}

func TestCommandOutputVerifier_Fail(t *testing.T) {
	v := &CommandOutputVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "command_output", Command: "echo wrong", Expect: "right"})
	if r.Passed {
		t.Fatal("expected fail, got pass")
	}
}

func TestCommandOutputVerifier_MultiLine(t *testing.T) {
	v := &CommandOutputVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "command_output", Command: `printf "line1\nline2"`, Expect: "line1\nline2"})
	if !r.Passed {
		t.Fatalf("expected pass, got fail: %s", r.Message)
	}
}

func TestExitCodeVerifier_Pass(t *testing.T) {
	v := &ExitCodeVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "exit_code", Command: "true", Expect: "0"})
	if !r.Passed {
		t.Fatalf("expected pass, got fail: %s", r.Message)
	}
}

func TestExitCodeVerifier_Fail(t *testing.T) {
	v := &ExitCodeVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "exit_code", Command: "false", Expect: "0"})
	if r.Passed {
		t.Fatal("expected fail, got pass")
	}
}

func TestExitCodeVerifier_NonZeroExpected(t *testing.T) {
	v := &ExitCodeVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "exit_code", Command: "false", Expect: "1"})
	if !r.Passed {
		t.Fatalf("expected pass, got fail: %s", r.Message)
	}
}
