<script setup>
import { FileUp, File, CheckLine } from 'lucide-vue-next';
import { ref, computed } from 'vue'
import axios from 'axios'
import { useChunkUpload } from './composable/useChunkUpload';
import { useMergeChunk } from './composable/useMergeChunk'
import { status } from './model/uploadStatus'
import { toMB } from './lib/utils'
import ProgressBar from './components/ProgressBar.vue';
import ButtonUpload from './components/ButtonUpload.vue';

const { chunkUpload } = useChunkUpload()
const { mergeChunks } = useMergeChunk()

const CHUNK_SIZE = 1 * 1024 * 1024

// {
//   "fileName": {
//     isPause: false
//     file: null,
//      status: pending | process | success | failed
//       currentChunkIndex: 0,
//         progress: 0
//   }
// }
const files = ref({})

const filesCollection = computed(() => {
  let temp = []
  for (const [key, value] of Object.entries(files.value ?? [])) {
    temp = [...temp, { name: key, ...value }]
  }
  return temp
})

const handleUploadFile = (event) => {
  const file = event.target.files[0]
  files.value[file.name] = {
    file,
    status: status.PENDING,
    progress: 0,
    currentChunkIndex: 0,
    abortController: new AbortController()
  }
}

const chunkUploadAPI = async (fileKey) => {
  try {
    const fileState = files.value[fileKey]
    if (!fileState) return

    if (fileState.status == status.PAUSE) {
      fileState.abortController = new AbortController()
    }

    const fileName = fileState.file.name
    const fileSize = fileState.file.size
    const currentChunkIndex = fileState.currentChunkIndex
    const abortController = fileState.abortController
    const totalChunk = Math.ceil(fileSize / CHUNK_SIZE)

    fileState.status = status.FETCHING

    // do chunk upload first
    const resChunk = await chunkUpload({
      file: fileState.file,
      fileName,
      fileSize,
      currentChunkIndex,
      abortController,
      chunkSize: CHUNK_SIZE,
      totalChunk,
      uploadProgress: ({ currentChunkIndex, progress }) => {
        fileState.currentChunkIndex = currentChunkIndex
        fileState.progress = progress
      }
    })
    if (resChunk != "complete") return

    // then, if chunk upload not error do merge all chunks
    const resMerge = await mergeChunks({
      file: fileState.file,
      chunkSize: CHUNK_SIZE,
      abortController
    })

    if (resMerge == "complete") {
      fileState.status = status.SUCCESS
    }

  } catch (error) {
    if (axios.isAxiosError(error)) {
      const cancelError = "ERR_CANCELED"
      if (error.code = cancelError) {
        console.log(`Upload file ${fileKey} pause!`)
      } else {
        console.log(error)
        alert("error while pipe process")
      }
    }
  }
}

const pauseUpload = (fileKey) => {
  const fileState = files.value[fileKey]
  fileState.abortController.abort()
  fileState.status = status.PAUSE
}

</script>

<template>
  <main class="flex items-center justify-center min-h-screen w-full bg-[#F3F5F7]">
    <div class="min-w-[500px] rounded-xl shadow p-4 bg-white flex flex-col gap-5">
      <h1 class="font-medium">
        Upload file
      </h1>
      <label>
        <div
          class="py-8 flex flex-col items-center justify-center gap-3 border-2 border-dashed bg-blue-400/5 hover:bg-blue-400/20 border-blue-400 transition duration-300 cursor-pointer rounded-lg"
          aria-label="area">
          <FileUp class="text-blue-500 size-7" />
          <p class="text-sm text-gray-500">Upload file here</p>
        </div>
        <input type="file" accept="*" class="hidden" @change="handleUploadFile">
      </label>
      <div class="space-y-3">
        <figure v-for="(collection, index) of filesCollection" :key="index"
          class="bg-gray-100 rounded-md px-3 py-3 space-y-2">
          <div class="flex justify-between">
            <div class="flex gap-3">
              <div class="p-2 bg-white rounded-md w-max">
                <File class="size-4" />
              </div>
              <div>
                <h2 class="text-[0.75rem] font-medium">{{ collection.name }}</h2>
                <p class="text-xs text-gray-500">{{ toMB(collection.file.size) }} MB</p>
              </div>
            </div>
            <ButtonUpload v-if="collection.status === status.FETCHING" @click="pauseUpload(collection.name)">
              Pause
            </ButtonUpload>
            <ButtonUpload v-if="[status.PENDING, status.PAUSE].includes(collection.status)"
              @click="chunkUploadAPI(collection.name)">
              {{ collection.status == status.PAUSE ? "Continue" : "Upload" }}
            </ButtonUpload>
            <div v-if="collection.status == status.SUCCESS" class="inline-flex gap-2">
              <CheckLine class="text-green-500 size-4" />
            </div>
          </div>
          <ProgressBar v-if="collection.status != status.SUCCESS" :progress="collection.progress ?? 0" />
        </figure>
      </div>
    </div>
  </main>
</template>