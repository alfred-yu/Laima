<template>
  <div class="input-group">
    <label v-if="label" :for="id" class="input-label">{{ label }}</label>
    <input
      :id="id"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :required="required"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      class="input"
    />
    <div v-if="error" class="input-error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: () => `input-${Math.random().toString(36).substr(2, 9)}`
  },
  modelValue: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'text',
    validator: (value: string) => ['text', 'password', 'email', 'number', 'tel', 'url'].includes(value)
  },
  placeholder: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  readonly: {
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
.input-group {
  margin-bottom: 16px;
}

.input-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s, box-shadow 0.3s;
  font-family: inherit;
}

.input:focus {
  outline: none;
  border-color: #0366d6;
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
}

.input:disabled {
  background: #f8f9fa;
  border-color: #e2e8f0;
  color: #6c757d;
  cursor: not-allowed;
}

.input:readonly {
  background: #f8f9fa;
  border-color: #e2e8f0;
}

.input-error {
  margin-top: 4px;
  font-size: 12px;
  color: #dc3545;
}
</style>