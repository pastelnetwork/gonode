// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

#include "utils.hpp"

namespace wdm {

namespace impl {

inline void normalize_weights(std::vector<double>& w)
{
    if (w.size() > 0) {
        double s = utils::sum(w);
        for (size_t i = 0; i < w.size(); i++)
            w[i] /= s;
    }
}

//! fast calculation of the weighted Kendall's tau.
//! @param x, y input data.
//! @param weights an optional vector of weights for the data.
inline double ktau(std::vector<double> x,
                   std::vector<double> y,
                   std::vector<double> weights = std::vector<double>())
{
    utils::check_sizes(x, y, weights);

    // 1.1 Sort x, y, and weights in x order; break ties in according to y.
    utils::sort_all(x, y, weights);

    // 1.2 Count pairs of tied x and simultaneous ties in x and y.
    double ties_x = utils::count_tied_pairs(x, weights);
    double ties_both = utils::count_joint_ties(x, y, weights);

    // 2.1 Sort y again and count exchanges (= number of discordant pairs).
    double num_d = 0.0;
    utils::merge_sort(y, weights, num_d);

    // 2.2 Count pairs of tied y.
    double ties_y = utils::count_tied_pairs(y, weights);

    // 3. Calculate Kendall's tau.
    if (weights.size() == 0)
        weights = std::vector<double>(x.size(), 1.0);
    double num_pairs = utils::perm_sum(weights, 2);
    double num_c = num_pairs - (num_d + ties_x + ties_y - ties_both);
    double tau = num_c - num_d;
    tau /= std::sqrt((num_pairs - ties_x) * (num_pairs - ties_y));

    return tau;
}

//! tie adjustment for Kendall's test statistic
inline double ktau_stat_adjust(
    std::vector<double> x,
    std::vector<double> y,
    std::vector<double> weights)
{
    utils::check_sizes(x, y, weights);

    // 1.1 Sort x, y, and weights in x order; break ties in according to y.
    utils::sort_all(x, y, weights);

    // 1.2 Count pairs and triplets of tied x and simultaneous ties in x and y.
    double pair_x = utils::count_tied_pairs(x, weights);
    double trip_x = utils::count_tied_triplets(x, weights);
    double v_x = utils::count_ties_v(x, weights);

    // 2.1 Sort y and weights in y order; break ties according to x.
    utils::sort_all(y, x, weights);

    // 2.2 Count pairs and triplets of tied y.
    double pair_y = utils::count_tied_pairs(y, weights);
    double trip_y = utils::count_tied_triplets(y, weights);
    double v_y = utils::count_ties_v(y, weights);

    // 3. Calculate adjustment factor.
    if (weights.size() == 0)
        weights = std::vector<double>(x.size(), 1.0);
    double s = utils::sum(weights);
    double s2 = utils::perm_sum(weights, 2);
    double s3 = utils::perm_sum(weights, 3);
    double r = s / utils::sum(utils::pow(weights, 2));
    double v_0 = 2 * s2 * (2 * s) * std::pow(r, 3);
    double v_1 = 2 * pair_x * 2 * pair_y / (2 * 2 * s2) * std::pow(r, 2);
    double v_2 = 6 * trip_x * 6 * trip_y / (9 * 6 * s3) * std::pow(r, 3);
    double v = (v_0 - std::pow(r, 3) * (v_x - v_y)) / 18 + (v_1 + v_2);
    return std::pow(r, 2) * std::sqrt((s2 - pair_x) * (s2 - pair_y) / v);
}

}

}
