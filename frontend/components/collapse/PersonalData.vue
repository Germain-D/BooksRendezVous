<template>
  <details class="collapse bg-base-200 mb-4" open>
    <summary class="collapse-title text-xl font-medium p-6">
      Informations personnelles
    </summary>
    <div class="collapse-content px-6">
      <div class="space-y-4">
        <div>
          <label class="label">Pseudo</label>
          <input
            type="text"
            v-model="pseudo"
            class="input input-bordered w-full"
            disabled
          />
        </div>
        <div>
          <label class="label">Email</label>
          <input
            type="email"
            v-model="email"
            class="input input-bordered w-full"
            disabled
          />
        </div>
      </div>
    </div>
  </details>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useAuthStore } from "@/stores/auth";

const authStore = useAuthStore();
const config = useRuntimeConfig();

// Profile data
const pseudo = ref(authStore.getPseudo);
const email = ref(authStore.getEmail);
</script>

<style scoped>
.collapse {
  @apply rounded-lg overflow-hidden transition-all duration-200;
}

.collapse-title {
  @apply cursor-pointer hover:bg-base-300 transition-colors;
}

.collapse-content {
  @apply pt-0;
}

.collapse[open] summary {
  @apply border-b border-base-300;
}
</style>
