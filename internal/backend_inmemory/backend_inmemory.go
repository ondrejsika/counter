package backend_inmemory

var COUNTER = 0

func DoCountInMemory() (int, error) {
	COUNTER += 1
	return COUNTER, nil
}

func GetCountInMemory() (int, error) {
	return COUNTER, nil
}
