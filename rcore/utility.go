package rcore

// --------------------------------------------------------------------
//
func Contains(keys []string, key string) bool {
	//
	for _, v := range keys {
		//
		if v == key {
			return true
		}
	}

	//
	return false
}
