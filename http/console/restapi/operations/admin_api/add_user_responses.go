// Code generated by go-swagger; DO NOT EDIT.

// This file is part of MinIO Console Server
// Copyright (c) 2021 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

package admin_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/filedag-project/filedag-storage/http/console/models"
	"net/http"

	"github.com/go-openapi/runtime"
)

// AddUserCreatedCode is the HTTP code returned for type AddUserCreated
const AddUserCreatedCode int = 201

/*AddUserCreated A successful response.

swagger:response addUserCreated
*/
type AddUserCreated struct {

	/*
	  In: Body
	*/
	Payload *models.User `json:"body,omitempty"`
}

// NewAddUserCreated creates AddUserCreated with default headers values
func NewAddUserCreated() *AddUserCreated {

	return &AddUserCreated{}
}

// WithPayload adds the payload to the add user created response
func (o *AddUserCreated) WithPayload(payload *models.User) *AddUserCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add user created response
func (o *AddUserCreated) SetPayload(payload *models.User) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddUserCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AddUserDefault Generic error response.

swagger:response addUserDefault
*/
type AddUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddUserDefault creates AddUserDefault with default headers values
func NewAddUserDefault(code int) *AddUserDefault {
	if code <= 0 {
		code = 500
	}

	return &AddUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the add user default response
func (o *AddUserDefault) WithStatusCode(code int) *AddUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the add user default response
func (o *AddUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the add user default response
func (o *AddUserDefault) WithPayload(payload *models.Error) *AddUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add user default response
func (o *AddUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
