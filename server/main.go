package main

import (
	"chunks-upload-server/handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	var route = httprouter.New()

	route.POST("/upload", handler.UploadFile)
	route.POST("/merge-chunks", handler.MergeChunks)

	route.PanicHandler = func(w http.ResponseWriter, r *http.Request, i any) {
		fmt.Println(i)

		apiResponse := map[string]any{
			"status":  http.StatusInternalServerError,
			"message": i,
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
			panic(err)
		}
	}

	handler := cors.Default().Handler(route)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: handler,
	}

	log.Printf("Server is running on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
