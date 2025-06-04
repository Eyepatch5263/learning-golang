package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody(r *http.Request, x any)error{
	//closing the req body
	defer r.Body.Close()

	//reading the body
	body,err:=io.ReadAll(r.Body);

	if err != nil {
		return err
	}

	if err=json.Unmarshal(body, x); err != nil {
		return err
	}
	return nil
}