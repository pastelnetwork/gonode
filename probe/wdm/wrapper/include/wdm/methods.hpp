// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

namespace wdm {

namespace methods {

inline bool is_hoeffding(std::string method)
{
    return (method == "hoeffding") | (method == "hoeffd") | (method == "d");
}
inline bool is_kendall(std::string method)
{
    return (method == "kendall") | (method == "ktau") | (method == "tau");
}
inline bool is_pearson(std::string method)
{
    return (method == "pearson") | (method == "prho") | (method == "cor");
}
inline bool is_spearman(std::string method)
{
    return (method == "spearman") | (method == "srho") | (method == "rho");
}
inline bool is_blomqvist(std::string method)
{
    return (method == "blomqvist") | (method == "bbeta") | (method == "beta");
}

inline size_t get_min_nobs(std::string method)
{
    if (is_hoeffding(method)) {
        return 5;
    } else {
        return 2;
    }
}

}

}
