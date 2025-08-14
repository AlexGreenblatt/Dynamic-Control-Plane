package validate

# This rule generates a set of violations by checking required fields and data types.
violations contains message if {
    # Check for required fields
    some required_field in input.requestSchema.required
    not input.requestBody[required_field]
    message := sprintf("'%s' is a required field.", [required_field])
}

violations contains message if {
    # Check data types for all properties in the schema
    some property_name, property_def in input.requestSchema.properties
    value := input.requestBody[property_name]

    # Skip validation if the data isn't present
    not is_null(value)

    # Check if the value's type matches the schema's type
    not type_matches(value, property_def.type)
    message := sprintf("'%s' has an invalid type. Expected: %s", [property_name, property_def.type])
}

# Helper function to check if a value's type matches a given string type
type_matches(val, type) if {
    type == "string"
    is_string(val)
}

type_matches(val, type) if {
    type == "number"
    is_number(val)
}

type_matches(val, type) if {
    type == "boolean"
    is_boolean(val)
}