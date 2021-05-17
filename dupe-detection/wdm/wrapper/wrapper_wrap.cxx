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

#define SWIGMODULE wrapper

#ifdef __cplusplus
/* SwigValueWrapper is described in swig.swg */
template<typename T> class SwigValueWrapper {
  struct SwigMovePointer {
    T *ptr;
    SwigMovePointer(T *p) : ptr(p) { }
    ~SwigMovePointer() { delete ptr; }
    SwigMovePointer& operator=(SwigMovePointer& rhs) { T* oldptr = ptr; ptr = 0; delete oldptr; ptr = rhs.ptr; rhs.ptr = 0; return *this; }
  } pointer;
  SwigValueWrapper& operator=(const SwigValueWrapper<T>& rhs);
  SwigValueWrapper(const SwigValueWrapper<T>& rhs);
public:
  SwigValueWrapper() : pointer(0) { }
  SwigValueWrapper& operator=(const T& t) { SwigMovePointer tmp(new T(t)); pointer = tmp; return *this; }
  operator T&() const { return *pointer.ptr; }
  T *operator&() { return pointer.ptr; }
};

template <typename T> T SwigValueInit() {
  return T();
}
#endif

/* -----------------------------------------------------------------------------
 *  This section contains generic SWIG labels for method/variable
 *  declarations/attributes, and other compiler dependent labels.
 * ----------------------------------------------------------------------------- */

/* template workaround for compilers that cannot correctly implement the C++ standard */
#ifndef SWIGTEMPLATEDISAMBIGUATOR
# if defined(__SUNPRO_CC) && (__SUNPRO_CC <= 0x560)
#  define SWIGTEMPLATEDISAMBIGUATOR template
# elif defined(__HP_aCC)
/* Needed even with `aCC -AA' when `aCC -V' reports HP ANSI C++ B3910B A.03.55 */
/* If we find a maximum version that requires this, the test would be __HP_aCC <= 35500 for A.03.55 */
#  define SWIGTEMPLATEDISAMBIGUATOR template
# else
#  define SWIGTEMPLATEDISAMBIGUATOR
# endif
#endif

/* inline attribute */
#ifndef SWIGINLINE
# if defined(__cplusplus) || (defined(__GNUC__) && !defined(__STRICT_ANSI__))
#   define SWIGINLINE inline
# else
#   define SWIGINLINE
# endif
#endif

/* attribute recognised by some compilers to avoid 'unused' warnings */
#ifndef SWIGUNUSED
# if defined(__GNUC__)
#   if !(defined(__cplusplus)) || (__GNUC__ > 3 || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4))
#     define SWIGUNUSED __attribute__ ((__unused__))
#   else
#     define SWIGUNUSED
#   endif
# elif defined(__ICC)
#   define SWIGUNUSED __attribute__ ((__unused__))
# else
#   define SWIGUNUSED
# endif
#endif

#ifndef SWIG_MSC_UNSUPPRESS_4505
# if defined(_MSC_VER)
#   pragma warning(disable : 4505) /* unreferenced local function has been removed */
# endif
#endif

#ifndef SWIGUNUSEDPARM
# ifdef __cplusplus
#   define SWIGUNUSEDPARM(p)
# else
#   define SWIGUNUSEDPARM(p) p SWIGUNUSED
# endif
#endif

/* internal SWIG method */
#ifndef SWIGINTERN
# define SWIGINTERN static SWIGUNUSED
#endif

/* internal inline SWIG method */
#ifndef SWIGINTERNINLINE
# define SWIGINTERNINLINE SWIGINTERN SWIGINLINE
#endif

/* exporting methods */
#if defined(__GNUC__)
#  if (__GNUC__ >= 4) || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4)
#    ifndef GCC_HASCLASSVISIBILITY
#      define GCC_HASCLASSVISIBILITY
#    endif
#  endif
#endif

#ifndef SWIGEXPORT
# if defined(_WIN32) || defined(__WIN32__) || defined(__CYGWIN__)
#   if defined(STATIC_LINKED)
#     define SWIGEXPORT
#   else
#     define SWIGEXPORT __declspec(dllexport)
#   endif
# else
#   if defined(__GNUC__) && defined(GCC_HASCLASSVISIBILITY)
#     define SWIGEXPORT __attribute__ ((visibility("default")))
#   else
#     define SWIGEXPORT
#   endif
# endif
#endif

/* calling conventions for Windows */
#ifndef SWIGSTDCALL
# if defined(_WIN32) || defined(__WIN32__) || defined(__CYGWIN__)
#   define SWIGSTDCALL __stdcall
# else
#   define SWIGSTDCALL
# endif
#endif

/* Deal with Microsoft's attempt at deprecating C standard runtime functions */
#if !defined(SWIG_NO_CRT_SECURE_NO_DEPRECATE) && defined(_MSC_VER) && !defined(_CRT_SECURE_NO_DEPRECATE)
# define _CRT_SECURE_NO_DEPRECATE
#endif

