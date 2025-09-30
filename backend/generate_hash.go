package main

  import (
        "fmt"
        "golang.org/x/crypto/bcrypt"
  )

  func main() {
        // ***********************************
        // * YOU MUST CHANGE THE TEXT INSIDE THE QUOTES BELOW            *
        // * Replace "YOUR_CHOSEN_PASSWORD" with the actual password     *
        // * you want for your Reception Manager.                        *
        // *** For example, if you want the password to be "MySecretPass123",
        // *** change the line to: password := "MySecretPass123"
        // ***********************************
        password := "testpass123" // <--- CHANGE THIS LINE

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
                fmt.Println("Error hashing password:", err)
                return
        }
        fmt.Println("Hashed Password:", string(hashedPassword))
  }
  