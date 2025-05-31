<script setup>
</script>

<template>
    <button v-bind="$attrs" class="text-xs bg-white text-black rounded-md px-3 h-6 cursor-pointer hover:outline-2">
        <slot></slot>
    </button>
</template>