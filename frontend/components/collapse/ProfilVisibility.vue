<template>
  <div v-if="showCopyAlert" role="alert" class="alert alert-success mt-4">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-6 w-6 shrink-0 stroke-current"
      fill="none"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
    <span>Lien copié avec succès !</span>
  </div>
  <details class="collapse bg-base-200 mb-4" open>
    <summary class="collapse-title text-xl font-medium p-6">
      Visibilité du profil
    </summary>
    <div class="collapse-content px-6">
      <div class="flex gap-4 items-center">
        <h2 class="text-xl font-bold">Profil public</h2>
        <div class="inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            id="favorite"
            v-model="isPublic"
            @change="handleProfileVisibility"
            class="toggle toggle-primary"
          />
          <div
            class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300"
          >
            {{
              isPublic
                ? "Votre profil est visible publiquement"
                : "Votre profil est privé"
            }}
          </div>
        </div>
      </div>

      <div v-if="isPublic" class="mt-4">
        <label class="label">Lien de partage</label>
        <div class="flex gap-4">
          <input
            type="text"
            :value="sharelink"
            class="input input-bordered flex-1"
            readonly
          />
          <button @click="copy" class="btn btn-secondary">Copier</button>
        </div>
      </div>
    </div>
  </details>
  <!-- Confirmation Modal -->
  <div
    v-if="showConfirmation"
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-4"
  >
    <div role="alert" class="alert max-w-lg bg-base-200">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="stroke-info h-6 w-6 shrink-0"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
      <span>{{ confirmationMessage }}</span>
      <div>
        <button class="btn btn-sm" @click="cancelVisibilityChange">
          Annuler
        </button>
        <button
          class="btn btn-sm btn-primary ml-2"
          @click="confirmVisibilityChange"
        >
          Confirmer
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useAuthStore } from "@/stores/auth";

const authStore = useAuthStore();
const config = useRuntimeConfig();

// Profile data
const pseudo = ref(authStore.getPseudo);
const email = ref(authStore.getEmail);
const isPublic = ref(false);
const sharelink = ref(
  `${config.public.BASE_URL}/publicprofile/${authStore.getPseudo}_${authStore.getShareLink}` ||
    ""
);

// Confirmation modal
const showConfirmation = ref(false);
const confirmationMessage = ref("");
const pendingVisibility = ref(false);

const showCopyAlert = ref(false);

const handleProfileVisibility = async () => {
  pendingVisibility.value = isPublic.value;
  confirmationMessage.value = isPublic.value
    ? "Êtes-vous sûr de vouloir rendre votre profil public ? Vos lectures seront visibles par tous."
    : "Êtes-vous sûr de vouloir rendre votre profil privé ? Vos lectures ne seront plus visibles.";
  showConfirmation.value = true;
};

const cancelVisibilityChange = () => {
  isPublic.value = !isPublic.value;
  showConfirmation.value = false;
};

const confirmVisibilityChange = async () => {
  try {
    await authStore.changePublicVisibility();
    showConfirmation.value = false;
  } catch (error) {
    console.error("Error updating visibility:", error);
    isPublic.value = !isPublic.value;
  }
};

const copy = () => {
  navigator.clipboard.writeText(sharelink.value);
  showCopyAlert.value = true;
  setTimeout(() => {
    showCopyAlert.value = false;
  }, 3000);
};

onMounted(async () => {
  isPublic.value = await authStore.getPublicVisibility();
});
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
