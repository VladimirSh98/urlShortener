package memory

// GlobalURLStorage variable to store urls
var GlobalURLStorage = map[string]string{}

// Get from memory url by mask
func Get(mask string) (string, bool) {
	result, ok := GlobalURLStorage[mask]
	return result, ok
}

// Delete from memory url by mask
func Delete(mask string) {
	delete(GlobalURLStorage, mask)
}

// CreateInMemory create record in memory for new url
func CreateInMemory(mask string, originalURL string) {
	GlobalURLStorage[mask] = originalURL
}
