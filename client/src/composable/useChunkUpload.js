import axios from "axios"
import { sleep } from "../lib/utils"

export const useChunkUpload = () => {
    /**
     * Upload file in chunks to the server.
     * 
     * @param {Object} params
     * @param {File} params.file
     * @param {String} params.fileName
     * @param {Number} params.fileSize
     * @param {Number} params.currentChunkIndex
     * @param {Number} params.chunkSize
     * @param {Number} params.totalChunk
     * @param {AbortController} params.abortController
     * @param {function ({currentChunkIndex: number, progress: number})} params.uploadProgress
     * 
     * @returns {Promise<string>} 
     */

    async function chunkUpload(params) {
        console.log(params)
        try {
            for (let i = params.currentChunkIndex; i < params.totalChunk; i++) {
                const
                    start = i * params.chunkSize,
                    end = Math.min(start + params.chunkSize, params.fileSize),
                    chunk = params.file.slice(start, end)

                const formdata = new FormData()
                formdata.append("chunk", chunk)

                await axios.post("http://localhost:3000/upload", formdata, {
                    params: {
                        fileName: params.fileName,
                        chunkIndex: i,
                    },
                    headers: {
                        "Content-Type": "multipart/form-data"
                    },
                    signal: params.abortController.signal
                })

                params.uploadProgress({
                    currentChunkIndex: i,
                    progress: Math.ceil(Math.abs((i + 1) / params.totalChunk) * 100)
                })

                if (i == params.totalChunk - 1) {
                    return "complete"
                }

                await sleep(400)
            }

        } catch (error) {
            throw error
        }
    }

    return { chunkUpload }
}