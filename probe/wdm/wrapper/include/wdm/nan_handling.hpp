// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

#include <limits>
#include <sstream>

namespace wdm {

namespace utils {

inline void remove_incomplete(std::vector<double>& x,
                              std::vector<double>& y,
                              std::vector<double>& w)
{
    // if observation conatins nan, move it to the end
    size_t last = x.size() - 1;
    for (size_t i = 0; i < last + 1; i++) {
        bool row_has_nan = (std::isnan(x[i]) | std::isnan(y[i]));
        if (w.size() > 0)
            row_has_nan = (row_has_nan |  std::isnan(w[i]));
        if (row_has_nan) {
            if (w.size() > 0)
                std::swap(w[i], w[last]);
            std::swap(x[i], x[last]);
            std::swap(y[i--], y[last--]);
        }
    }

    x.resize(last + 1);
    y.resize(last + 1);
    if (w.size() > 0)
        w.resize(last + 1);
}

inline bool any_nan(const std::vector<double>& x) {
    for (size_t i = 0; (i < x.size()); i++) {
        if (std::isnan(x[i]))
            return true;
    }

    return false;
}

inline std::string preproc(std::vector<double>& x,
                           std::vector<double>& y,
                           std::vector<double>& weights,
                           std::string method,
                           bool remove_missing)
{
    size_t min_nobs = (method == "hoeffding") ? 5 : 2;
    if (remove_missing) {
        utils::remove_incomplete(x, y, weights);
        if (x.size() < min_nobs)
            return "return_nan";
    } else {
        std::stringstream msg;
        if (utils::any_nan(x) | utils::any_nan(y) | utils::any_nan(weights)) {
            msg << "there are missing values in the data; " <<
                   "try remove_missing = TRUE";
        } else if (x.size() < min_nobs) {
            msg << "need at least " << min_nobs << "observations.";
        }
        if (!msg.str().empty())
            throw std::runtime_error(msg.str());
    }

    return "continue";
}

} // end utils

} // end wdm
