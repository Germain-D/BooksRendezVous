<template>
  <div class="navbar bg-base-100">
    <div class="navbar-start">
      <div class="dropdown">
        <div tabindex="0" role="button" class="btn btn-ghost lg:hidden">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h8m-8 6h16"
            />
          </svg>
        </div>
        <ul
          tabindex="0"
          class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
        >
          <li><NuxtLink to="/">Mes Lectures</NuxtLink></li>
          <li><NuxtLink to="/admin">Gérer ses livres</NuxtLink></li>
          <li><NuxtLink to="/profile">Mon Profil</NuxtLink></li>
        </ul>
      </div>
      <NuxtLink to="/">
        <a class="btn btn-ghost text-xl">
          <span
            class="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent"
            >{{ pseudo || "X" }}</span
          >
          <span>BRDV</span>
        </a></NuxtLink
      >
    </div>
    <div class="navbar-center hidden lg:flex">
      <ul class="menu menu-horizontal px-1">
        <li><NuxtLink to="/">Mes Lectures</NuxtLink></li>
        <li><NuxtLink to="/admin">Gérer ses livres</NuxtLink></li>
        <li><NuxtLink to="/profile">Mon Profil</NuxtLink></li>
      </ul>
    </div>
    <div class="navbar-end">
      <button
        class="btn border-none ml-10 btn-error btn-outline"
        @click="handleLogout"
      >
        Se déconnecter
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useBooksStore } from "@/stores/books";
const authStore = useAuthStore();
const booksStore = useBooksStore();

const pseudo = computed(() => {
  try {
    return authStore.getPseudo || "X";
  } catch (error) {
    return "X";
  }
});

console.log(authStore.getPseudo);

const open = ref(false);

const handleLogout = async () => {
  await authStore.logout();
  booksStore.clearBooks();

  navigateTo("/login");
};
</script>
