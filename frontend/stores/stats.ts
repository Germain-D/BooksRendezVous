import { defineStore } from 'pinia';

interface StatsState {
    totalBooks: number;
    completedBooks: number;
    toReadBooks: number;
    readingBooks: number;
    favoriteBooks: number;
    totalPages: number;
    averageRating: number;
}

export const useStatsStore = defineStore('stats', {
    state: (): StatsState => ({
        totalBooks: 0,
        completedBooks: 0,
        toReadBooks: 0,
        readingBooks: 0,
        favoriteBooks: 0,
        totalPages: 0,
        averageRating: 0
    }),

    actions: {
        async fetchStats() {
            try {
                const config = useRuntimeConfig();
                const response = await fetch(`${config.public.BACKEND_URL}/api/stats`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    credentials: 'include'
                });
                
                const statsData = await response.json();
                console.log(statsData);
                this.totalBooks = statsData.totalBooks;
                this.completedBooks = statsData.completedBooks;
                this.toReadBooks = statsData.toReadBooks;
                this.readingBooks = statsData.readingBooks;
                this.favoriteBooks = statsData.favoriteBooks;
                this.totalPages = statsData.totalPages;
                this.averageRating = statsData.averageRating;
            } catch (error) {
                console.error('Error while fetching stats', error);
            }
        },

        clearStats() {
            this.totalBooks = 0;
            this.completedBooks = 0;
            this.toReadBooks = 0;
            this.readingBooks = 0;
            this.favoriteBooks = 0;
            this.totalPages = 0;
            this.averageRating = 0;
        }
    },

    getters: {
        getTotalBooks: (state) => state.totalBooks,
        getCompletedBooks: (state) => state.completedBooks,
        getToReadBooks: (state) => state.toReadBooks,
        getReadingBooks: (state) => state.readingBooks,
        getFavoriteBooks: (state) => state.favoriteBooks,
        getTotalPages: (state) => state.totalPages,
        getAverageRating: (state) => state.averageRating,
    }
});