// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"
	"strings"

	"github.com/this-is-a-fork-remove-me-asap/tools-golang/spdx"
)

func (parser *tvParser2_1) parsePairFromCreationInfo2_1(tag string, value string) error {
	// fail if not in Creation Info parser state
	if parser.st != psCreationInfo2_1 {
		return fmt.Errorf("got invalid state %v in parsePairFromCreationInfo2_1", parser.st)
	}

	// create an SPDX Creation Info data struct if we don't have one already
	if parser.doc.CreationInfo == nil {
		parser.doc.CreationInfo = &spdx.CreationInfo2_1{}
	}

	ci := parser.doc.CreationInfo
	switch tag {
	case "LicenseListVersion":
		ci.LicenseListVersion = value
	case "Creator":
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}

		creator := spdx.Creator{Creator: subvalue}
		switch subkey {
		case "Person", "Organization", "Tool":
			creator.CreatorType = subkey
		default:
			return fmt.Errorf("unrecognized Creator type %v", subkey)
		}

		ci.Creators = append(ci.Creators, creator)
	case "Created":
		ci.Created = value
	case "CreatorComment":
		ci.CreatorComment = value

	// tag for going on to package section
	case "PackageName":
		// error if last file does not have an identifier
		// this may be a null case: can we ever have a "last file" in
		// the "creation info" state? should go on to "file" state
		// even when parsing unpackaged files.
		if parser.file != nil && parser.file.FileSPDXIdentifier == nullSpdxElementId2_1 {
			return fmt.Errorf("file with FileName %s does not have SPDX identifier", parser.file.FileName)
		}
		parser.st = psPackage2_1
		parser.pkg = &spdx.Package2_1{
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: false,
		}
		return parser.parsePairFromPackage2_1(tag, value)
	// tag for going on to _unpackaged_ file section
	case "FileName":
		// leave pkg as nil, so that packages will be placed in Files
		parser.st = psFile2_1
		parser.pkg = nil
		return parser.parsePairFromFile2_1(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_1
		return parser.parsePairFromOtherLicense2_1(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_1
		return parser.parsePairFromReview2_1(tag, value)
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_1{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_1(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_1(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_1{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_1(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in CreationInfo section", tag)
	}

	return nil
}

// ===== Helper functions =====

func extractExternalDocumentReference(value string) (spdx.DocElementID, string, string, string, error) {
	sp := strings.Split(value, " ")
	// remove any that are just whitespace
	keepSp := []string{}
	for _, s := range sp {
		ss := strings.TrimSpace(s)
		if ss != "" {
			keepSp = append(keepSp, ss)
		}
	}

	var documentRefID spdx.DocElementID
	var uri, alg, checksum string

	// now, should have 4 items (or 3, if Alg and Checksum were joined)
	// and should be able to map them
	if len(keepSp) == 4 {
		documentRefID = spdx.MakeDocElementID(keepSp[0], "")
		uri = keepSp[1]
		alg = keepSp[2]
		// check that colon is present for alg, and remove it
		if !strings.HasSuffix(alg, ":") {
			return documentRefID, "", "", "", fmt.Errorf("algorithm does not end with colon")
		}
		alg = strings.TrimSuffix(alg, ":")
		checksum = keepSp[3]
	} else if len(keepSp) == 3 {
		documentRefID = spdx.MakeDocElementID(keepSp[0], "")
		uri = keepSp[1]
		// split on colon into alg and checksum
		parts := strings.SplitN(keepSp[2], ":", 2)
		if len(parts) != 2 {
			return documentRefID, "", "", "", fmt.Errorf("missing colon separator between algorithm and checksum")
		}
		alg = parts[0]
		checksum = parts[1]
	} else {
		return documentRefID, "", "", "", fmt.Errorf("expected 4 elements, got %d", len(keepSp))
	}

	return documentRefID, uri, alg, checksum, nil
}
