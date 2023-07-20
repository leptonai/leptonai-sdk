/*
Lambda Cloud API

API for interacting with the Lambda GPU Cloud

API version: 1.4.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package lambdalabs

import (
	"encoding/json"
)

// checks if the ListFileSystems200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListFileSystems200Response{}

// ListFileSystems200Response struct for ListFileSystems200Response
type ListFileSystems200Response struct {
	Data []FileSystem `json:"data"`
}

// NewListFileSystems200Response instantiates a new ListFileSystems200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListFileSystems200Response(data []FileSystem) *ListFileSystems200Response {
	this := ListFileSystems200Response{}
	this.Data = data
	return &this
}

// NewListFileSystems200ResponseWithDefaults instantiates a new ListFileSystems200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListFileSystems200ResponseWithDefaults() *ListFileSystems200Response {
	this := ListFileSystems200Response{}
	return &this
}

// GetData returns the Data field value
func (o *ListFileSystems200Response) GetData() []FileSystem {
	if o == nil {
		var ret []FileSystem
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *ListFileSystems200Response) GetDataOk() ([]FileSystem, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *ListFileSystems200Response) SetData(v []FileSystem) {
	o.Data = v
}

func (o ListFileSystems200Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListFileSystems200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableListFileSystems200Response struct {
	value *ListFileSystems200Response
	isSet bool
}

func (v NullableListFileSystems200Response) Get() *ListFileSystems200Response {
	return v.value
}

func (v *NullableListFileSystems200Response) Set(val *ListFileSystems200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableListFileSystems200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableListFileSystems200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListFileSystems200Response(val *ListFileSystems200Response) *NullableListFileSystems200Response {
	return &NullableListFileSystems200Response{value: val, isSet: true}
}

func (v NullableListFileSystems200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListFileSystems200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


