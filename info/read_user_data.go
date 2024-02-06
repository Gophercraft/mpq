package info

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const MaxUserDataSize = 0x17FF

func ReadUserData(reader io.Reader, user_data *UserData) (err error) {
	// read user data size
	var user_data_size uint32
	if err = binary.Read(reader, binary.LittleEndian, &user_data_size); err != nil {
		return err
	}
	if user_data_size > MaxUserDataSize {
		return fmt.Errorf("info: user data size is too large (%d)", user_data)
	}
	// read user data
	user_data_bytes := make([]byte, user_data_size)
	if _, err = io.ReadFull(reader, user_data_bytes); err != nil {
		return
	}
	user_data_reader := bytes.NewReader(user_data_bytes)
	// read first segment of user data
	if err = binary.Read(user_data_reader, binary.LittleEndian, &user_data.UserDataPrefix); err != nil {
		return
	}
	return
}
