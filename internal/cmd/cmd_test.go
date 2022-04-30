package cmd_test

import (
	"os"
	"testing"

	"github.com/paulhammond/sup/internal/cmd"
	"github.com/paulhammond/sup/internal/remote"
	"github.com/paulhammond/sup/internal/remote/testutil"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestIntegration(t *testing.T) {

	testscript.Run(t, testscript.Params{
		Dir: "../../tests",
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"s3init": func(ts *testscript.TestScript, neg bool, args []string) {

				t2 := ts.Value("t").(testing.TB)
				s3 := testutil.S3Remote(t2)

				ts.Setenv("S3_REMOTE", s3)
				ts.Setenv("AWS_REGION", os.Getenv("AWS_REGION"))
				ts.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
				ts.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
				ts.Setenv("AWS_SESSION_TOKEN", os.Getenv("AWS_SESSION_TOKEN"))
			},

			"fakeinit": func(ts *testscript.TestScript, neg bool, args []string) {
				fakePath := ts.MkAbs(args[0])
				err := remote.CreateFake(fakePath)
				if err != nil {
					ts.Fatalf("fakeinit error %s", err)
				}
			},
		},

		Setup: func(env *testscript.Env) error {
			// testscript doesn't expose T to cmds so this is a workaround
			// this is a sign we're being too clever
			env.Values["t"] = env.T()
			return nil
		},
	})
}

func TestMain(m *testing.M) {
	remote.Timer = testutil.Timer
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"sup": cmd.Run,
	}))
}
