// HIDEME: a tool for hiding large amounts of data on your computer using AES encryption
//
// notice: if someone does something bad with this its not my fault
// and also the hiding functionality is for really stupid people so someone might find out :|


package main

// probably a good idea to wipe this file afterwards
type hidemeConfig struct {
	randomFilenames  []string `json:"randomFilenames"`
	randomExtentions []string `json:"randomExtensions"`
	locationsTry     []string `json:"locationsTry"`
	split            int      `json:"split"`
	passphrase       string   `json:"passphrase"`
}


// load the config
