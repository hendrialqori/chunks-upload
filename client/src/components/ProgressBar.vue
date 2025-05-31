<script setup>
const props = defineProps({
    progress: Number
})
</script>

<template>
    <div class="flex items-center gap-3" aria-label="progress">
        <div class="w-full bg-white rounded-lg h-1">
            <div :style="`width: ${props.progress}%`" class="h-full bg-blue-500 rounded-md"></div>
        </div>
        <span class="text-xs font-medium text-gray-600">{{ props.progress }}%</span>
    </div>
</template>