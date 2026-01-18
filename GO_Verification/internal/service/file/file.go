package file

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"verif/app/internal/payload"
)

func CreateJsFile(verify *payload.Verification)error{
	file, errCr := os.Create("localDb.json")
	if errCr != nil{
		 return errCr
	}
	defer file.Close()
	data, errJs := json.Marshal(verify)
	if errJs != nil{
		return errJs
	}
	_, errWr := file.Write(data)
	if errWr != nil {
		return errWr
	}
	return nil
}
func CheckHash(hash string) error{
	data, errRead := os.ReadFile("localDb.json")
	if errRead != nil {
		 return errRead
	}
	var dataUser payload.Verification
	errJs := json.Unmarshal(data, &dataUser)
	if errJs != nil {
		return errJs
	}
	if dataUser.Hash != hash{
		return  errors.New("The specified hash does not match")
	}
	errRem := os.Remove("localDb.json")
	if errRem != nil {
		log.Println("failed to delete file: localDb.json")
	}
	return nil
}