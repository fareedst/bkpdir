package contract

import "testing"

func TestEvalBool_AtomicWriterPost_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	env := WriterEnv("/tmp/out.txt", "/tmp/out.txt.tmp", true, true, true)
	ok, err := EvalBool("writer.committed == true AND writer.closed == true", env)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected commit post to hold")
	}
}

func TestEvalBool_ValidatePathPre_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	env := Env{"targetpath": "/safe/path"}
	ok, err := EvalBool(`VALIDATE_PATH(targetPath) == ok`, env)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected valid path pre")
	}
	env["targetpath"] = "../bad"
	ok, err = EvalBool(`VALIDATE_PATH(targetPath) == ok`, env)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected invalid path to fail pre")
	}
}

func TestEvalBool_BranchAlreadyCommitted_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	env := WriterEnv("/tmp/x", "", false, true, true)
	ok, err := EvalBool("writer.committed == true", env)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("branch B001 guard")
	}
}
