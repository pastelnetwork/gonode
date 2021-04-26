// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

#include "utils.hpp"
#include "ranks.hpp"
#include <cmath>

namespace wdm {

namespace impl {

const double pi = std::acos(-1);

//! fast calculation of the weighted Hoeffdings's D.
//! @param x, y input data.
//! @param weights an optional vector of weights for the data.
inline double hoeffd(std::vector<double> x,
                     std::vector<double> y,
                     std::vector<double> weights = std::vector<double>())
{
    utils::check_sizes(x, y, weights);

    // 1. Compute (weighted) ranks
    std::vector<double> R_X = rank(x, weights);
    std::vector<double> R_Y = rank(y, weights);
    std::vector<double> S_X, S_Y;
    if (weights.size() > 0) {
        S_X = rank(x, utils::pow(weights, 2));
        S_Y = rank(y, utils::pow(weights, 2));
    } else {
        S_X = R_X;
        S_Y = R_Y;
    }

    // 2. Compute (weighted) bivariate ranks (number of points w/ both columns
    // less than the ith row).
    std::vector<double> R_XY, S_XY, T_XY, U_XY;
    R_XY = bivariate_rank(x, y, weights);
    if (weights.size() > 0) {
        S_XY = bivariate_rank(x, y, utils::pow(weights, 2));
        T_XY = bivariate_rank(x, y, utils::pow(weights, 3));
        U_XY = bivariate_rank(x, y, utils::pow(weights, 4));
    } else {
        S_XY = R_XY;
        T_XY = R_XY;
        U_XY = R_XY;
    }


    // 3. Compute (weighted) Hoeffdings' D
    if (weights.size() == 0)
        weights = std::vector<double>(x.size(), 1.0);
    double A_1 = 0.0, A_2 = 0.0, A_3 = 0.0;
    for (size_t i = 0; i < x.size(); i++) {
        A_1 += (R_XY[i] * R_XY[i] - S_XY[i]) * weights[i];
        A_2 += (
            (R_X[i] * R_Y[i] - S_XY[i]) * R_XY[i] -
                S_XY[i] * (R_X[i] + R_Y[i]) + 2 * T_XY[i]
        ) * weights[i];
        A_3 += (
            (R_X[i] * R_X[i] - S_X[i]) * (R_Y[i] * R_Y[i] - S_Y[i])  -
                4 * ((R_X[i] * R_Y[i] - S_XY[i]) * S_XY[i] -
                T_XY[i] * (R_X[i] + R_Y[i]) + 2 * U_XY[i]) -
                2 * (S_XY[i] * S_XY[i] - U_XY[i])
        ) * weights[i];
    }
    double D = 0.0;
    D += A_1 / (utils::perm_sum(weights, 3) * 6);
    D -= 2 * A_2 / (utils::perm_sum(weights, 4) * 24);
    D += A_3 / (utils::perm_sum(weights, 5) * 120);

    return 30.0 * D;
}

//! calculates the (approximate) asymptotic distribution function of Hoeffding's
//! B (as in Blum, Kiefer, and Rosenblatt) under the null hypothesis of
//! independence.
//! @param B sample estimate of Hoeffding's B.
//! @param n the sample size.
inline double phoeffb(double B, double n) {
    B *= 0.5 * std::pow(pi, 4) * (n - 1);

    // obtain approximate p values by interpolation of tabulated values
    double p;
    if ((B <= 1.1) | (B >= 8.5)) {
        p = std::min(1.0, std::exp(0.3885037 - 1.164879 * B));
        p = std::max(1e-12, p);
    } else {
        std::vector<double> grid{
            1.1, 1.15, 1.2, 1.25, 1.3, 1.35, 1.4, 1.45, 1.5, 1.55, 1.6,
            1.65, 1.7, 1.75, 1.8, 1.85, 1.9, 1.95, 2, 2.05, 2.1, 2.15, 2.2,
            2.25, 2.3, 2.35, 2.4, 2.45, 2.5, 2.55, 2.6, 2.65, 2.7, 2.75,
            2.8, 2.85, 2.9, 2.95, 3, 3.05, 3.1, 3.15, 3.2, 3.25, 3.3, 3.35,
            3.4, 3.45, 3.5, 3.55, 3.6, 3.65, 3.7, 3.75, 3.8, 3.85, 3.9, 3.95,
            4, 4.05, 4.1, 4.15, 4.2, 4.25, 4.3, 4.35, 4.4, 4.45, 4.5, 4.55,
            4.6, 4.65, 4.7, 4.75, 4.8, 4.85, 4.9, 4.95, 5, 5.5, 6, 6.5, 7,
            7.5, 8, 8.5
        };
        std::vector<double> vals{
            0.5297, 0.4918, 0.4565, 0.4236, 0.3930, 0.3648, 0.3387, 0.3146,
            0.2924, 0.2719, 0.2530, 0.2355, 0.2194, 0.2045, 0.1908, 0.1781,
            0.1663, 0.1554, 0.1453, 0.1359, 0.1273, 0.1192, 0.1117, 0.1047,
            0.0982, 0.0921, 0.0864, 0.0812, 0.0762, 0.0716, 0.0673, 0.0633,
            0.0595, 0.0560, 0.0527, 0.0496, 0.0467, 0.0440, 0.0414, 0.0390,
            0.0368, 0.0347, 0.0327, 0.0308, 0.0291, 0.0274, 0.0259, 0.0244,
            0.0230, 0.0217, 0.0205, 0.0194, 0.0183, 0.0173, 0.0163, 0.0154,
            0.0145, 0.0137, 0.0130, 0.0123, 0.0116, 0.0110, 0.0104, 0.0098,
            0.0093, 0.0087, 0.0083, 0.0078, 0.0074, 0.0070, 0.0066, 0.0063,
            0.0059, 0.0056, 0.0053, 0.0050, 0.0047, 0.0045, 0.0042, 0.00025,
            0.00014, 0.0008, 0.0005, 0.0003, 0.0002, 0.0001
        };
        p = utils::linear_interp(B, grid, vals);
    }

    return p;
}

}

}
