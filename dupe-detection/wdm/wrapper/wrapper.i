%module wrapper

%include <typemaps.i>
%include "std_string.i"

%{
#include <include/wdm.hpp>
%}

namespace wdm {
   double wdm(double* x,
                  int sizeX,
                  double* y,
                  int sizeY,
                  std::string method);
}
