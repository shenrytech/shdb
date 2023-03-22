// Copyright 2023 Shenry Tech AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shdb

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TypeId returns the TypeId representation of the type and Id
func (m *Metadata) TypeId() TypeId {
	res := TypeId{}
	res.SetType(TypeKey(m.Type))
	res.SetUuidBytes(m.Uuid)
	return res
}

// Return the Metadata as an *ObjRef
func (m *Metadata) Ref() *ObjRef {
	return &ObjRef{
		Type: m.Type,
		Uuid: m.Uuid,
	}
}

// Fill the metadata with fields that are missing
func (m *Metadata) Fill() error {
	if m.CreatedAt == nil {
		m.CreatedAt = timestamppb.Now()
	}
	if m.Uuid == nil {
		id, err := uuid.New().MarshalBinary()
		if err != nil {
			return err
		}
		m.Uuid = id
	}
	if m.Labels == nil {
		m.Labels = make([]string, 0)
	}
	return nil
}

// GetUuidAsString returns the UUID as a string
func (m *Metadata) GetUuidAsString() (string, error) {
	us, err := m.GetUuidAsUUID()
	if err != nil {
		return "", err
	}
	return us.String(), nil
}

// GetUuidAsUUID returns the UUID as a uuid.UUID
func (m *Metadata) GetUuidAsUUID() (uuid.UUID, error) {
	return uuid.FromBytes(m.Uuid)
}
