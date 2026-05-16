package main

import (
    "crypto/aes"
    "crypto/cipher"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

const aes_key = "12345678123456781234567812345678"

func main() {
    path := "pdfs"
    
    err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println("Error accesing path: ", path, err)
            return err
        }
        
        // Brišemo ransom note prilikom dekripcije
        if !info.IsDir() && info.Name() == "RANSOM_NOTE.txt" {
            os.Remove(path)
            fmt.Println("Ransom poruka obrisana.")
            return nil
        }

        if !info.IsDir() && strings.HasSuffix(info.Name(), ".enc") {
            data, err := os.ReadFile(path)
            if err != nil {
                fmt.Println("Error reading file: ", path, err)
                return err
            }
            newFileName := strings.TrimSuffix(path, ".enc")
            file, err := os.Create(newFileName)
            if err != nil {
                fmt.Println(err)
                return err
            }
            defer file.Close()

            aes_key_byted := []byte(aes_key)
            aes_key_cipher, _ := aes.NewCipher(aes_key_byted)
            gcm, err := cipher.NewGCM(aes_key_cipher)
            if err != nil {
                fmt.Println("Greska pri GCM")
                return err
            }
            nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
            plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
            if err != nil {
                fmt.Println("Greska pri dekripciji fajla: ", path)
                return err
            }

            _, err = file.Write(plainText)
            if err != nil {
                fmt.Println("Greska pri pisanju")
                return err
            }
            err = os.Remove(path)
            if err != nil {
                fmt.Println("Greska pri brisanju .enc fajla")
                return err
            }
            fmt.Printf("Dekriptovan: %s\n", newFileName)
        }
        return err
    })

    if err != nil {
        fmt.Println("Error walking through the files/folders", err)
        return
    }
    fmt.Println("Dekripcija zavrsena. Fajlovi povraceni.")
}