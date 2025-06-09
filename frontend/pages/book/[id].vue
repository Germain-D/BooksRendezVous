<template>

      <div class="max-w-4xl mt-16 mx-auto">
        <NuxtLink to="/" class="inline-flex items-center gap-2 text-gray-400 hover:text-primary mb-8 transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Retour à la liste
        </NuxtLink>
  
        <div class="bg-dark-light rounded-xl p-6 md:p-8">
          <div class="grid md:grid-cols-[300px,1fr] gap-8">
            <div class="aspect-[2/3] rounded-lg overflow-hidden bg-dark-lighter">
              <img
                v-if="book.imageUrl"
                :src="book.imageUrl"
                :alt="`Couverture de ${book.title}`"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-gray-500">
                Image non disponible
              </div>
            </div>
  
            <div class="space-y-6">
              <div>
                <h1 class="text-3xl font-bold text-gray-100 mb-2">
                  {{ book.title }}
                  <span v-if="book.favorite" class="inline-flex items-center text-primary ml-2" title="Coup de cœur">
                    <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  </span>
                </h1>
                <p v-if="book.authors" class="text-xl text-gray-400">
                  par {{ book.authors.join(', ') }}
                </p>
              </div>
  
              <div class="flex flex-wrap gap-2">
                <span v-if="book.status" :class="statusClass(book.status)">
                  {{ statusText(book.status) }}
                </span>
                <span v-if="book.abandoned" class="px-3 py-1 rounded-full text-sm bg-error/20 text-error">
                  Abandonné
                </span>
              </div>
  
              <div v-if="book.status === 'reading' && book.progress" class="bg-dark rounded-lg p-4">
                <div class="w-full bg-dark-lighter rounded-full h-3 mb-2">
                  <div
                    class="bg-gradient-to-r from-primary to-accent h-3 rounded-full transition-all duration-300"
                    :style="{ width: `${book.progress}%` }"
                  />
                </div>
                <p class="text-gray-400">{{ book.progress }}% lu</p>
              </div>
  
              <div v-if="book.description" class="prose prose-invert prose-gray max-w-none">
                <h2 class="text-xl font-semibold mb-3">Résumé</h2>
                <p class="text-gray-300" v-html="book.description" />
              </div>
  
              <div v-if="book.rating || book.comment" class="bg-dark rounded-lg p-6 border-l-4 border-primary">
                <h2 class="text-xl font-semibold mb-3">Mon avis</h2>
                <p v-if="book.comment" class="text-gray-300 italic mb-4">"{{ book.comment }}"</p>
                <div v-if="book.rating" class="flex items-center gap-1">
                  <svg
                    v-for="i in 5"
                    :key="i"
                    :class="`w-6 h-6 ${i <= book.rating ? 'text-primary' : 'text-dark-lighter'}`"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                </div>
              </div>
  
              <div class="flex flex-wrap gap-4 text-sm text-gray-400">
                <p v-if="book.startDate">Commencé le {{ formatDate(book.startDate) }}</p>
                <p v-if="book.endDate">Terminé le {{ formatDate(book.endDate) }}</p>
              </div>
            </div>
          </div>
        </div>
  
        <div v-if="book.googleBooksId" class="mt-8">
          <a 
            :href="googleBooksPreviewUrl"
            target="_blank"
            rel="noopener noreferrer" 
            class="btn btn-primary w-full md:w-auto"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 488 512">
              <path fill="currentColor" d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 123 24.5 166.3 64.9l-67.5 64.9C258.5 52.6 94.3 116.6 94.3 256c0 86.5 69.1 156.6 153.7 156.6 98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 24.9 3.9 41.4z"/>
            </svg>
            Voir sur Google Books
          </a>
        </div>
      </div>

  </template>
  
  <script setup lang="ts">
  definePageMeta({
  layout: 'landing',
});


  import { useRoute, useRouter } from 'vue-router';

  import booksData from '@/assets/data/books.json';
  import type { Book } from '../../types/book';
  
  const route = useRoute();
  const router = useRouter();
  const bookId = route.params.id as string;
  const book = booksData.books.find((b) => b.id === bookId) as unknown as Book;
  
  if (!book) {
    router.push('/404');
  }
  
  const googleBooksPreviewUrl = `https://books.google.fr/books?id=${book.googleBooksId}&printsec=frontcover`;
  
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('fr-FR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };
  
  const statusClass = (status: string) => {
    return `px-3 py-1 rounded-full text-sm ${
      status === 'reading' ? 'bg-primary/20 text-primary' :
      status === 'finished' ? 'bg-success/20 text-success' :
      'bg-warning/20 text-warning'
    }`;
  };
  
  const statusText = (status: string) => {
    return status === 'reading' ? 'En lecture' :
           status === 'finished' ? 'Terminé' :
           'À lire';
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
  
    /* Décalage de l'animation pour chaque carte */
    .animate-fade-in:nth-child(n) {
      animation-delay: calc(n * 0.05s);
    }
  </style>