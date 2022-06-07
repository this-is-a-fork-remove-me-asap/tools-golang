// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfloader

import (
	"io"

	"github.com/spdx/gordf/rdfloader"
	"github.com/this-is-a-fork-remove-me-asap/tools-golang/rdfloader/parser2v2"
	"github.com/this-is-a-fork-remove-me-asap/tools-golang/spdx"
)

// Takes in a file Reader and returns the pertaining spdx document
// or the error if any is encountered while setting the doc.
func Load2_2(content io.Reader) (*spdx.Document2_2, error) {
	var rdfParserObj, err = rdfloader.LoadFromReaderObject(content)
	if err != nil {
		return nil, err
	}

	doc, err := parser2v2.LoadFromGoRDFParser(rdfParserObj)
	return doc, err
}
