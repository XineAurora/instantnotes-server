package main

import (
	"log"

	"github.com/XineAurora/instantnotes-server/internal/app/instantnotesapi"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	s := instantnotesapi.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
