package docs

import "github.com/cubny/httpqueue/internal/infra/http/api"

// SetTimersRequestWrapper is the wrapper.
// swagger:parameters setTimersRequest
type SetTimersRequestWrapper struct {
	// in:body
	RequestBody api.SetTimersRequest
}

// SetTimerResponseWrapper is the wrapper.
// swagger:response setTimers
type SetTimerResponseWrapper struct {
	// in:body
	RequestBody api.SetTimerResponse
}

// swagger:parameters getTimerRequest
type GetTimerRequestWrapper struct {
	// TimerID that identifies a timer.
	//
	// in:path
	TimerID string `json:"timer_id"`
}

// GetTimerResponseWrapper is the wrapper.
// swagger:response getTimer
type GetTimerResponseWrapper struct {
	// in:body
	RequestBody api.GetTimerResponse
}

// InvalidRequestBody is an error that is used when the request body fails to be decoded.
// swagger:response invalidRequestBody
type InvalidRequestBody struct {
	// The error body.
	// in: body
	ResponseBody struct {
		// The validation message
		//
		// Required: true
		// Code of the error
		Code int
		// Required: true
		// Details of the error
		Details string
	}
}

// InvalidParams is an error that is used when the required input fails validation..
// swagger:response invalidParams
type InvalidParams struct {
	// The error body.
	// in: body
	ResponseBody struct {
		// The validation message
		//
		// Required: true
		// Code of the error
		Code int
		// Required: true
		// Details of the error
		Details string
	}
}

// NotFoundError is an error that is used when the requested resource cannot be found.
// swagger:response notFoundError
type NotFoundError struct {
	// The error body.
	// in: body
	ResponseBody struct {
		// The not found message
		//
		// Required: true
		// Code of the error
		Code int
		// Required: true
		// Details of the error
		Details string
	}
}

// ServerError is a 500 error used to show that there is a problem with the server in processing the request.
// swagger:response serverError
type ServerError struct {
	// The error body.
	// in: body
	ResponseBody struct {
		// The server error message
		//
		// Required: true
		// Code of the error
		Code int
		// Required: true
		// Details of the error
		Details string
	}
}
