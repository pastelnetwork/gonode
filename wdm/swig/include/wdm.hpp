// Copyright © 2020 Thomas Nagler
//
// This file is part of the wdm library and licensed under the terms of
// the MIT license. For a copy, see the LICENSE file in the root directory
// or https://github.com/tnagler/wdm/blob/master/LICENSE.

#pragma once

#include "wdm/ktau.hpp"
#include "wdm/hoeffd.hpp"
#include "wdm/prho.hpp"
#include "wdm/srho.hpp"
#include "wdm/bbeta.hpp"
#include "wdm/methods.hpp"
#include "wdm/nan_handling.hpp"

//! Weighted dependence measures
namespace wdm {

//! calculates (weighted) dependence measures.
//! @param x, y input data.
//! @param method the dependence measure; see details for possible values.
//! @param weights an optional vector of weights for the data.
//! @param remove_missing if `true`, all observations containing a `nan` are
//!    removed; otherwise throws an error if `nan`s are present.
//!
//! @details
//! Available methods:
//!   - `"pearson"`, `"prho"`, `"cor"`: Pearson correlation
//!   - `"spearman"`, `"srho"`, `"rho"`: Spearman's \f$ \rho \f$
//!   - `"kendall"`, `"ktau"`, `"tau"`: Kendall's \f$ \tau \f$
//!   - `"blomqvist"`, `"bbeta"`, `"beta"`: Blomqvist's \f$ \beta \f$
//!   - `"hoeffding"`, `"hoeffd"`, `"d"`: Hoeffding's \f$ D \f$
//!
//! @return the dependence measure
inline double wdm(std::vector<double> x,
                  std::vector<double> y,
                  std::string method,
                  std::vector<double> weights = std::vector<double>(),
                  bool remove_missing = true)
{
    utils::check_sizes(x, y, weights);
    // na handling
    if (utils::preproc(x, y, weights, method, remove_missing) == "return_nan")
        return std::numeric_limits<double>::quiet_NaN();

    if (methods::is_hoeffding(method))
        return impl::hoeffd(x, y, weights);
    if (methods::is_kendall(method))
        return impl::ktau(x, y, weights);
    if (methods::is_pearson(method))
        return impl::prho(x, y, weights);
    if (methods::is_spearman(method))
        return impl::srho(x, y, weights);
    if (methods::is_blomqvist(method))
        return impl::bbeta(x, y, weights);
    throw std::runtime_error("method not implemented.");
}


//! Independence test
//!
//! The test calcualtes asymptotic p-values of independence tests based on
//! (weighted) dependence measures.
//!
//! @details
//! Available methods:
//!   - `"pearson"`, `"prho"`, `"cor"`: Pearson correlation
//!   - `"spearman"`, `"srho"`, `"rho"`: Spearman's \f$ \rho \f$
//!   - `"kendall"`, `"ktau"`, `"tau"`: Kendall's \f$ \tau \f$
//!   - `"blomqvist"`, `"bbeta"`, `"beta"`: Blomqvist's \f$ \beta \f$
//!   - `"hoeffding"`, `"hoeffd"`, `"d"`: Hoeffding's \f$ D \f$
//!
class Indep_test {
public:
    Indep_test() = delete;

    //! @param x, y input data.
    //! @param method the dependence measure; see class details for possible values.
    //! @param weights an optional vector of weights for the data.
    //! @param remove_missing if `true`, all observations containing a `nan` are
    //!    removed; otherwise throws an error if `nan`s are present.
    //! @param alternative indicates the alternative hypothesis and must be one
    //!    of `"two-sided"``, `"greater"` or `"less"`; `"greater"` corresponds
    //!    to positive association, `"less"` to negative association. For
    //!    Hoeffding's \f$ D \f$, only `"two-sided"` is allowed.
    Indep_test(std::vector<double> x,
               std::vector<double> y,
               std::string method,
               std::vector<double> weights = std::vector<double>(),
               bool remove_missing = true,
               std::string alternative = "two-sided") :
        method_(method),
        alternative_(alternative)
    {
        utils::check_sizes(x, y, weights);
        if (utils::preproc(x, y, weights, method, remove_missing) == "return_nan") {
            n_eff_ = utils::effective_sample_size(x.size(), weights);
            estimate_  = std::numeric_limits<double>::quiet_NaN();
            statistic_ = std::numeric_limits<double>::quiet_NaN();
            p_value_   = std::numeric_limits<double>::quiet_NaN();
        } else {
            n_eff_ = utils::effective_sample_size(x.size(), weights);
            estimate_ = wdm(x, y, method, weights, false);
            statistic_ = compute_test_stat(estimate_, method, n_eff_, x, y, weights);
            p_value_ = compute_p_value(statistic_, method, alternative, n_eff_);
        }
    }

    //! the method used for the test
    std::string method() const {return method_;}

    //! the alternative hypothesis used for the test
    std::string alternative() const {return alternative_;}

    //! the effective sample size in the test
    double n_eff() const {return n_eff_;}

    //! the estimated dependence measure
    double estimate() const {return estimate_;}

    //! the test statistic
    double statistic() const {return statistic_;}

    //! the p-value
    double p_value() const {return p_value_;}

private:

    inline double compute_test_stat(double estimate,
                                    std::string method,
                                    double n_eff,
                                    const std::vector<double>& x,
                                    const std::vector<double>& y,
                                    const std::vector<double>& weights)
    {
        // prevent overflow in atanh
        if (estimate == 1.0)
            estimate = 1 - 1e-12;
        if (estimate == -1.0)
            estimate = 1e-12;

        double stat;
        if (methods::is_hoeffding(method)) {
            stat = estimate / 30.0 + 1.0 / (36.0 * n_eff);
        } else if (methods::is_kendall(method)) {
            stat = estimate * impl::ktau_stat_adjust(x, y, weights);
        } else if (methods::is_pearson(method)) {
            stat = std::atanh(estimate) * std::sqrt(n_eff - 3);
        } else if (methods::is_spearman(method)) {
            stat = std::atanh(estimate) * std::sqrt((n_eff - 3) / 1.06);
        }  else if (methods::is_blomqvist(method)) {
            stat = std::atanh(estimate) * std::sqrt(n_eff);
        } else {
            throw std::runtime_error("method not implemented.");
        }

        return stat;
    }

    inline double compute_p_value(double statistic,
                                  std::string method,
                                  std::string alternative,
                                  double n_eff = 0.0)
    {
        double p_value;
        if (methods::is_hoeffding(method)) {
            if (n_eff == 0.0)
                throw std::runtime_error("must provide n_eff for method 'hoeffd'.");
            if (alternative != "two-sided")
                throw std::runtime_error("only two-sided test available for Hoeffding's D.");
            p_value = impl::phoeffb(statistic, n_eff);
        } else {
            if (alternative == "two-sided") {
                p_value = 2 * utils::normalCDF(-std::abs(statistic));
            } else if (alternative == "less") {
                p_value = utils::normalCDF(statistic);
            } else if (alternative == "greater") {
                p_value = 1 - utils::normalCDF(statistic);
            } else {
                throw std::runtime_error("alternative not implemented.");
            }
        }

        return p_value;
    }

    std::string method_;
    std::string alternative_;
    double n_eff_;
    double estimate_;
    double statistic_;
    double p_value_;
};


}
