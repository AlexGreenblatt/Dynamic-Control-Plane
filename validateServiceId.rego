package validate


allow if {
	count(violations) == 0
}


violations contains "serviceId is a required field" if {
	not input.serviceId
}

violations contains "serviceId has a maximum length of 32" if {
	count(input.serviceId) >= 33
}