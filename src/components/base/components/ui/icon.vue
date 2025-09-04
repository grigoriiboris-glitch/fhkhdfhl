<template>
  <i :class="name" :style="{ fontSize: `${size}px`, color: color }"></i>
</template>

<script>
export default {
  props: {
    name: {
      required: true,
    },
    color: {
      defautl: null,
    },
    size: {
      type: Number,
      default: 16,
    },
  },
};
</script>

<style scoped></style>
