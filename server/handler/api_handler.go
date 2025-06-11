package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/julienschmidt/httprouter"
)

func HTTPError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  status,
		"message": message,
	})
}

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
		HTTPError(w, http.StatusBadRequest, "invalid chunk index")
		return
	}

	// take chunks that sending from client
	file, _, err := r.FormFile("chunk")
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "failed to get file from request")
		return
	}
	defer file.Close()

	tempFile := filepath.Join("./temp", fmt.Sprintf("%s.part%d", fileName, chunkIndex))
	// Prepare empty file fo save chunks
	outputFile, err := os.Create(tempFile)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, "failed to create temporary file")
		return
	}
	defer outputFile.Close()

	// create buffer of 4MB, to save temporary data while copying
	buf := make([]byte, 4*1024*1024)
	// copying file (containng chunk from client) to outputFile (temp/chunk.partN) using that buffer
	_, err = io.CopyBuffer(outputFile, file, buf)
	if err != nil {
		// If it fails while copying (for example, full disk, corrupt file, etc.), stop the program and show the error
		HTTPError(w, http.StatusInternalServerError, "failed to write chunk to temporary file")
		return
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
		HTTPError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	outPath := filepath.Join("./uploads", fileRequest.FileName)
	outFile, err := os.Create(outPath)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, "failed to create output file")
		return
	}
	defer outFile.Close()

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := range fileRequest.TotalChunks {
		wg.Add(1)

		go func(chunkIndex int) {
			defer wg.Done()
			chunkPath := filepath.Join("./temp", fmt.Sprintf("%s.part%d", fileRequest.FileName, chunkIndex))
			chunkFile, err := os.Open(chunkPath)
			if err != nil {
				if os.IsNotExist(err) {
					HTTPError(w, http.StatusBadRequest, fmt.Sprintf("chunk %d does not exist", i))
					return
				}
				HTTPError(w, http.StatusInternalServerError, "failed to open chunk file")
				return
			}
			defer chunkFile.Close()

			chunkData, err := io.ReadAll(chunkFile)
			if err != nil {
				HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to read chunk %d file", i))
				return
			}

			mu.Lock()
			defer mu.Unlock()

			_, err = outFile.Write(chunkData)
			if err != nil {
				HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write chunk %d to output file", i))
				return
			}

			os.Remove(chunkPath) // Remove the chunk file after merging
		}(i)
	}

	wg.Wait()

	if err := CleanupTempFiles(w); err != nil {
		HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to clean up temporary files: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "Chunks merged successfully",
	})

}

func CleanupTempFiles(w http.ResponseWriter) error {
	if _, err := os.Stat("./temp"); os.IsNotExist(err) {
		return nil // No temp directory to clean up
	}

	files, err := filepath.Glob("./temp/*.part*")
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to list temporary files: %v", err))
		return err
	}

	for i, file := range files {
		err := os.Remove(file)
		if err != nil {
			HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to remove temporary file-%d: %v\n", i, err))
			return err
		}
	}

	return nil
}
