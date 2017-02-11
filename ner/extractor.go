package ner

/*
#cgo LDFLAGS: -lmitie

#include <stdlib.h>
#include <stdio.h>
#include "mitie.h"
*/
import "C"
import (
	"unsafe"
	"strings"

	mitie "github.com/michlabs/gomitie"
	"github.com/michlabs/cgoutil"
)

// Extractor extracts named entities from text, based on a pre-trained model
// Call method Free() to free the C underlying structure of this object
type Extractor struct {
	ner *C.mitie_named_entity_extractor
}

// NewExtractor returns pointer to a new named entity extractor
// path: path to a pre-trained model file
func NewExtractor(path string) (*Extractor, error) {
	model := C.CString(path)
	defer C.free(unsafe.Pointer(model))

	mitie.Lock()
	ner := C.mitie_load_named_entity_extractor(model)
	mitie.Unlock()

	if ner == nil {
		return nil, ErrCannotOpenModel
	}

	return &Extractor{
		ner: ner,
	}, nil
}

// Free frees underlying C object
// MUST be called before your application exists
func (e *Extractor) Free() {
	mitie.Lock()
	C.mitie_free(unsafe.Pointer(e.ner))
	mitie.Unlock()
}

// Tags returns a slice of tags in the pre-trained model
// Such as PERSON, ORG,...
func (e *Extractor) Tags() []string {
	mitie.Lock()
	num := int(C.mitie_get_num_possible_ner_tags(e.ner))
	mitie.Unlock()

	var tags []string
	for i := 0; i < num; i++ {
		tags = append(tags, e.tagStr(i))
	}
	return tags
}

func (e *Extractor) tagStr(idx int) string {
	mitie.Lock()
	cstr := C.mitie_get_named_entity_tagstr(e.ner, C.ulong(idx))
	mitie.Unlock()

	return C.GoString(cstr)
}

// Extract runs the extractor and returns a slice of Entities found in the
// given tokens.
func (e *Extractor) Extract(tokens []string) ([]Entity, error) {
	cArr := cgoutil.NewCStringArrayFromSlice(tokens)
	defer cArr.Free()

	mitie.Lock()
	dets := C.mitie_extract_entities(e.ner, (**C.char)(cArr.Pointer))
	mitie.Unlock()

	if dets == nil {
		return nil, ErrExtractionFailed
	}

	mitie.Lock()
	n := int(C.mitie_ner_get_num_detections(dets))
	mitie.Unlock()
	entities := make([]Entity, n, n)
	for i := 0; i < n; i++ {
		mitie.Lock()
		position := int(C.mitie_ner_get_detection_position(dets, C.ulong(i)))
		length := int(C.mitie_ner_get_detection_length(dets, C.ulong(i)))
		tagID := int(C.mitie_ner_get_detection_tag(dets, C.ulong(i)))
		score := float64(C.mitie_ner_get_detection_score(dets, C.ulong(i)))
		mitie.Unlock()

		entities[i] = Entity{
			Tag:   Tag{
				ID: tagID,
				Name: e.tagStr(tagID),
				},
			Value:  strings.Join(tokens[position:position+length], " "),
			Score: score,
		}
	}

	mitie.Lock()
	C.mitie_free(unsafe.Pointer(dets))
	mitie.Unlock()

	return entities, nil
}

// SaveToFile rerializes this object to a binary file for later use
func (e *Extractor) SaveToFile(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	mitie.Lock()
	result := int(C.mitie_save_named_entity_extractor(cs, e.ner))
	mitie.Unlock()
	if result != 0 {
		return ErrSavingToFileFailed
	}
	return nil
}