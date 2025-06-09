<template>
  <details class="collapse bg-base-200 mb-4">
    <summary class="collapse-title text-xl font-medium p-6">
      Changer le mot de passe
    </summary>
    <div class="collapse-content px-6">
      <form @submit.prevent="handlePasswordChange" class="space-y-4">
        <div>
          <label class="label">Mot de passe actuel</label>
          <input
            type="password"
            v-model="currentPassword"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Nouveau mot de passe</label>
          <div class="relative">
            <input
              v-model="newPassword"
              :type="showNewPassword ? 'text' : 'password'"
              class="input input-bordered w-full pr-10"
              @input="validateInRealTime"
            />
            <button
              type="button"
              class="absolute right-2 top-1/2 -translate-y-1/2"
              @click="showNewPassword = !showNewPassword"
            >
              <Icon :name="showNewPassword ? 'ph:eye-slash' : 'ph:eye'" />
            </button>
            <div class="mt-2 space-y-1 text-sm">
              <div class="flex items-center gap-2">
                <Icon
                  :name="lengthValid ? 'ph:check-circle' : 'ph:x-circle'"
                  :class="lengthValid ? 'text-success' : 'text-error'"
                />
                <span>10 caractères minimum</span>
              </div>
              <div class="flex items-center gap-2">
                <Icon
                  :name="numberValid ? 'ph:check-circle' : 'ph:x-circle'"
                  :class="numberValid ? 'text-success' : 'text-error'"
                />
                <span>Au moins 1 chiffre</span>
              </div>
              <div class="flex items-center gap-2">
                <Icon
                  :name="specialCharValid ? 'ph:check-circle' : 'ph:x-circle'"
                  :class="specialCharValid ? 'text-success' : 'text-error'"
                />
                <span>Au moins 1 caractère spécial</span>
              </div>
            </div>
          </div>
        </div>
        <div>
          <label class="label">Confirmer le mot de passe</label>
          <input
            type="password"
            v-model="confirmPassword"
            class="input input-bordered w-full"
          />
        </div>
        <button type="submit" class="btn btn-primary w-full">
          Mettre à jour le mot de passe
        </button>
      </form>
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
        <button class="btn btn-sm" @click="cancelPasswordChange">
          Annuler
        </button>
        <button
          class="btn btn-sm btn-primary ml-2"
          @click="confirmPasswordChange"
        >
          Confirmer
        </button>
      </div>
    </div>
  </div>

  <!-- error modal -->
  <div
    v-if="ShowerrorPassword"
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-4"
  >
    <div role="alert" class="alert max-w-lg bg-base-200">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="stroke-error h-6 w-6 shrink-0"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M6 18L18 6M6 6l12 12"
        ></path>
      </svg>
      <span>{{ errorPassword }}</span>
      <div>
        <button
          class="btn btn-sm"
          @click="
            errorPassword = '';
            ShowerrorPassword = false;
          "
        >
          Fermer
        </button>
      </div>
    </div>
  </div>

  <!-- success modal -->
  <div
    v-if="showSuccessModal"
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-4"
  >
    <div role="alert" class="alert max-w-lg bg-base-200">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="stroke-success h-6 w-6 shrink-0"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
      <span>{{ successMessage }}</span>
      <div>
        <button class="btn btn-sm" @click="showSuccessModal = false">
          Fermer
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

// Password change
const currentPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const errorPassword = ref("");
const ShowerrorPassword = ref(false);

// Confirmation modal
const showConfirmation = ref(false);
const confirmationMessage = ref("");
const pendingVisibility = ref(false);

const lengthValid = ref(false);
const numberValid = ref(false);
const specialCharValid = ref(false);

const showNewPassword = ref(false);

const showSuccessModal = ref(false);
const successMessage = ref("");

const validateInRealTime = () => {
  // Length validation
  lengthValid.value = newPassword.value.length >= 10;

  // Number validation
  numberValid.value = /\d/.test(newPassword.value);

  // Special character validation
  specialCharValid.value = /[!@#$%^&*(),.?":{}|<>]/.test(newPassword.value);
};

const validatePassword = (password: string): boolean => {
  validateInRealTime();

  if (!lengthValid.value || !numberValid.value || !specialCharValid.value) {
    ShowerrorPassword.value = true;
    errorPassword.value =
      "Veuillez respecter tous les critères du mot de passe.";
    return false;
  }

  return true;
};

const handlePasswordChange = async () => {
  //pendingVisibility.value = isPublic.value;
  if (!validatePassword(newPassword.value)) {
    return;
  }

  if (newPassword.value !== confirmPassword.value) {
    errorPassword.value = "Les mots de passe ne correspondent pas.";
    ShowerrorPassword.value = true;
    return;
  }

  confirmationMessage.value =
    "Êtes-vous sûr de vouloir changer votre mot de passe ?";
  showConfirmation.value = true;
};

const cancelPasswordChange = () => {
  //isPublic.value = !isPublic.value;
  showConfirmation.value = false;
};

const confirmPasswordChange = async () => {
  try {
    const success = await authStore.passwordChange(
      currentPassword.value,
      newPassword.value
    );
    if (success) {
      showConfirmation.value = false;
      // Reset form
      currentPassword.value = "";
      newPassword.value = "";
      confirmPassword.value = "";
      // Show success message
      successMessage.value = "Mot de passe modifié avec succès";
      showSuccessModal.value = true;
    } else {
      errorPassword.value = "Erreur lors du changement de mot de passe";
      ShowerrorPassword.value = true;
      showConfirmation.value = false;
    }
  } catch (error: any) {
    console.error("Error updating password:", error);
    errorPassword.value =
      error.message || "Erreur lors du changement de mot de passe";
    ShowerrorPassword.value = true;
    showConfirmation.value = false;
  }
};
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
