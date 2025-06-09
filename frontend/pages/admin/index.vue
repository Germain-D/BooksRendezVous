<template>
  <div class="max-w-7xl m-16">
    <!-- Back button -->
    <NuxtLink
      to="/"
      class="inline-flex items-center gap-2 text-gray-400 hover:text-primary mb-8 transition-colors"
    >
      <svg
        class="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 19l-7-7m0 0l7-7m-7 7h18"
        />
      </svg>
      Retour à la liste
    </NuxtLink>

    <h1 class="text-3xl font-bold mb-8">Gestion des livres</h1>

    <!-- Add Book Form -->
    <div class="card bg-base-200 p-6 mb-8">
      <h2 class="text-xl font-bold mb-4">Ajouter un nouveau livre</h2>

      <div class="flex gap-4 mb-4 flex-col md:flex-row">
        <input
          type="text"
          v-model="searchQuery"
          @keyup.enter="searchBooks"
          placeholder="Rechercher un livre..."
          class="input input-bordered flex-1"
        />
        <button @click="searchBooks" class="btn btn-primary">Rechercher</button>
      </div>

      <div v-if="searchResults.length > 0" class="mb-4">
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div
            v-for="book in searchResults"
            :key="book.id"
            class="card bg-base-300 hover:bg-base-100 cursor-pointer transition-colors"
            @click="
              selectedBook = book;
              searchResults = [];
            "
          >
            <figure class="px-4 pt-4">
              <img
                :src="book.volumeInfo.imageLinks?.thumbnail"
                :alt="book.volumeInfo.title"
                class="h-40 object-cover rounded"
              />
            </figure>
            <div class="card-body p-4">
              <h3 class="card-title text-sm">{{ book.volumeInfo.title }}</h3>
              <p class="text-xs opacity-70">
                {{ book.volumeInfo.authors?.join(", ") }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <form
        v-if="selectedBook"
        @submit.prevent="handleSubmit(selectedBook)"
        class="space-y-4"
      >
        <!-- Book preview -->
        <div class="flex gap-4">
          <img
            :src="selectedBook.volumeInfo.imageLinks?.thumbnail"
            :alt="selectedBook.volumeInfo.title"
            class="h-40 object-cover rounded"
          />
          <div class="space-y-4">
            <h3 class="text-lg font-bold">
              {{ selectedBook.volumeInfo.title }}
            </h3>
            <p class="text-sm opacity-70">
              {{ selectedBook.volumeInfo.authors?.join(", ") }}
            </p>
            <p class="text-sm">{{ selectedBook.volumeInfo.description }}</p>
          </div>
        </div>

        <!-- Status select -->
        <div>
          <label for="status" class="label">Statut</label>
          <select
            v-model="status"
            id="status"
            class="select select-bordered w-full"
          >
            <option value="to-read">À lire</option>
            <option value="reading">En cours</option>
            <option value="finished">Terminé</option>
          </select>
        </div>

        <!-- Rating -->
        <div>
          <label class="label">Note</label>
          <div class="rating gap-1">
            <input
              v-for="i in 5"
              :key="i"
              type="radio"
              :name="`rating-${i}`"
              class="mask mask-star-2 bg-orange-400"
              :value="i"
              v-model="rating"
            />
          </div>
        </div>

        <!-- favorite -->
        <div>
          <label for="favorite" class="label">Coup de coeur</label>
          <input
            type="checkbox"
            v-model="favorite"
            id="favorite"
            class="toggle toggle-primary"
          />
        </div>

        <!-- Comment -->
        <div>
          <label for="comment" class="label">Commentaire</label>
          <textarea
            v-model="comment"
            id="comment"
            class="textarea textarea-bordered h-24 w-full"
          ></textarea>
        </div>

        <button type="submit" class="btn btn-primary w-full">
          Ajouter à ma bibliothèque
        </button>
      </form>
    </div>

    <!-- Books Grid -->
    <div class="space-y-6">
      <h2 class="text-xl font-bold">Mes livres</h2>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div
          v-for="book in books"
          :key="book.id"
          class="card bg-base-300 hover:bg-base-100 transition-colors"
        >
          <NuxtLink :to="`/admin/${book.id}`" class="h-full flex flex-col">
            <figure class="px-4 pt-4">
              <img
                :src="book.imageUrl"
                :alt="book.title"
                class="h-40 object-cover rounded"
              />
            </figure>
            <div class="card-body p-4">
              <h3 class="card-title text-sm">{{ book.title }}</h3>
              <p class="text-xs opacity-70">{{ book.authors?.join(", ") }}</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<!-- pages/admin/books.vue -->
<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { Book, GoogleBook, BookStatus } from "@/types/book";

import { useBooksStore } from "@/stores/books";

// Store initialization

const booksStore = useBooksStore();

// Form state
const searchQuery = ref("");
const searchResults = ref<GoogleBook[]>([]);
const selectedBook = ref<GoogleBook | null>(null);
const status = ref<BookStatus>("to-read");
const rating = ref<number>(3);
const comment = ref("");
const favorite = ref(false);

const isLoading = ref(false);

// Search books from Google API
const searchBooks = async () => {
  if (!searchQuery.value.trim() || isLoading.value) return;

  isLoading.value = true;
  try {
    const response = await fetch(
      `https://www.googleapis.com/books/v1/volumes?q=${encodeURIComponent(
        searchQuery.value
      )}&fields=items(id,volumeInfo(title,authors,description,imageLinks(*),pageCount,publishedDate,categories))`
    );
    const data = await response.json();
    searchResults.value = data.items || [];
  } catch (error) {
    console.error("Error searching books:", error);
    searchResults.value = [];
  } finally {
    isLoading.value = false;
  }
};

// Initialize books and visibility
const initializeBooks = async () => {
  let booksData = booksStore.getBooks;
  if (booksData.length === 0) {
    await booksStore.fetchBooks();
    booksData = booksStore.getBooks;
  }
  books.value = booksData;
};

// Books state
const books = ref<Book[]>([]);

onMounted(async () => {
  await initializeBooks();
});

// Handle book submission
const handleSubmit = async (book: GoogleBook) => {
  const newBook: Book = {
    id: crypto.randomUUID(),
    googleBooksId: book.id,
    status: status.value,
    title: book.volumeInfo.title,
    authors: book.volumeInfo.authors || ["Auteur inconnu"],
    description: book.volumeInfo.description,
    imageUrl: book.volumeInfo.imageLinks?.thumbnail,
    rating: parseInt(String(rating.value)),
    comment: comment.value,
    genres: book.volumeInfo.categories || [],
    pageCount: parseInt(String(book.volumeInfo.pageCount)) || 0,
    publishedDate: book.volumeInfo.publishedDate || "",
    favorite: favorite.value,
  };

  try {
    await booksStore.addBook(newBook);
    books.value = booksStore.getBooks;
    // Reset form
    selectedBook.value = null;
    status.value = "to-read";
    rating.value = 3;
    comment.value = "";
    searchQuery.value = "";
    searchResults.value = [];
  } catch (error) {
    console.error("Error adding book:", error);
  }
};

// Page configuration
definePageMeta({
  layout: "connected",
  middleware: ["guard"],
});
</script>

<style scoped>
.rating {
  @apply flex gap-2 my-4;
}

.rating input {
  @apply cursor-pointer;
}
</style>
