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

	mitie "github.com/michlabs/gomitie"
	"github.com/michlabs/cgoutil"
)

// TrainingInstance repesents a sample for training the entity extractor
// Call method Free() to free C underlying structure of this object
type TrainingInstance struct {
	trainingInstance *C.mitie_ner_training_instance
	NumOfEntities int
	NumOfTokens int
}

// NewTrainingInstance returns an instance of TrainingInstance and an error
func NewTrainingInstance(tokens []string) (*TrainingInstance, error) {
	strArr := cgoutil.NewCStringArrayFromSlice(tokens)
	defer strArr.Free()

	mitie.Lock()
	i := C.mitie_create_ner_training_instance((**C.char)(strArr.Pointer))
	mitie.Unlock()
	if i == nil {
		return nil, ErrCreatingTrainingInstanceFailed
	}

	return &TrainingInstance{
		trainingInstance: i,
		NumOfTokens: strArr.Length,
		NumOfEntities: 0,
	}, nil
}

// AddEntity adds new entity to the TrainingInstance object
func (t *TrainingInstance) AddEntity(start int, length int, label string) error {
	if length <= 0 {
		return ErrInvalidLength
	}
	if (start + length) > t.NumOfTokens {
		return ErrInvalidStartOrLength
	}

	mitie.Lock()
	overlap := C.mitie_overlaps_any_entity(t.trainingInstance, C.ulong(start), C.ulong(length))
	mitie.Unlock()
	if int(overlap) != 0 {
		return ErrEntityOverlap
	}

	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))

	mitie.Lock()
	if int(C.mitie_add_ner_training_entity(t.trainingInstance, C.ulong(start), C.ulong(length), clabel)) != 0 {
		return ErrAddingEntityFailed
	}
	t.NumOfEntities += 1
	mitie.Unlock()

	return nil
}

// Free frees C underlying structure of TrainingInstance object
func (t *TrainingInstance) Free() {
	mitie.Lock()
	C.mitie_free(unsafe.Pointer(t.trainingInstance))
	mitie.Unlock()
}