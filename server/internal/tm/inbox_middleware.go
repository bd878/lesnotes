package tm

import (
	"fmt"
)

type ErrDuplicateMessage string

func (e ErrDuplicateMessage) Error() string {
	return fmt.Sprintf("duplicate message id encountered: %s", string(e))
}