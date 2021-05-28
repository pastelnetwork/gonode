// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

#include "utils.hpp"

namespace wdm {
    
namespace impl {

//! calculates the weighted Blomqvists's beta.
//! @param x, y input data.
//! @param weights an optional vector of weights for the data.
inline double bbeta(const std::vector<double>& x,
                    const std::vector<double>& y,
                    std::vector<double> weights = std::vector<double>())
{
    utils::check_sizes(x, y, weights);
    size_t n = x.size();

    // find the medians
    double med_x = impl::median(x, weights);
    double med_y = impl::median(y, weights);

    if (weights.size() == 0)
        weights = std::vector<double>(n, 1.0);

    // count elements in lower left and upper right quadrants
    double w_acc{0.0};
    for (size_t i = 0; i < n; i++) {
        if ((x[i] <= med_x) && (y[i] <= med_y))
            w_acc += weights[i];
        else if ((x[i] > med_x) && (y[i] > med_y))
            w_acc += weights[i];
    }

    return 2 * w_acc / utils::sum(weights) - 1;
}

}

}
