// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

#include "utils.hpp"
#include "ranks.hpp"
#include "prho.hpp"

namespace wdm {
    
namespace impl {

//! fast calculation of the weighted Spearman's rho.
//! @param x, y input data.
//! @param weights an optional vector of weights for the data.
inline double srho(std::vector<double> x,
                   std::vector<double> y,
                   std::vector<double> weights = std::vector<double>())
{
    utils::check_sizes(x, y, weights);
    x = rank(x, weights, "average");
    y = rank(y, weights, "average");
    return prho(x, y, weights);
}

}

}
