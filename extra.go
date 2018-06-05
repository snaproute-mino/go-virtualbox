package virtualbox

// SetExtra sets extra data. Name could be "global"|<uuid>|<vmname>
func SetExtra(name, key, val string) error {
	return VBM("setextradata", name, key, val)
}

// DelExtraData deletes extra data. Name could be "global"|<uuid>|<vmname>
func DelExtra(name, key string) error {
	return VBM("setextradata", name, key)
}
