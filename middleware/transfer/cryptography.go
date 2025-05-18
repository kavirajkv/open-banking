package transfer

import (
	"github.com/kavirajkv/security/encrypt"
	"github.com/kavirajkv/security/sign"
	"github.com/kavirajkv/security/digest"

)

func Generatekey()(string,string,error){
	publickey,privatekey,err:=sign.GenerateKeypair()
	if err != nil {
		return "", "",err
	}
	return string(publickey),string(privatekey), nil
}


func EncryptData(key []byte, data string) (string,string, error) {
	encryptedData,nonce, err := encrypt.AESencrypt(key,data)
	if err != nil {
		return "", "",err
	}
	return string(encryptedData),string(nonce), nil
}


func DecryptData(key []byte,encrypted string,nonce string)string{
	decryprteddata:=encrypt.AESdecrypt(key,encrypted,nonce)
	return decryprteddata
}


func Createdigest(data string)(string){
	digest:=digest.ShaDigest(data)
	return digest
}


func SignData(key []byte, data string) (string, error){
	digest:=Createdigest(data)
	sign,err:=sign.Digitalsign(string(key),digest)
	if err != nil {
		return "", err
	}	
	return string(sign), nil
}


func VerifySign(key []byte, data string) (bool, error){
	digest:=Createdigest(data)
	sign,err:=sign.Verifysign(string(key),digest,data)
	if err != nil {
		return false, err
	}	
	return sign, nil
}