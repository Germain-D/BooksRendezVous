<template>
    <div 
      class="card bg-base-200 shadow-xl h-full hover:shadow-2xl transition-all duration-300 group"
      :data-favorite="book.favorite ? 'true' : 'false'"
    >
    <NuxtLink :to="isConnected ? `/admin/${book.id}` : `/book/${book.id}`" class="h-full flex flex-col">

        <figure class="relative w-full overflow-hidden">
          <div v-if="book.imageUrl" class="w-full h-[280px] sm:h-[320px] lg:h-[360px]">
            <img
              :src="book.imageUrl"
              :alt="`Couverture du livre ${book.title}`"
              class="w-full h-full object-cover bg-base-300 transition-transform duration-300 group-hover:scale-105"
              loading="lazy"
            />
          </div>
          <div v-else class="w-full h-[280px] sm:h-[320px] lg:h-[360px] bg-base-300 flex items-center justify-center">
            <span class="text-base-content/50">Image non disponible</span>
          </div>
          <div class="absolute top-2 left-2 flex gap-1 flex-wrap max-w-[calc(100%-1rem)]">
            <span class="badge badge-sm md:badge-md badge-primary shadow-lg">{{ getStatusLabel(book.status) }}</span>
            <span v-if="book.abandoned" class="badge badge-sm md:badge-md badge-error shadow-lg">Abandonné</span>
            <span v-if="book.favorite" class="badge badge-sm md:badge-md badge-secondary shadow-lg">
              <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              Coup de cœur
            </span>
          </div>
        </figure>
        
        <div class="card-body p-3 sm:p-4">
          <h2 class="card-title text-sm sm:text-base line-clamp-2 group-hover:text-primary transition-colors">
            {{ book.title }}
          </h2>
          
          <p v-if="book.authors" class="text-xs sm:text-sm opacity-70 line-clamp-1">
            par {{ book.authors.join(', ') }}
          </p>
  
          <div v-if="book.status === 'reading' && book.progress" class="mt-2">
            <progress 
              class="progress progress-primary w-full" 
              :value="book.progress" 
              max="100"
            ></progress>
            <p class="text-xs opacity-70 mt-1">{{ book.progress }}% lu</p>
          </div>
  
          <div v-if="book.rating" class="rating rating-xs sm:rating-sm mt-2">
            <template v-for="i in 5">
              <input
                type="radio"
                class="mask mask-star-2 bg-primary"
                disabled
                :checked="i === book.rating"
              />
            </template>
          </div>
        </div>
      </NuxtLink>
    </div>
  </template>
  
  <script setup lang="ts">
  import { slugify } from '../utils/slugify';

  import { useAuthStore } from '@/stores/auth';

  const authStore = useAuthStore();
  const isConnected = authStore.isLoggedIn;
  
  interface Book {
    id: string;
    googleBooksId: string;
    status: string;
    progress?: number;
    startDate?: string;
    endDate?: string;
    rating?: number | null;
    comment?: string | null;
    abandoned?: boolean;
    genres?: string[];
    favorite?: boolean;
    title: string;
    authors: string[];
    imageUrl?: string;
  }
  
  interface Props {
    book: Book;
    compact?: boolean;
  }
  
  const props = defineProps<Props>();
  const { book, compact = false } = props;
  
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('fr-FR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };
  
  const getStatusBadgeClass = (status: string) => {
    const baseClasses = "text-sm px-3 py-1 rounded-full font-medium";
    switch (status) {
      case 'reading':
        return `${baseClasses} bg-primary/20 text-primary border border-primary/30`;
      case 'finished':
        return `${baseClasses} bg-green-500/20 text-green-400 border border-green-500/30`;
      case 'to-read':
        return `${baseClasses} bg-dark-lighter text-gray-300 border border-gray-600`;
      default:
        return baseClasses;
    }
  };
  
  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'reading': return 'En cours';
      case 'finished': return 'Terminé';
      case 'to-read': return 'À lire';
      default: return status;
    }
  };
  </script>
  
  <style scoped>
  .card {
    transition: all 0.3s ease;
  }
  
  .card:hover {
    transform: translateY(-5px);
  }
  
  .line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  </style>