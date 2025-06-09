
<template>
    <div class="card  max-w-md  mt-10 mx-auto ">
      <div class="card-body p-8 text-center">
        <div class="text-center space-y-2">
          <h1 class="card-title text-2xl font-semibold justify-center">
            Se Connecter
          </h1>
          <p >
            Connectez-vous pour accéder à votre espace personnel
          </p>
        </div>
  
        <div class="divider my-6"></div>
  
        <form @submit.prevent="handleSimpleLogin" class="space-y-4 flex flex-col justify-center mx-auto">
  
          <label class="input input-bordered flex items-center gap-2">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                class="h-4 w-4 opacity-70">
                <path
                d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
            </svg>
            <input type="email" class="grow" placeholder="Email"  v-model="form.email"/>
            </label>
          <label class="input input-bordered flex items-center gap-2">
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
            <input type="password"             placeholder="Mot de passe"
            required class="grow" v-model="form.password" />
            </label>
  
          <button type="submit" class="btn btn-primary btn-lg w-full mt-6">
            Se Connecter
          </button>
  
          <div id="result" class="text-center text-sm text-neutral/70 mt-4"></div>
        </form>
        <NuxtLink to="/register" class="text-primary">Créer un compte</NuxtLink>
        <NuxtLink to="/password/forget" class="text-primary">Mot de passe oublié ?</NuxtLink>
      </div>
    </div>
  </template>
  

<script setup>

definePageMeta({
  layout: 'landing',
});


import { useAuthStore } from "@/stores/auth";
import { useBooksStore } from "@/stores/books";


const form = ref({
  email: "",
  password: "",
});


const authStore = useAuthStore();
const booksStore = useBooksStore();

const handleSimpleLogin = async () => {

const success = await authStore.login({ 
  email: form.value.email, 
  password: form.value.password 
});

if (success) {
  booksStore.fetchBooks();
  navigateTo('/admin');
} else {
  alert('Login failed');
}
};
  


</script>

<style>
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