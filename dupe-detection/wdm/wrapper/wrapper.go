/* ----------------------------------------------------------------------------
 * This file was automatically generated by SWIG (http://www.swig.org).
 * Version 3.0.12
 *
 * This file is not intended to be easily readable and contains a number of
 * coding conventions designed to improve portability and efficiency. Do not make
 * changes to this file unless you know what you are doing--modify the SWIG
 * interface file instead.
 * ----------------------------------------------------------------------------- */

// source: wrapper.i

package wrapper

/*
#define intgo swig_intgo
typedef void *swig_voidp;

#include <stdint.h>


typedef long long intgo;
typedef unsigned long long uintgo;



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;


typedef _gostring_ swig_type_1;
extern void _wrap_Swig_free_wrapper_510c6f704a9db668(uintptr_t arg1);
extern uintptr_t _wrap_Swig_malloc_wrapper_510c6f704a9db668(swig_intgo arg1);
extern double _wrap_wdm_wrapper_510c6f704a9db668(swig_voidp arg1, swig_intgo arg2, swig_voidp arg3, swig_intgo arg4, swig_type_1 arg5);
#undef intgo
*/
import "C"

import (
	_ "runtime/cgo" // Autogenerated
	"sync"
	"unsafe"
)

type _ unsafe.Pointer

// Swigescapealwaysfalse is Autogenerated
var Swigescapealwaysfalse bool

// Swigescapeval is Autogenerated
var Swigescapeval interface{}

type swigfnptr *byte
type swigmemberptr *byte

type _ sync.Mutex

// Swigfree is Autogenerated
func Swigfree(arg1 uintptr) {
	swigi0 := arg1
	C._wrap_Swig_free_wrapper_510c6f704a9db668(C.uintptr_t(swigi0))
}

// Swigmalloc is Autogenerated
func Swigmalloc(arg1 int) (swigret uintptr) {
	var swigr uintptr
	swigi0 := arg1
	swigr = (uintptr)(C._wrap_Swig_malloc_wrapper_510c6f704a9db668(C.swig_intgo(swigi0)))
	return swigr
}

// Wdm is Autogenerated
func Wdm(arg1 *float64, arg2 int, arg3 *float64, arg4 int, arg5 string) (swigret float64) {
	var swigr float64
	swigi0 := arg1
	swigi1 := arg2
	swigi2 := arg3
	swigi3 := arg4
	swigi4 := arg5
	swigr = (float64)(C._wrap_wdm_wrapper_510c6f704a9db668(C.swig_voidp(swigi0), C.swig_intgo(swigi1), C.swig_voidp(swigi2), C.swig_intgo(swigi3), *(*C.swig_type_1)(unsafe.Pointer(&swigi4))))
	if Swigescapealwaysfalse {
		Swigescapeval = arg5
	}
	return swigr
}
