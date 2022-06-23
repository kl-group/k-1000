package main

var needReloadAsterisk = false

func main() {
	entryFlags()
	entryLog()
	entryConfig()
	entryLdap()
	defer GConn.Close()
	//entryQL()
	//entryExportFromLdap()
	err := VarMapExtensionLdap.loadFromLdap()
	if err != nil {
		logme.Fatal(err)
	}

	err = VarMapExtensionAsterisk.loadFromAsterisk()
	if err != nil {
		logme.Fatal(err)
	}
	VarMapExtensionLdap.prepareLdap()
	VarMapExtensionAsterisk.prepareAsterisk()

	reloadAsterisk()
}
