<template>
  <HeaderBar :pseudo="pseudo" />
  <div
    v-if="FetchStatus === 'error'"
    class="flex items-center justify-center my-8"
  >
    <div class="text-center">
      <h1 class="text-3xl font-bold mb-4">Erreur de chargement</h1>
      <p class="text-lg mb -4">
        Une erreur est survenue lors du chargement des livres. Veuillez
        réessayer.
      </p>
      <p class="text-lg mb-4">- Le profil n'est pas public ou n'existe pas</p>
      <p class="text-lg mb-4">- Aucun livre n'a été ajouté</p>
      <button @click="fetchBooks" class="btn btn-primary mt-4">
        Réessayer
      </button>
    </div>
  </div>
  <FilterBar />

  <div
    class="grid gap-4 sm:gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 max-w-[1800px] mx-auto px-4"
  >
    <div
      v-for="book in shuffledBooks"
      :key="book.id"
      :data-book-status="book.status"
      :data-book-title="book.title"
      :data-favorite="book.favorite ? 'true' : 'false'"
      class="animate-fade-in"
    >
      <PublicBookCard
        :book="book"
        :compact="true"
        :publicid="publicid"
        :pseudo="rawPseudo"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import type { Book } from "@/types/book";
import { useBooksStore } from "@/stores/books";

// Page meta
definePageMeta({
  layout: "static",
});

// Route and navigation
const route = useRoute();
const pseudo = computed(() => {
  const rawPseudo = route.params.pseudo as string;
  return rawPseudo.charAt(0).toUpperCase() + rawPseudo.slice(1);
});
const rawPseudo = route.params.pseudo as string;
const publicid = route.params.publicid as string;

const FetchStatus = ref("");

// Books state
const books = ref<Book[]>([]); // Initialize with empty array
const booksStore = useBooksStore();

// Fetch books
const fetchBooks = async () => {
  try {
    const response = await booksStore.fetchPublicBooks(publicid);
    books.value = response?.books || [];
    if (books.value.length) FetchStatus.value = "success";
    else FetchStatus.value = "error";
  } catch (error) {
    console.error("Error fetching books:", error);
    FetchStatus.value = "error";
    books.value = []; // Ensure books is always an array
  }
};

// Call fetch books immediately
await fetchBooks();

// Computed properties
const shuffledBooks = computed(() => {
  if (!books.value?.length) return [];
  return [...books.value].sort(() => Math.random() - 0.5);
});

// Utility functions
const getBooksByStatus = (status: string): Book[] => {
  return shuffledBooks.value.filter((book) => book.status === status);
};

// Page metadata
useHead({
  title: `${pseudo.value} - Mes lectures`,
  meta: [
    {
      name: "description",
      content: `Explorez la bibliothèque personnelle de ${pseudo.value} : découvrez les livres en cours de lecture, terminés et à venir.`,
    },
  ],
});
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.5s ease-out forwards;
  opacity: 0;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in:nth-child(n) {
  animation-delay: calc(n * 0.05s);
}
</style>
