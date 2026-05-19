package consts

type Status string

const (
	PENDING_STATUS Status = "pending"
	FAILURE_STATUS Status = "failure"
	READY_STATUS   Status = "ready"
	DONE_STATUS    Status = "done"
)
