package validate


allow if {
	count(violations) == 0
}


violations contains "verbose must be true or false" if {
	input.verbose != true
    input.verbose != false
}