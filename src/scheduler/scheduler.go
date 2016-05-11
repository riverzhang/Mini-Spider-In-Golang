/* scheduler.go - the interface of mini spider's scheduler */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
scheduler manage request.
*/

package scheduler

import (
	"request"
)

type Scheduler interface {
	Push(req *request.Request)
	Pop() (*request.Request, error)
	Count() int
	Empty() bool
}
