package message

import "fmt"

type CodeClassType uint8
type CodeDetailType uint8

// Code: class, detail (c.dd)
type CodeType struct {
	CodeClass  CodeClassType
	CodeDetail CodeDetailType
}

// Stringify code
func (c CodeType) String() string {
	return fmt.Sprintf("%d.%02d", c.CodeClass, c.CodeDetail)
}
