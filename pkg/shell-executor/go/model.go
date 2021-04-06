package shell_executor

type IShellExecutor interface {
	Execute(args ...string) (stdOut string, stdErr string, err error)
}
