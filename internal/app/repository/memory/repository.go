package memory

var GlobalURLStorage = map[string]string{}

func Get(mask string) (string, bool) {
	result, ok := GlobalURLStorage[mask]
	return result, ok
}

func Delete(mask string) {
	delete(GlobalURLStorage, mask)
}

func CreateInMemory(mask string, originalURL string) {
	GlobalURLStorage[mask] = originalURL
}
