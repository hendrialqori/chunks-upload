# Chunked File Upload using Go 
![Chunk upload preview](assets/chunks-upload.gif)

Demonstrates how to upload large files in chunks using a Go backend and Vue frontend. The backend is built with `net/http` and `httprouter`, and the frontend uses plain `Vue` with `Axios` for HTTP requests.

## 📦 Tech Stack

- **Go**: Handles chunk upload and file merge logic.
- **httprouter**: Lightweight HTTP router for Go.
- **Vue**: Frontend file selection and chunked upload.
- **Axios**: Makes HTTP requests to the backend.


## 🚀 Features
- Upload large files in chunks (default 1MB).
- Multiple upload
- Resumable uploads (pause/continue).
- Automatic file merge after upload at folder `/uploads` (server)
- Supports progress tracking.

## 🚀 API Endpoints

### POST /upload

Upload a single file chunk.  
Send the chunk as **form-data** with the key: `chunk`.

Required query parameters:  
- `fileName` — original file name (e.g., `document.pdf`)  
- `chunkIndex` — chunk index (starting from 0)

---

### POST /merge-chunks

Merge all uploaded chunks into a complete file.

Request body (JSON):

```json
{
  "fileName": "filename.pdf",
  "totalChunk": 10
}
```

## 🧑‍💻 How to Clone and Run

### 1. Clone the Repository

```
git clone https://github.com/hendrialqori/chunks-upload.git

cd chunks-upload
```
### 2. Run the Backend (Go Server)

```
cd server

go mod tidy

make go
```

### 3. Run the Frontend (Vue Client)

```
cd client

npm install

npm run dev
```

🌐 Visit the app at: http://localhost:5173
🔗 Ensure the API calls in the client point to http://localhost:3000

## MIT License
Please feel free to use for personal and commercial purposes.
