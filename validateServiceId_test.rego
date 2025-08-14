package validate_test

import future.keywords.if
import future.keywords.in

import data.validate
import data.test.assert

test_allow_with_serviceId if {
    validate.allow with input as {"serviceId": "abc"}
}

test_violation_when_missing if {
    v := validate.violations with input as {}
    some i
    v[i] == "serviceId is a required field"
}

test_violation_when_serviceId_is_too_long if {
    v := validate.violations with input as {"serviceId": "a string that is longer than allowed"}
    some i
    v[i] == "serviceId has a maximum length of 32"
}

