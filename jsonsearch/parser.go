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

package jsonsearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrJson = errors.New("invalid json")
)

type Parser struct {
	FieldPaths []string
	dec        *json.Decoder
	tok        json.Token
	query      string
	selector   func(string) bool
}

func (p *Parser) next() (err error) {
	p.tok, err = p.dec.Token()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) addHit(fpath string) {
	p.FieldPaths = append(p.FieldPaths, fpath)

}

func (p *Parser) handleElement(fpath string) (err error) {
	switch t := p.tok.(type) {
	case json.Delim:
		if t == '{' {
			return p.parseObj(fpath)
		}
		if t == '[' {
			return p.parseArray(fpath)
		}
		return ErrJson
	case json.Number:
		if p.selector(t.String()) {
			p.addHit(fpath)
		}
	case bool:
		if p.selector(fmt.Sprintf("%v", t)) {
			p.addHit(fpath)
		}
	case float64:
		if p.selector(fmt.Sprintf("%f", t)) {
			p.addHit(fpath)
		}
	case string:
		if p.selector(t) {
			p.addHit(fpath)
		}
	case nil:
		if p.selector("null") {
			p.addHit(fpath)
		}
	}
	return nil
}

func (p *Parser) Parse(fpath string) (err error) {
	if err = p.next(); err != nil {
		return
	}
	return p.handleElement(fpath)
}

func (p *Parser) parseObj(fpath string) (err error) {
	for {
		if err = p.next(); err != nil {
			return
		}
		switch t := p.tok.(type) {
		case json.Delim:
			if t == '}' {
				return nil
			}
			return ErrJson
		case string:
			// This is a field name
			childPath := fpath + "/" + t
			// Next token should be the element
			if err = p.Parse(childPath); err != nil {
				return
			}
		default:
			return ErrJson
		}
	}
}

func (p *Parser) parseArray(fpath string) (err error) {
	idx := 0
	for {
		if err = p.next(); err != nil {
			return
		}
		switch t := p.tok.(type) {
		case json.Delim:
			if t == ']' {
				return nil
			}
			if t == '{' {
				if err = p.parseObj(fmt.Sprintf("%s/@%d", fpath, idx)); err != nil {
					return err
				}
			}
		default:
			if err = p.handleElement(fmt.Sprintf("%s/@%d", fpath, idx)); err != nil {
				return err
			}
		}
		idx++
	}
}

func NewParser(jsonData []byte, selector func(string) bool) *Parser {
	p := &Parser{
		FieldPaths: []string{},
		selector:   selector,
		dec:        json.NewDecoder(bytes.NewReader(jsonData)),
	}
	return p
}
