// A program consists of a sequence of function calls that are evaluated in order.
type Program struct {
    Steps []FuncCall `json:"@steps"`
}

// A function call specifies a function name and a list of argument expressions. Arguments may contain
// nested function calls and result references.
type FuncCall struct {
    // Name of the function
    Func string       `json:"@func"`
    // Arguments for the function, if any
    Args []Expression `json:"@args,omitempty"`
}

// An expression is a int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string | FuncCall | ResultReference.
type Expression any
type Result any

// A result reference represents the value of an expression from a preceding step.
type ResultReference struct {
    // Index of the previous expression in the "@steps" array
    Ref int `json:"@ref"`
}