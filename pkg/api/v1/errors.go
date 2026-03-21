package api

import (
	"encoding/json"
	"fmt"
)

//===========================================================================
// Standard Error Handling
//===========================================================================

var (
	Unsuccessful  = Reply{Success: false}
	NotFound      = Reply{Success: false, Error: "resource not found"}
	NotAllowed    = Reply{Success: false, Error: "method not allowed"}
	InternalError = Reply{Success: false, Error: "an internal error occurred"}
)

// Construct a new response for an error or simply return unsuccessful.
func Error(err any) Reply {
	if err == nil {
		return Unsuccessful
	}

	rep := Reply{Success: false}
	switch err := err.(type) {
	case error:
		rep.Error = err.Error()
	case string:
		rep.Error = err
	case fmt.Stringer:
		rep.Error = err.String()
	case json.Marshaler:
		data, e := err.MarshalJSON()
		if e != nil {
			panic(err)
		}
		rep.Error = string(data)
	default:
		rep.Error = "unhandled error response"
	}

	return rep
}
