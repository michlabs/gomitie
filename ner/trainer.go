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
)

const DefaultNumOfThreads int = 1

// Trainer represents a MITIE NER trainer
type Trainer struct {
	trainer *C.mitie_ner_trainer
	Size int
}

// NewTrainer returns pointer to a new MITIE NER instance and an error
// path: path to a saved total word feature model file
func NewTrainer(path string) (*Trainer, error) {
	totalWordFeature := C.CString(path)
	defer C.free(unsafe.Pointer(totalWordFeature))

	mitie.Lock()
	trainer := C.mitie_create_ner_trainer(totalWordFeature)
	mitie.Unlock()

	if trainer == nil {
		return nil, ErrCannotOpenModel
	}

	return &Trainer{
		trainer: trainer,
		Size: 0,
	}, nil
}

// AddTrainingInstance adds a TrainingInstance object to Trainer
// and increases the size of Trainer by 1
func (t *Trainer) AddTrainingInstance(ins *TrainingInstance) error {
	mitie.Lock()
	defer mitie.Unlock()
	if int(C.mitie_add_ner_training_instance(t.trainer, ins.trainingInstance)) != 0 {
		return ErrAddingTrainingInstanceFailed
	}
	t.Size += 1
	return nil
}

/*
This parameter controls the trade-off between trying to avoid false alarms 
but also detecting everything.
Different values of beta have the following interpretations:
- beta < 1 indicates that you care more about avoiding false alarms than
  missing detections.  The smaller you make beta the more the trainer will
  try to avoid false alarms.
- beta == 1 indicates that you don't have a preference between avoiding
  false alarms or not missing detections.  That is, you care about these
  two things equally.
- beta > 1 indicates that care more about not missing detections than
  avoiding false alarms.
*/
func (t *Trainer) SetBeta(beta float64) {
	mitie.Lock()
	C.mitie_ner_trainer_set_beta(t.trainer, C.double(beta))
	mitie.Unlock()
}

// Beta returns beta value of the trainer
func (t *Trainer) Beta() float64 {
	mitie.Lock()
	defer mitie.Unlock()

	return float64(C.mitie_ner_trainer_get_beta(t.trainer))
}

// SetNumOfThreads sets number of CPU threads used for training
func (t *Trainer) SetNumOfThreads(n int) {
	if n <= 0 {
		n = DefaultNumOfThreads
	}
	mitie.Lock()
	C.mitie_ner_trainer_set_num_threads(t.trainer, C.ulong(n))
	mitie.Unlock()
}

// NumOfThreads returns number of CPU threads used for training
func (t *Trainer) NumOfThreads() int {
	mitie.Lock()
	defer mitie.Unlock()
	return int(C.mitie_ner_trainer_get_num_threads(t.trainer))
}

// Train trains with saved training instances, then returns pointer to 
// new the Extractor and an error object
func (t *Trainer) Train() (*Extractor, error) {
	if t.Size == 0 {
		return nil, ErrTrainerEmpty
	}

	mitie.Lock()
	x := C.mitie_train_named_entity_extractor(t.trainer)
	mitie.Unlock()
	
	if x == nil {
		return nil, ErrTrainingFailed
	}

	return &Extractor{
		ner: x,
	}, nil
}