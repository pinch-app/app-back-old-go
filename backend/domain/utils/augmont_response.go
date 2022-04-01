package utils

import "github.com/cockroachdb/errors"

type AugmontResponse struct {
	StatusCode int    `json:"statusCode" mapstructure:"statusCode"`
	Message    string `json:"message" mapstructure:"message"`
	Errors     Dict   `json:"errors" mapstructure:"errors"`
	Result     Any    `json:"result" mapstructure:"result"`
}

func (r *AugmontResponse) IsError() bool {
	return r.StatusCode != 200
}

func (r *AugmontResponse) Error() error {
	if !r.IsError() {
		return nil
	}
	errStr := r.Errors.ToString()
	err := errors.New(errStr)
	return errors.WithStackDepth(err, 1)
}
