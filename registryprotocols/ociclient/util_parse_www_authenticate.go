// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "fmt"

// This function parses a single WWW-Authenticate header and returns the auth schemes encountered.
//
// For implementation details see: https://www.rfc-editor.org/rfc/rfc7235#section-4.1
func parseWWWAuthenticate(wwwAuthenticate string) ([]OCIRawAuthScheme, error) {
	type parserState int

	const (
		parserStateStart             parserState = iota
		parserStateBeforeSchemeOrKey parserState = iota
		parserStateInScheme          parserState = iota
		parserStateBeforeKey         parserState = iota
		parserStateInKey             parserState = iota
		parserStateInSchemeOrKey     parserState = iota
		parserStateAfterKey          parserState = iota
		parserStateAfterKeyOrScheme  parserState = iota
		parserStateBeforeValue       parserState = iota
		parserStateInValue           parserState = iota
		parserStateInQuotes          parserState = iota
		parserStateAfterEscape       parserState = iota
		parserStateAfterValue        parserState = iota
	)

	state := parserStateStart

	isWhitespace := func(c int32) bool {
		return c == 9 || // Tab
			c == 10 || // Line feed
			c == 13 || // Newline
			c == 32 // Space
	}

	var authSchemes []OCIRawAuthScheme

	buf := ""
	currentKey := ""
	for i, c := range wwwAuthenticate {
		if (c < 33 || c > 126) && !isWhitespace(c) {
			// Non-printable character. This shouldn't happen if the HTTP client did its job,
			// but we will ignore this header.
			return nil, fmt.Errorf("invalid www-authenticate header, encountered non-printable character at position %d", i)
		}
		switch state {
		case parserStateStart:
			switch {
			case c == ',':
				continue
			case c == '=':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered equals sign at the start of the scheme")
			case isWhitespace(c):
				continue
			default:
				buf += string(c)
				state = parserStateInScheme
			}
		case parserStateBeforeSchemeOrKey:
			switch {
			case c == ',':
				continue
			case c == '=':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered equals sign at the start of the scheme or key")
			case isWhitespace(c):
				continue
			default:
				buf += string(c)
				state = parserStateInSchemeOrKey
			}
		case parserStateInScheme:
			switch {
			case isWhitespace(c):
				authSchemes = append(authSchemes, OCIRawAuthScheme{
					buf,
					map[string]string{},
				})
				buf = ""
				state = parserStateBeforeKey
			case c == '=':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered equals sign in scheme at position %d", i)
			case c == '"':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered quotes in scheme at position %d", i)
			default:
				buf += string(c)
				continue
			}
		case parserStateInSchemeOrKey:
			switch {
			case isWhitespace(c):
				currentKey = buf
				buf = ""
				state = parserStateAfterKeyOrScheme
				continue
			case c == '=':
				currentKey = buf
				buf = ""
				state = parserStateBeforeValue
				continue
			case c == '"':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered quotes in scheme or key at position %d", i)
			default:
				buf += string(c)
				continue
			}
		case parserStateBeforeKey:
			switch {
			case isWhitespace(c):
				continue
			case c == '=':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered equals sign before key at position %d", i)
			case c == '"':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered quotes in key at position %d", i)
			default:
				buf += string(c)
				state = parserStateInKey
				continue
			}
		case parserStateInKey:
			switch {
			case isWhitespace(c):
				state = parserStateAfterKey
				continue
			case c == '=':
				currentKey = buf
				buf = ""
				state = parserStateBeforeValue
				continue
			case c == '"':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered quotes in key at position %d", i)
			default:
				buf += string(c)
				continue
			}
		case parserStateAfterKey:
			switch {
			case c == '"':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered quotes after key at position %d", i)
			case c == '=':
				currentKey = buf
				buf = ""
				state = parserStateBeforeValue
				continue
			case isWhitespace(c):
				continue
			default:
				return nil, fmt.Errorf("invalid www-authenticate header, encountered a key without a value at position %d", i)
			}
		case parserStateAfterKeyOrScheme:
			switch {
			case c == '"':
				return nil, fmt.Errorf("invalid www-authenticate header, encountered quotes after scheme or key at position %d", i)
			case c == '=':
				currentKey = buf
				buf = ""
				state = parserStateBeforeValue
				continue
			case isWhitespace(c):
				continue
			default:
				// We have encountered a new scheme
				authSchemes = append(authSchemes, OCIRawAuthScheme{
					currentKey,
					map[string]string{},
				})
				currentKey = ""
				buf = string(c)
				state = parserStateInKey
			}
		case parserStateBeforeValue:
			switch {
			case c == '"':
				state = parserStateInQuotes
				continue
			case isWhitespace(c):
				continue
			}
			fallthrough
		case parserStateInValue:
			state = parserStateInValue
			switch {
			case isWhitespace(c):
				// End of value since we encountered a whitespace
				authSchemes[len(authSchemes)-1].Params[currentKey] = buf
				currentKey = ""
				buf = ""
				state = parserStateAfterValue
				continue
			case c == ',':
				// End of value, after comma
				authSchemes[len(authSchemes)-1].Params[currentKey] = buf
				currentKey = ""
				buf = ""
				state = parserStateBeforeSchemeOrKey
				continue
			default:
				buf += string(c)
				continue
			}
		case parserStateInQuotes:
			switch c {
			case '"':
				authSchemes[len(authSchemes)-1].Params[currentKey] = buf
				currentKey = ""
				buf = ""
				state = parserStateAfterValue
				continue
			case '\\':
				state = parserStateAfterEscape
				continue
			default:
				buf += string(c)
				continue
			}
		case parserStateAfterEscape:
			switch c {
			case '\\':
				buf += "\\"
			case '"':
				buf += "\""
			default:
				// Assume incorrect escape sequence:
				buf += "\\" + string(c)
			}
			state = parserStateInQuotes
			continue
		case parserStateAfterValue:
			switch {
			case c == ',':
				state = parserStateBeforeSchemeOrKey
				continue
			case isWhitespace(c):
				continue
			default:
				return nil, fmt.Errorf("invalid www-authenticate header, encountered non-whitespace, non-comma character after value at position %d", i)
			}
		}
	}
	switch state {
	case parserStateStart:
	case parserStateAfterValue:
	case parserStateInValue:
		authSchemes[len(authSchemes)-1].Params[currentKey] = buf
	case parserStateInScheme:
		fallthrough
	case parserStateInSchemeOrKey:
		authSchemes = append(authSchemes, OCIRawAuthScheme{
			buf,
			map[string]string{},
		})
	default:
		return nil, fmt.Errorf("invalid www-authenticate header, unexpected end of input")
	}
	return authSchemes, nil
}
