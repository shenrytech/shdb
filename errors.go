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
	"errors"
)

// Todo: Make it possible to add arguments to the errors.
// For instance NewErrNotFound(item string) error and then
// also a new IsErrNotFound(err).

var (
	ErrNotAnObject      = errors.New("not an object type")
	ErrInvalidType      = errors.New("invalid type")
	ErrNotFound         = errors.New("not found")
	ErrSessionInvalid   = errors.New("session invalid")
	ErrContextCancelled = errors.New("context cancelled")
	errJson             = errors.New("invalid json data")
)
