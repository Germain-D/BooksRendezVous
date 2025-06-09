<template>
  <HeaderBar :pseudo="pseudo || 'X'" />
  <FilterBar />
  <div
    class="grid gap-4 sm:gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 max-w-[1800px] mx-auto px-4"
  >
    <div
      v-for="book in shuffledBooks"
      :key="book.id"
      :data-book-status="book.status"
      :data-book-title="book.title"
      :data-book-genres="book.genres"
      :data-favorite="book.favorite ? 'true' : 'false'"
      class="animate-fade-in"
    >
      <BookCard :book="book" :compact="true" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Book } from "@/types/book";
import Xbooks from "@/public/data/books.json";
import { ref, computed, onMounted } from "vue";

useHead({
  title: "BooksRDV",
  meta: [
    {
      name: "description",
      content:
        "Créez votre bibliothèque personnelle et partagez vos lectures préférées.",
    },
  ],
});

// Reactive variables
const isConnected = ref(false);
const pseudo = ref("X");
const books = ref<Book[]>([]);
const shuffledBooks = computed(() =>
  [...books.value].sort(() => Math.random() - 0.5)
);

// Stores (initialized in onMounted)
let authStore: ReturnType<typeof useAuthStore>;
let booksStore: ReturnType<typeof useBooksStore>;

onMounted(async () => {
  // Initialize stores client-side
  authStore = useAuthStore();
  booksStore = useBooksStore();

  // Check authentication state
  isConnected.value = authStore.isLoggedIn;

  setPageLayout(isConnected.value ? "connected" : "landing");

  // Get user pseudo
  try {
    pseudo.value = authStore.getPseudo || "X";
  } catch (error) {
    pseudo.value = "X";
  }

  // Initialize books data
  await initializeBooks();
});

async function initializeBooks() {
  if (isConnected.value) {
    if (booksStore.getBooks.length === 0) {
      await booksStore.fetchBooks();
    }
    books.value = booksStore.getBooks;
  } else {
    books.value = Xbooks.books as Book[];
  }
}

// Utility function (if needed elsewhere)
const getBooksByStatus = (status: string): Book[] => {
  return shuffledBooks.value.filter((book) => book.status === status);
};
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
