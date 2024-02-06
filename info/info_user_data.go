package info

type UserData struct {
	UserDataPrefix
}

type UserDataPrefix struct {
	// Offset of the MPQ header, relative to the begin of this header
	HeaderOffset uint32
	// Appears to be size of user data header (Starcraft II maps)
	UserDataHeader uint32
}
