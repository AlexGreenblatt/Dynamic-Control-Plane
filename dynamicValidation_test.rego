package validate_test

test_valid_input if {
    # Define a valid input that should pass validation
    test_input := {
        "requestSchema": {
            "type": "object",
            "properties": {
                "serviceId": {"type": "string"},
                "verbose": {"type": "boolean"}
            },
            "required": ["serviceId"]
        },
        "requestBody": {
            "serviceId": "abc123",
            "verbose": true
        }
    }

    count(data.validate.violations) == 0 with input as test_input
}


test_invalid_input if {
    test_input := {
        "requestSchema": {
            "type": "object",
            "properties": {
                "serviceId": {"type": "string"},
                "verbose": {"type": "boolean"}
            },
            "required": ["serviceId"]
        },
        "requestBody": {
            # "serviceId" is missing
            "verbose": "true" # This is an invalid type (string instead of boolean)
        }
    }

    violations := data.validate.violations with input as test_input

    count(violations) == 2
}