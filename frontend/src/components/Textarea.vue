<template>
  <div class="textarea-group">
    <label v-if="label" :for="id" class="textarea-label">{{ label }}</label>
    <textarea
      :id="id"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :required="required"
      :rows="rows"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
      class="textarea"
    ></textarea>
    <div v-if="error" class="textarea-error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: () => `textarea-${Math.random().toString(36).substr(2, 9)}`
  },
  modelValue: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
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
  rows: {
    type: Number,
    default: 3
  },
  error: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])
</script>

<style scoped>
.textarea-group {
  margin-bottom: 16px;
}

.textarea-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s, box-shadow 0.3s;
  font-family: inherit;
  resize: vertical;
}

.textarea:focus {
  outline: none;
  border-color: #0366d6;
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
}

.textarea:disabled {
  background: #f8f9fa;
  border-color: #e2e8f0;
  color: #6c757d;
  cursor: not-allowed;
}

.textarea:readonly {
  background: #f8f9fa;
  border-color: #e2e8f0;
}

.textarea-error {
  margin-top: 4px;
  font-size: 12px;
  color: #dc3545;
}
</style>