package repository

var globalURLStorage = map[string]string{}

func Create(mask string, originalURL string) string {
	globalURLStorage[mask] = originalURL
	return mask
}

func Get(mask string) (string, bool) {
	result, ok := globalURLStorage[mask]
	return result, ok
}

func Delete(mask string) {
	delete(globalURLStorage, mask)
}
