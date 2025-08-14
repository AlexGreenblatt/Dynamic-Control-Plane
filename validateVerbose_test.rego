package validate_test

import future.keywords.if
import future.keywords.in

import data.validate
import data.test.assert

test_allow_with_verbose_true if {
    validate.allow with input as {"verbose": true}
}

test_allow_with_verbose_false if {
    validate.allow with input as {"verbose": false}
}

test_no_violation_when_missing if {
    v := validate.violations with input as {}
    count(validate.violations) == 0
}

test_violation_when_verbose_is_neither_true_nor_fale if {
    v := validate.violations with input as {"verbose": "a non-boolean value"}
    some i
    v[i] == "verbose must be true or false"
}

