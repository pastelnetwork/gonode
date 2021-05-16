%module wrapper

%include <typemaps.i>
%include "std_string.i"
%include "std_vector.i"

%{
#include <include/wdm.hpp>
%}

namespace std {
   %template(DoubleVector) vector<double>;
}

namespace wdm {
    double wdm(std::vector<double> x,
                  std::vector<double> y,
                  std::string method,
                  std::vector<double> weights = std::vector<double>(),
                  bool remove_missing = true);

   double wdm(double* x,
                  int sizeX,
                  double* y,
                  int sizeY,
                  std::string method);
}
