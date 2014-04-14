package main

import (
	"fmt"
	"os"
	"time"
	"math/rand"
	"io"
)

func main() {
	args := os.Args;
	
	seeder := (time.Now().UTC().UnixNano())
	rand.Seed(seeder)
	
	if len(args) > 2 {
		if args[1] == "encrypt" {
			fmt.Println("encrypting")
			
			fi, err := os.Open(args[2])
			if err != nil{ panic(err) }
			
			defer func() {
				if err := fi.Close(); err != nil {
					panic(err)	
				}
			}()
			
			//open key file
			keyout, err := os.Create("key.txt")
			if err != nil { panic(err) }
			
			defer func() {
				if err := keyout.Close(); err != nil {
					panic(err)
				}
			}()
			
			//open encrypted file
			cryptout, err := os.Create("out.txt")
			if err != nil { panic(err) }
			
			defer func() {
				if err := cryptout.Close(); err != nil {
					panic(err)
				}
			}()
			
			buf := make([]byte, 1024)
			for {
				//read chunk
				n, err := fi.Read(buf)
				if err != nil && err != io.EOF { panic(err) }
				if n == 0 { break }
				
				var keybit = make([]byte,n)
				//gen key bit
				for i := 0; i < n; i++ {
					keybit[i] = byte(rand.Intn(126))
					buf[i]=buf[i]^keybit[i]
					buf[i]=buf[i]^26
					buf[i]=buf[i]^32
				}
				
				//write encrypted
				if _, err := keyout.Write(buf[:n]); err != nil {
					panic(err)
				}
				//write key
				if _, err := cryptout.Write(keybit[:n]); err != nil {
					panic(err)
				}
			}
			fmt.Println("done!")			
		} else if args[1] == "decrypt" {
			fmt.Println("trying to decrypt")
			
			fi, err := os.Open(args[2])
			if err != nil{ panic(err) }
			
			defer func() {
				if err := fi.Close(); err != nil {
					panic(err)	
				}
			}()
			
			keyin, err := os.Open(args[3])
			if err != nil{ panic(err) }
			
			defer func() {
				if err := keyin.Close(); err != nil {
					panic(err)	
				}
			}()
			
			//open decrypted file
			cryptout, err := os.Create("outDecrypted.txt")
			if err != nil { panic(err) }
			
			defer func() {
				if err := cryptout.Close(); err != nil {
					panic(err)
				}
			}()
			
			buf := make([]byte, 1024)
			keybit := make([]byte, 1024)
			for {
				//read chunk
				n, err := fi.Read(buf)
				if err != nil && err != io.EOF { panic(err) }
				if n == 0 { break }
				
				//read keybit
				n2, err2 := keyin.Read(keybit)
				if err2 != nil && err2 != io.EOF { panic(err2) }
				if n2 == 0 { break }
				
				for i := 0; i < n; i++ {
					buf[i] = 32^buf[i]
					buf[i] = 26^buf[i]
					buf[i] = keybit[i]^buf[i]
				}
				
				//write decrypted
				if _, err := cryptout.Write(buf[:n]); err != nil {
					panic(err)
				}
			}
			fmt.Println("done!")
			
		} else {
			fmt.Println("\nPlease use one of the following formats..\n-bitwise-crypt encrypt infile.txt\n-bitwise-crypt decrypt encrypted.txt key.txt")
		}
	} else {
		fmt.Println("\nPlease use one of the following formats..\n-bitwise-crypt encrypt infile.txt\n-bitwise-crypt decrypt encrypted.txt key.txt")
	}
}
