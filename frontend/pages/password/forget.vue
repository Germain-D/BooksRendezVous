<template>
    <div class="card max-w-md mt-10 mx-auto">
      <div class="card-body p-8 text-center">
        <div class="text-center space-y-2">
          <h1 class="card-title text-2xl font-semibold justify-center">
            Réinitialisation du mot de passe
          </h1>
          <p>
            Entrez votre adresse e-mail pour recevoir un lien de réinitialisation
          </p>
        </div>
  
        <div class="divider my-6"></div>
  
        <form @submit.prevent="submitForm" class="space-y-4 flex flex-col justify-center mx-auto">
          <label class="input input-bordered flex items-center gap-2" :class="{ 'input-error': errors.email }">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="h-4 w-4 opacity-70">
              <path d="M2.5 3A1.5 1.5 0 0 0 1 4.5v.793c.026.009.051.02.076.032L7.674 8.51c.206.1.446.1.652 0l6.598-3.185A.755.755 0 0 1 15 5.293V4.5A1.5 1.5 0 0 0 13.5 3h-11Z" />
              <path d="M15 6.954 8.978 9.86a2.25 2.25 0 0 1-1.956 0L1 6.954V11.5A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5V6.954Z" />
            </svg>
            <input 
              type="email" 
              id="email" 
              v-model="email" 
              required 
              placeholder="Votre adresse e-mail"
              class="grow"
            />
          </label>
          <p v-if="errors.email" class="text-sm text-error text-left">{{ errors.email }}</p>
  
          <button 
            type="submit" 
            class="btn btn-primary btn-lg w-full mt-6"
            :disabled="isLoading"
          >
            <span v-if="isLoading">Envoi en cours...</span>
            <span v-else>Envoyer le lien</span>
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
  </template>
  
  <script>
  definePageMeta({
  layout: 'landing',
});

import { useAuthStore } from "@/stores/auth";
  export default {
    name: 'ForgetPassword',
    layout: 'auth',
    data() {
      return {
        email: '',
        errors: {},
        isLoading: false,
        successMessage: ''
      }
    },
    methods: {
      validateForm() {
        this.errors = {};
        let isValid = true;
        
        // Email validation
        if (!this.email) {
          this.errors.email = "L'adresse e-mail est requise";
          isValid = false;
        } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(this.email)) {
          this.errors.email = "Veuillez entrer une adresse e-mail valide";
          isValid = false;
        }
        
        return isValid;
      },
      async submitForm() {
        if (!this.validateForm()) return;
        
        this.isLoading = true;
        try {
          // Here you would call your API to send the reset password email
          // For example:
          // await this.$axios.post('/api/auth/reset-password', { email: this.email });
          
          const authStore = useAuthStore();
            await authStore.sendResetPasswordEmail(this.email);
          // Simulated response delay
          await new Promise(resolve => setTimeout(resolve, 1000));
          
          this.successMessage = `Un email avec les instructions de réinitialisation a été envoyé à ${this.email}`;
          this.email = '';
        } catch (error) {
          console.error('Error sending reset password email:', error);
          this.errors.email = "Une erreur s'est produite. Veuillez réessayer.";
        } finally {
          this.isLoading = false;
        }
      }
    }
  }
  </script>