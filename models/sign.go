package models


import (
    "crypto/hmac"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/base64"
    "fmt"
)




func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ComputeHmac512(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha512.New, key)
    h.Write([]byte(message))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}


func IsCorrectSign(userInfo *UserInfo,params map[string]string)bool{

    return true
}

func CreatSignString(){

}


func test() {
    fmt.Println(ComputeHmac256("Message", "secret"))
}

