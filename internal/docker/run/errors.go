package run

import (
	"fmt"
)

type MisconfiguredTemplate struct {
	Message string
}

func (mt *MisconfiguredTemplate) Error() string {
	return mt.Message
}

func (mt *MisconfiguredTemplate) Is(target error) bool {
	return fmt.Sprintf("%T", target) == fmt.Sprintf("%T", mt)
}
