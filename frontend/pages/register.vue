<template>
  <div class="card max-w-md mt-10 mx-auto">
    <div class="card-body p-8 text-center">
      <div class="text-center space-y-2">
        <h1 class="card-title text-2xl font-semibold justify-center">
          Créer un compte
        </h1>
        <p>
          Inscrivez-vous pour accéder à votre espace personnel
        </p>
        <h1 class="card-title text-2xl font-semibold justify-center">
          Seul les utilisateurs autorisés peuvent s'inscrire
        </h1>
      </div>

      <div class="divider my-6"></div>

      <form @submit.prevent="handleRegister" class="space-y-4 flex flex-col justify-center mx-auto">
        <!-- Username Input -->
        <label class="input input-bordered flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="h-4 w-4 opacity-70">
            <path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
          </svg>
          <input 
            type="text" 
            class="grow" 
            placeholder="Nom d'utilisateur" 
            v-model="form.username"
            required
          />
        </label>

        <!-- Email Input -->
        <label class="input input-bordered flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="h-4 w-4 opacity-70">
            <path d="M2.5 3A1.5 1.5 0 0 0 1 4.5v.793c.026.009.051.02.076.032L7.674 8.51c.206.1.446.1.652 0l6.598-3.185A.755.755 0 0 1 15 5.293V4.5A1.5 1.5 0 0 0 13.5 3h-11Z" />
            <path d="M15 6.954 8.978 9.86a2.25 2.25 0 0 1-1.956 0L1 6.954V11.5A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5V6.954Z" />
          </svg>
          <input 
            type="email" 
            class="grow" 
            placeholder="Email" 
            v-model="form.email"
            required
          />
        </label>

        <!-- Password Input -->
        <label class="input input-bordered flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="h-4 w-4 opacity-70">
            <path fill-rule="evenodd" d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z" clip-rule="evenodd" />
          </svg>
          <input 
            type="password" 
            class="grow" 
            placeholder="Mot de passe" 
            v-model="form.password"
            required
          />
        </label>

        <button type="submit" class="btn btn-primary btn-lg w-full mt-6">
          S'inscrire
        </button>

        <div class="text-center text-sm text-neutral/70 mt-4">
          Déjà inscrit ? 
          <NuxtLink to="/login" class="text-primary hover:underline">
            Se connecter
          </NuxtLink>
        </div>

        <div v-if="error" class="alert alert-error">
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from "@/stores/auth"

definePageMeta({
  layout: 'landing',
})

const authStore = useAuthStore()
const error = ref('')

const form = ref({
  username: '',
  email: '',
  password: '',
})

const handleRegister = async () => {
  try {
    const success = await authStore.register({ 
      username: form.value.username,
      email: form.value.email, 
      password: form.value.password 
    })

    if (success) {
      navigateTo('/login')
    }
  } catch (err) {
    error.value = 'Une erreur est survenue lors de l\'inscription'
  }
}
</script>

<style scoped>
.invalid-feedback,
.empty-feedback {
  display: none;
}

.was-validated :placeholder-shown:invalid ~ .empty-feedback {
  display: block;
}

.was-validated :not(:placeholder-shown):invalid ~ .invalid-feedback {
  display: block;
}

.is-invalid,
.was-validated :invalid {
  border-color: #dc3545;
}
</style>