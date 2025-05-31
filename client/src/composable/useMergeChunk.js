import axios from "axios"

export const useMergeChunk = () => {
    /**
     * @param {Object} params
     * @param {File} params.file
     * @param {Number} params.chunkSize
     * @param {AbortController} params.abortController 
     */
    async function mergeChunks({ file, chunkSize, abortController }) {
        try {
            const fileName = file.name
            const totalChunk = Math.ceil(file.size / chunkSize)

            const payload = { fileName, totalChunk }
            await axios.post("http://localhost:3000/merge-chunks", payload, {
                headers: {
                    "Content-Type": "application/json"
                },
                signal: abortController.signal
            })

            return "complete"

        } catch (error) {
            throw error
        }
    }

    return { mergeChunks }
}