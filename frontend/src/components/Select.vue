<template>
  <div class="select-group">
    <label v-if="label" :for="id" class="select-label">{{ label }}</label>
    <select
      :id="id"
      :value="modelValue"
      :disabled="disabled"
      :required="required"
      @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
      class="select"
    >
      <option v-if="placeholder" :value="''">{{ placeholder }}</option>
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label }}
      </option>
    </select>
    <div v-if="error" class="select-error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">

interface Option {
  value: string | number
  label: string
}

const props = defineProps({
  id: {
    type: String,
    default: () => `select-${Math.random().toString(36).substr(2, 9)}`
  },
  modelValue: {
    type: [String, Number],
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  options: {
    type: Array as () => Option[],
    default: () => []
  },
  placeholder: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])
</script>

<style scoped>
.select-group {
  margin-bottom: 16px;
}

.select-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s, box-shadow 0.3s;
  font-family: inherit;
  background: #fff;
  cursor: pointer;
}

.select:focus {
  outline: none;
  border-color: #0366d6;
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
}

.select:disabled {
  background: #f8f9fa;
  border-color: #e2e8f0;
  color: #6c757d;
  cursor: not-allowed;
}

.select-error {
  margin-top: 4px;
  font-size: 12px;
  color: #dc3545;
}
</style>