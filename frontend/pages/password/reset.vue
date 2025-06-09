<template>
    <div class="card max-w-md mt-10 mx-auto">
      <div class="card-body p-8 text-center">
        <div class="text-center space-y-2">
          <h1 class="card-title text-2xl font-semibold justify-center">
            Créer un nouveau mot de passe
          </h1>
          <p>
            Veuillez entrer votre nouveau mot de passe
          </p>
        </div>
  
        <div class="divider my-6"></div>
  
        <form @submit.prevent="submitForm" class="space-y-4 flex flex-col justify-center mx-auto">
          <div>
            <label class="input input-bordered flex items-center gap-2" :class="{ 'input-error': errors.password }">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                class="h-4 w-4 opacity-70">
                <path
                  fill-rule="evenodd"
                  d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                  clip-rule="evenodd" />
              </svg>
              <input 
                :type="showPassword ? 'text' : 'password'"
                id="password" 
                v-model="password" 
                required 
                placeholder="Nouveau mot de passe"
                class="grow"
                @input="validateInRealTime"
              />
              <button
                type="button"
                class="btn btn-ghost btn-sm p-0"
                @click="showPassword = !showPassword"
              >
                <Icon :name="showPassword ? 'ph:eye-slash' : 'ph:eye'" />
              </button>
            </label>
            
            <div class="mt-2 space-y-1 text-sm text-left">
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
  
          <label class="input input-bordered flex items-center gap-2" :class="{ 'input-error': errors.confirmPassword }">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="h-4 w-4 opacity-70">
              <path
                fill-rule="evenodd"
                d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                clip-rule="evenodd" />
            </svg>
            <input 
              type="password" 
              id="confirmPassword" 
              v-model="confirmPassword" 
              required 
              placeholder="Confirmer le mot de passe"
              class="grow"
            />
          </label>
          <p v-if="errors.confirmPassword" class="text-sm text-error text-left">{{ errors.confirmPassword }}</p>
  
          <button 
            type="submit" 
            class="btn btn-primary btn-lg w-full mt-6"
            :disabled="isLoading || !isFormValid"
          >
            <span v-if="isLoading">Traitement en cours...</span>
            <span v-else>Réinitialiser le mot de passe</span>
          </button>
        </form>
        
        <div v-if="successMessage" class="alert alert-success mt-6">
          {{ successMessage }}
        </div>
        
        <div class="mt-6 text-center">
          <NuxtLink to="/login" class="text-primary">
            Retourner à la connexion
          </NuxtLink>
        </div>
      </div>
    </div>
  
    <!-- Error Modal -->
    <div
      v-if="showErrorModal"
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
        <span>{{ errorMessage }}</span>
        <div>
          <button
            class="btn btn-sm"
            @click="closeErrorModal"
          >
            Fermer
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>

definePageMeta({
  layout: 'landing',
});


  import { ref, computed, onMounted } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useAuthStore } from '@/stores/auth';
  
  const route = useRoute();
  const router = useRouter();
  const authStore = useAuthStore();
  
  // Form data
  const password = ref('');
  const confirmPassword = ref('');
  const token = ref('');
  const errors = ref({});
  const isLoading = ref(false);
  const successMessage = ref('');
  
  // Password validation
  const showPassword = ref(false);
  const lengthValid = ref(false);
  const numberValid = ref(false);
  const specialCharValid = ref(false);
  const showErrorModal = ref(false);
  const errorMessage = ref('');
  
  const isFormValid = computed(() => {
    return lengthValid.value && 
           numberValid.value && 
           specialCharValid.value && 
           password.value === confirmPassword.value &&
           password.value.length > 0;
  });
  
  // Get token from URL when component mounts
  onMounted(async () => {
  token.value = route.query.token;
  if (!token.value) {
    // Redirect to forget password page if no token is present
    router.push('/password/forget');
  } else {
    // verify token
    try {
      const isValid = await authStore.verifyResetToken(token.value);
      if (isValid) {
        console.log('Token is valid');
      } else {
        console.log('Token is invalid');
        errorMessage.value = "Le lien de réinitialisation est invalide ou a expiré.";
        showErrorModal.value = true;
        setTimeout(() => {
          router.push('/password/forget');
        }, 3000);
      }
    } catch (error) {
      console.error('Error verifying token:', error);
      errorMessage.value = "Erreur lors de la vérification du lien de réinitialisation.";
      showErrorModal.value = true;
      setTimeout(() => {
        router.push('/password/forget');
      }, 3000);
    }
  }
});
  
  const validateInRealTime = () => {
    // Length validation - at least 10 characters
    lengthValid.value = password.value.length >= 10;
  
    // Number validation - at least one digit
    numberValid.value = /\d/.test(password.value);
  
    // Special character validation
    specialCharValid.value = /[!@#$%^&*(),.?":{}|<>]/.test(password.value);
  };
  
  const closeErrorModal = () => {
    showErrorModal.value = false;
    errorMessage.value = '';
  };
  
  const validateForm = () => {
    errors.value = {};
    let isValid = true;
    
    // Password validation is handled by validateInRealTime
    if (!lengthValid.value || !numberValid.value || !specialCharValid.value) {
      errorMessage.value = "Veuillez respecter tous les critères du mot de passe.";
      showErrorModal.value = true;
      isValid = false;
    }
    
    // Confirmation password validation
    if (password.value !== confirmPassword.value) {
      errors.value.confirmPassword = "Les mots de passe ne correspondent pas";
      isValid = false;
    }
    
    return isValid;
  };
  
  const submitForm = async () => {
    if (!validateForm()) return;
    
    isLoading.value = true;
    try {
      // Here you would call your API to reset the password
      await authStore.resetPassword(token.value, password.value).then(() =>
        console.log('Password reset successfully')
      );
      
      successMessage.value = "Votre mot de passe a été réinitialisé avec succès. Vous pouvez maintenant vous connecter avec votre nouveau mot de passe.";
      password.value = '';
      confirmPassword.value = '';
      
      // Redirect to login page after a short delay
      setTimeout(() => {
        router.push('/login');
      }, 3000);
    } catch (error) {
      console.error('Error resetting password:', error);
      errorMessage.value = "Une erreur s'est produite lors de la réinitialisation du mot de passe. Veuillez réessayer.";
      showErrorModal.value = true;
    } finally {
      isLoading.value = false;
    }
  };
  </script>