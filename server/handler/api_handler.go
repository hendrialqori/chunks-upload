package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func UploadFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// create uploads directory
	os.MkdirAll("./uploads", os.ModePerm)
	// create temporary directory
	os.MkdirAll("./temp", os.ModePerm)

	// take fileName query params from request
	fileName := r.URL.Query().Get("fileName")
	// take chunkIndexString query params from request
	chunkIndexString := r.URL.Query().Get("chunkIndex")
	// parse into int
	chunkIndex, err := strconv.Atoi(chunkIndexString)
	if err != nil {
		panic(err)
	}

	// take chunks that sending from client
	file, _, err := r.FormFile("chunk")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tempFile := filepath.Join("./temp", fmt.Sprintf("%s.part%d", fileName, chunkIndex))
	// Prepare empty file fo save chunks
	outputFile, err := os.Create(tempFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// create buffer of 4MB, to save temporary data while copying
	buf := make([]byte, 4*1024*1024)
	// copying file (containng chunk from client) to outputFile (temp/chunk.partN) using that buffer
	_, err = io.CopyBuffer(outputFile, file, buf)
	if err != nil {
		// If it fails while copying (for example, full disk, corrupt file, etc.), stop the program and show the error
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}

func MergeChunks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var fileRequest struct {
		FileName    string `json:"fileName"`
		TotalChunks int    `json:"totalChunk"`
	}

	if err := json.NewDecoder(r.Body).Decode(&fileRequest); err != nil {
		panic("invalid request")
	}

	outPath := filepath.Join("./uploads", fileRequest.FileName)
	outFile, err := os.Create(outPath)
	if err != nil {
		panic("failed to create merged file")
	}
	defer outFile.Close()

	for i := range fileRequest.TotalChunks {
		chunkPath := filepath.Join("./temp", fmt.Sprintf("%s.part%d", fileRequest.FileName, i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			panic("failed to open chunk")
		}

		_, err = io.Copy(outFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			panic("failed to write chunk")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "Chunks merged successfully",
	})

}
