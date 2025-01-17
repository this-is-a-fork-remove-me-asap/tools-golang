// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"fmt"
	"io"

	"github.com/this-is-a-fork-remove-me-asap/tools-golang/spdx"
)

func renderAnnotation2_2(ann *spdx.Annotation2_2, w io.Writer) error {
	if ann.Annotator.Annotator != "" && ann.Annotator.AnnotatorType != "" {
		fmt.Fprintf(w, "Annotator: %s: %s\n", ann.Annotator.AnnotatorType, ann.Annotator.Annotator)
	}
	if ann.AnnotationDate != "" {
		fmt.Fprintf(w, "AnnotationDate: %s\n", ann.AnnotationDate)
	}
	if ann.AnnotationType != "" {
		fmt.Fprintf(w, "AnnotationType: %s\n", ann.AnnotationType)
	}
	annIDStr := ann.AnnotationSPDXIdentifier.String()
	if annIDStr != "SPDXRef-" {
		fmt.Fprintf(w, "SPDXREF: %s\n", annIDStr)
	}
	if ann.AnnotationComment != "" {
		fmt.Fprintf(w, "AnnotationComment: %s\n", textify(ann.AnnotationComment))
	}

	return nil
}
