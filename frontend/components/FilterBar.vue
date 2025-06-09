<template>
  <div class="sticky top-0 z-10">
    <div
      class="bg-base-200/95 backdrop-blur-xl border-b border-base-300 shadow-lg py-4 px-4 md:px-6 mb-8"
    >
      <div
        class="flex flex-col md:flex-row gap-4 items-stretch md:items-center"
      >
        <!-- Filtres -->
        <div class="flex-1 flex items-center gap-2">
          <div class="md:hidden">
            <button
              @click="state.mobileMenuOpen = !state.mobileMenuOpen"
              class="btn btn-ghost btn-sm"
              :class="{ 'btn-active': state.mobileMenuOpen }"
            >
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
                  d="M4 6h16M4 12h16m-7 6h7"
                />
              </svg>
              <span class="ml-2">Filtres</span>
            </button>
          </div>

          <!-- Menu mobile -->
          <div
            v-if="state.mobileMenuOpen"
            class="absolute top-full left-0 right-0 mt-2 p-4 bg-base-200 shadow-xl md:hidden"
            @click.outside="state.mobileMenuOpen = false"
          >
            <div class="flex flex-col gap-2">
              <button
                v-for="filter in filters"
                :key="filter.id"
                class="btn btn-ghost justify-start gap-3 h-auto py-3"
                :class="{ 'btn-active': state.activeFilter === filter.id }"
                @click="
                  setActiveFilter(filter.id);
                  state.mobileMenuOpen = false;
                "
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.5"
                    :d="filter.icon"
                  />
                </svg>
                {{ filter.label }}
              </button>
            </div>
          </div>

          <!-- Menu desktop -->
          <div class="hidden md:flex gap-2">
            <button
              v-for="filter in filters"
              :key="filter.id"
              class="btn btn-ghost gap-2 min-h-0 h-auto py-2"
              :class="{ 'btn-active': state.activeFilter === filter.id }"
              @click="setActiveFilter(filter.id)"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  :d="filter.icon"
                />
              </svg>
              {{ filter.label }}
            </button>

            <!-- Genre Filter Desktop -->
            <div
              class="dropdown dropdown-bottom flex items-center gap-2"
              ref="dropdownDesktopRef"
            >
              <button tabindex="0" class="btn btn-ghost btn-sm">
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
                    d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"
                  />
                </svg>
                <span class="ml-2"
                  >Genres {{ selectedGenre ? `(${selectedGenre})` : "" }}</span
                >
              </button>
              <button
                v-if="selectedGenre"
                @click="setGenreFilter('')"
                class="btn btn-ghost btn-sm btn-circle"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
              <ul
                tabindex="0"
                class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52 z-[1]"
              >
                <li>
                  <a
                    @click="setGenreFilter('')"
                    :class="{ active: !selectedGenre }"
                  >
                    Tous les genres
                  </a>
                </li>
                <li v-for="genre in availableGenres" :key="genre">
                  <a
                    @click="setGenreFilter(genre)"
                    :class="{ active: selectedGenre === genre }"
                  >
                    {{ genre }}
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Recherche -->
        <div
          class="relative"
          @keydown.window.prevent.meta.k.exact="focusSearchInput"
          @keydown.window.prevent.ctrl.k.exact="focusSearchInput"
          @keydown.window.escape="clearSearch"
        >
          <div class="relative">
            <input
              ref="searchInput"
              type="text"
              placeholder="Rechercher un livre..."
              class="input input-bordered w-full md:w-80 pl-10 pr-16 h-auto py-2 bg-base-300/50"
              v-model="state.search"
              @input="filterBooks"
              @focus="state.searchFocused = true"
              @blur="state.searchFocused = false"
            />
            <div
              class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none"
            >
              <svg
                class="w-4 h-4 text-base-content/50"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
            <kbd
              class="absolute right-3 top-1/2 -translate-y-1/2 hidden md:inline-flex items-center gap-1 px-1.5 py-0.5 text-xs font-mono text-base-content/50 bg-base-300/50 rounded border border-base-content/10"
              :class="{
                'opacity-0': state.searchFocused,
                'opacity-100': !state.searchFocused,
              }"
            >
              <span class="text-[10px]">{{ shortcutKey }}</span>
              <span class="text-[10px]">K</span>
            </kbd>
          </div>

          <!-- Résultats de recherche -->
          <div
            v-if="state.search.length > 0"
            class="absolute mt-2 right-0 bg-base-200 rounded-lg shadow-xl p-4 w-full border border-base-300"
          >
            <div class="flex items-center justify-between text-sm">
              <span class="text-base-content/70">
                Résultats pour "<span>{{ state.search }}</span
                >"
              </span>
              <button @click="clearSearch" class="btn btn-ghost btn-xs">
                Effacer
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from "vue";
import { useBooksStore } from "@/stores/books";

