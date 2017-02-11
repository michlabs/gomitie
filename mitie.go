package gomitie

/*
#cgo LDFLAGS: -lmitie

#include <stdlib.h>
#include <stdio.h>
#include "mitie.h"
*/
import "C"
import (
	"sync"
	"unsafe"

    "github.com/michlabs/cgoutil"
)

var mutex *sync.Mutex = &sync.Mutex{}

// lock MITIE
// To access MITIE in multiple goroutines,
// this function MUST be called for thread safety
func Lock() {
	mutex.Lock()
}

// unlock MITIE
func Unlock() {
	mutex.Unlock()
}

// Tokenize returns a slice that contains a tokenized copy of the input text.
func Tokenize(text string) []string {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	ctokens := C.mitie_tokenize(cstr)
	defer C.mitie_free(unsafe.Pointer(ctokens))

    cArrString := cgoutil.NewCStringArrayFromPointer(unsafe.Pointer(ctokens))
    
    return cArrString.ToSlice()
}