package customerror

import "errors"

type CustomError error

// KEY: proto.ProtoReflect().Descriptor().FullName()
// VALUE: customerror type

var ContextDeadlineExceed CustomError = errors.New("waiting too long for response")