/* Deal with Microsoft's attempt at deprecating methods in the standard C++ library */
#if !defined(SWIG_NO_SCL_SECURE_NO_DEPRECATE) && defined(_MSC_VER) && !defined(_SCL_SECURE_NO_DEPRECATE)
# define _SCL_SECURE_NO_DEPRECATE
#endif

/* Deal with Apple's deprecated 'AssertMacros.h' from Carbon-framework */
#if defined(__APPLE__) && !defined(__ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES)
# define __ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES 0
#endif

/* Intel's compiler complains if a variable which was never initialised is
 * cast to void, which is a common idiom which we use to indicate that we
 * are aware a variable isn't used.  So we just silence that warning.
 * See: https://github.com/swig/swig/issues/192 for more discussion.
 */
#ifdef __INTEL_COMPILER
# pragma warning disable 592
#endif


#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>



typedef long long intgo;
typedef unsigned long long uintgo;


# if !defined(__clang__) && (defined(__i386__) || defined(__x86_64__))
#   define SWIGSTRUCTPACKED __attribute__((__packed__, __gcc_struct__))
# else
#   define SWIGSTRUCTPACKED __attribute__((__packed__))
# endif



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;




#define swiggo_size_assert_eq(x, y, name) typedef char name[(x-y)*(x-y)*-2+1];
#define swiggo_size_assert(t, n) swiggo_size_assert_eq(sizeof(t), n, swiggo_sizeof_##t##_is_not_##n)

swiggo_size_assert(char, 1)
swiggo_size_assert(short, 2)
swiggo_size_assert(int, 4)
typedef long long swiggo_long_long;
swiggo_size_assert(swiggo_long_long, 8)
swiggo_size_assert(float, 4)
swiggo_size_assert(double, 8)

#ifdef __cplusplus
extern "C" {
#endif
extern void crosscall2(void (*fn)(void *, int), void *, int);
extern char* _cgo_topofstack(void) __attribute__ ((weak));
extern void _cgo_allocate(void *, int);
extern void _cgo_panic(void *, int);
#ifdef __cplusplus
}
#endif

static char *_swig_topofstack() {
  if (_cgo_topofstack) {
    return _cgo_topofstack();
  } else {
    return 0;
  }
}

static void _swig_gopanic(const char *p) {
  struct {
    const char *p;
  } SWIGSTRUCTPACKED a;
  a.p = p;
  crosscall2(_cgo_panic, &a, (int) sizeof a);
}




#define SWIG_contract_assert(expr, msg) \
  if (!(expr)) { _swig_gopanic(msg); } else


static void Swig_free(void* p) {
  free(p);
}

static void* Swig_malloc(int c) {
  return malloc(c);
}


#include <string>


#include <vector>
#include <stdexcept>


#include <include/wdm.hpp>

SWIGINTERN std::vector< double >::const_reference std_vector_Sl_double_Sg__get(std::vector< double > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_double_Sg__set(std::vector< double > *self,int i,std::vector< double >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
#ifdef __cplusplus
extern "C" {
#endif

void _wrap_Swig_free_wrapper_d8938e3a6795767e(void *_swig_go_0) {
  void *arg1 = (void *) 0 ;
  
  arg1 = *(void **)&_swig_go_0; 
  
  Swig_free(arg1);
  
}


void *_wrap_Swig_malloc_wrapper_d8938e3a6795767e(intgo _swig_go_0) {
  int arg1 ;
  void *result = 0 ;
  void *_swig_go_result;
  
  arg1 = (int)_swig_go_0; 
  
  result = (void *)Swig_malloc(arg1);
  *(void **)&_swig_go_result = (void *)result; 
  return _swig_go_result;
}


std::vector< double > *_wrap_new_DoubleVector__SWIG_0_wrapper_d8938e3a6795767e() {
  std::vector< double > *result = 0 ;
  std::vector< double > *_swig_go_result;
  
  
  result = (std::vector< double > *)new std::vector< double >();
  *(std::vector< double > **)&_swig_go_result = (std::vector< double > *)result; 
  return _swig_go_result;
}


std::vector< double > *_wrap_new_DoubleVector__SWIG_1_wrapper_d8938e3a6795767e(long long _swig_go_0) {
  std::vector< double >::size_type arg1 ;
  std::vector< double > *result = 0 ;
  std::vector< double > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< double > *)new std::vector< double >(arg1);
  *(std::vector< double > **)&_swig_go_result = (std::vector< double > *)result; 
  return _swig_go_result;
}


long long _wrap_DoubleVector_size_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  std::vector< double >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  
  result = ((std::vector< double > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_DoubleVector_capacity_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  std::vector< double >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  
  result = ((std::vector< double > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_DoubleVector_reserve_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, long long _swig_go_1) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  std::vector< double >::size_type arg2 ;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_DoubleVector_isEmpty_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  
  result = (bool)((std::vector< double > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_DoubleVector_clear_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_DoubleVector_add_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, double _swig_go_1) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  std::vector< double >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  arg2 = (std::vector< double >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< double >::value_type const &)*arg2);
  
}


double _wrap_DoubleVector_get_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, intgo _swig_go_1) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  int arg2 ;
  std::vector< double >::value_type *result = 0 ;
  double _swig_go_result;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< double >::value_type *) &std_vector_Sl_double_Sg__get(arg1,arg2);
  }
  catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  
  _swig_go_result = (double)*result; 
  return _swig_go_result;
}