const filters = [
  {
    id: "all",
    label: "Tous",
    icon: "M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z",
  },
  {
    id: "reading",
    label: "En cours",
    icon: "M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25",
  },
  {
    id: "to-read",
    label: "À lire",
    icon: "M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z",
  },
  {
    id: "finished",
    label: "Terminés",
    icon: "M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
  },
  {
    id: "favorites",
    label: "Coups de cœur",
    icon: "M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z",
  },
];

const booksStore = useBooksStore();
const selectedGenre = ref("");

const availableGenres = computed(() => {
  const genres = new Set<string>();
  booksStore.getBooks.forEach((book) => {
    book.genres?.forEach((genre) => {
      if (genre) genres.add(genre);
    });
  });
  return Array.from(genres).sort();
});

const shortcutKey = "⌘/Ctrl";
const searchInput = ref<HTMLInputElement | null>(null);
const state = reactive({
  activeFilter: "all",
  search: "",
  searchFocused: false,
  mobileMenuOpen: false,
});

const setActiveFilter = (filterId: string) => {
  state.activeFilter = filterId;
  filterBooks();
};

const dropdownDesktopRef = ref<HTMLElement | null>(null);
const dropdownMobileRef = ref<HTMLElement | null>(null);

const setGenreFilter = (genre: string) => {
  selectedGenre.value = genre;
  if (!genre) {
    state.activeFilter = "all";
  }
  // Ferme les deux dropdowns
  if (dropdownDesktopRef.value) {
    dropdownDesktopRef.value.removeAttribute("open");
  }
  if (dropdownMobileRef.value) {
    dropdownMobileRef.value.removeAttribute("open");
  }
  filterBooks();
};

const focusSearchInput = () => {
  searchInput.value?.focus();
};

const clearSearch = () => {
  state.search = "";
  searchInput.value?.blur();
};

const filterBooks = () => {
  const books = document.querySelectorAll("[data-book-status]");
  const searchTerm = state.search.toLowerCase();

  books.forEach((book) => {
    const status = book.getAttribute("data-book-status");
    const title = book.getAttribute("data-book-title")?.toLowerCase();
    const isFavorite = book.getAttribute("data-favorite") === "true";
    const genres = book.getAttribute("data-book-genres")?.split(",") || [];

    let matchesFilter = false;
    if (state.activeFilter === "favorites") {
      matchesFilter = isFavorite;
    } else {
      matchesFilter =
        state.activeFilter === "all" || status === state.activeFilter;
    }

    const matchesSearch = !searchTerm || title?.includes(searchTerm);
    const matchesGenre =
      !selectedGenre.value || genres.includes(selectedGenre.value);

    book.classList.toggle(
      "hidden",
      !(matchesFilter && matchesSearch && matchesGenre)
    );
  });
};

onMounted(() => {
  document.addEventListener("keydown", (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === "k") {
      e.preventDefault();
      focusSearchInput();
    }
  });
});
</script>

<style scoped lang="postcss">
.input:focus {
  @apply ring-2 ring-primary/20 ring-offset-0;
}

.input:focus + div svg {
  @apply text-primary;
}
</style>
