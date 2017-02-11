package ner

import (
	"errors"
)

// Some errors returned by this package
var (
	ErrCannotOpenModel = errors.New("cannot open model file")
	ErrSavingToFileFailed = errors.New("failed to save to file")
	ErrExtractionFailed = errors.New("failed to extract entities")
	ErrCreatingTrainingInstanceFailed = errors.New("failed to create training instance")
	ErrInvalidLength = errors.New("invalid length, must be greater than 0")
	ErrInvalidStartOrLength = errors.New("invalid start or length value")
	ErrEntityOverlap = errors.New("entity overlaps with an existing one")
	ErrAddingEntityFailed = errors.New("failed to add entity, might be caused by running out of memory")
	ErrAddingTrainingInstanceFailed = errors.New("failed to add training instance, might be caused by running out of memory")
	ErrTrainerEmpty = errors.New("trainer is empty")
	ErrTrainingFailed = errors.New("failed to train")
)