void _wrap_DoubleVector_set_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, intgo _swig_go_1, double _swig_go_2) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  int arg2 ;
  std::vector< double >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< double >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_double_Sg__set(arg1,arg2,(double const &)*arg3);
  }
  catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  
  
}


void _wrap_delete_DoubleVector_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0) {
  std::vector< double > *arg1 = (std::vector< double > *) 0 ;
  
  arg1 = *(std::vector< double > **)&_swig_go_0; 
  
  delete arg1;
  
}


double _wrap_wdm__SWIG_0_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, std::vector< double > *_swig_go_1, _gostring_ _swig_go_2, std::vector< double > *_swig_go_3, bool _swig_go_4) {
  std::vector< double > arg1 ;
  std::vector< double > arg2 ;
  std::string arg3 ;
  std::vector< double > arg4 ;
  bool arg5 ;
  std::vector< double > *argp1 ;
  std::vector< double > *argp2 ;
  std::vector< double > *argp4 ;
  double result;
  double _swig_go_result;
  
  
  argp1 = (std::vector< double > *)_swig_go_0;
  if (argp1 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg1 = (std::vector< double >)*argp1;
  
  
  argp2 = (std::vector< double > *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg2 = (std::vector< double >)*argp2;
  
  (&arg3)->assign(_swig_go_2.p, _swig_go_2.n); 
  
  argp4 = (std::vector< double > *)_swig_go_3;
  if (argp4 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg4 = (std::vector< double >)*argp4;
  
  arg5 = (bool)_swig_go_4; 
  
  result = (double)wdm::wdm(arg1,arg2,arg3,arg4,arg5);
  _swig_go_result = result; 
  return _swig_go_result;
}


double _wrap_wdm__SWIG_1_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, std::vector< double > *_swig_go_1, _gostring_ _swig_go_2, std::vector< double > *_swig_go_3) {
  std::vector< double > arg1 ;
  std::vector< double > arg2 ;
  std::string arg3 ;
  std::vector< double > arg4 ;
  std::vector< double > *argp1 ;
  std::vector< double > *argp2 ;
  std::vector< double > *argp4 ;
  double result;
  double _swig_go_result;
  
  
  argp1 = (std::vector< double > *)_swig_go_0;
  if (argp1 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg1 = (std::vector< double >)*argp1;
  
  
  argp2 = (std::vector< double > *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg2 = (std::vector< double >)*argp2;
  
  (&arg3)->assign(_swig_go_2.p, _swig_go_2.n); 
  
  argp4 = (std::vector< double > *)_swig_go_3;
  if (argp4 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg4 = (std::vector< double >)*argp4;
  
  
  result = (double)wdm::wdm(arg1,arg2,arg3,arg4);
  _swig_go_result = result; 
  return _swig_go_result;
}


double _wrap_wdm__SWIG_2_wrapper_d8938e3a6795767e(std::vector< double > *_swig_go_0, std::vector< double > *_swig_go_1, _gostring_ _swig_go_2) {
  std::vector< double > arg1 ;
  std::vector< double > arg2 ;
  std::string arg3 ;
  std::vector< double > *argp1 ;
  std::vector< double > *argp2 ;
  double result;
  double _swig_go_result;
  
  
  argp1 = (std::vector< double > *)_swig_go_0;
  if (argp1 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg1 = (std::vector< double >)*argp1;
  
  
  argp2 = (std::vector< double > *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< double >");
  }
  arg2 = (std::vector< double >)*argp2;
  
  (&arg3)->assign(_swig_go_2.p, _swig_go_2.n); 
  
  result = (double)wdm::wdm(arg1,arg2,arg3);
  _swig_go_result = result; 
  return _swig_go_result;
}


double _wrap_wdm__SWIG_3_wrapper_d8938e3a6795767e(double *_swig_go_0, intgo _swig_go_1, double *_swig_go_2, intgo _swig_go_3, _gostring_ _swig_go_4) {
  double *arg1 = (double *) 0 ;
  int arg2 ;
  double *arg3 = (double *) 0 ;
  int arg4 ;
  std::string arg5 ;
  double result;
  double _swig_go_result;
  
  arg1 = *(double **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = *(double **)&_swig_go_2; 
  arg4 = (int)_swig_go_3; 
  (&arg5)->assign(_swig_go_4.p, _swig_go_4.n); 
  
  result = (double)wdm::wdm(arg1,arg2,arg3,arg4,arg5);
  _swig_go_result = result; 
  return _swig_go_result;
}


#ifdef __cplusplus
}
#endif

