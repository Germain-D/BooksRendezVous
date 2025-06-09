import { defineStore } from 'pinia';
import type { Book } from '@/types/book';


interface BooksState {
    books: Book[];
}



export const useBooksStore = defineStore('books', {
    state: (): BooksState => ({
        books: [],
    }),

    actions: {
        async addBook(newBook: Book) {
            
            try {
                const config = useRuntimeConfig();
                const booksData = await fetch(`${config.public.BACKEND_URL}/api/addbook`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    credentials: 'include',
                    body:
                        JSON.stringify({
                            book: newBook
                        })

                }).then(res => res.json());
                this.books.push(newBook);
                useStatsStore().fetchStats();
            }
            catch (error) {
                console.error('Error while adding book', error);
            }
        },

        removeBook(bookId: string) {

            try {
                const config = useRuntimeConfig();
                fetch(config.public.BACKEND_URL + '/api/books/' + bookId, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    credentials: 'include',
                })



                const index = this.books.findIndex((book) => book.id === bookId);
                if (index !== -1) {
                    this.books.splice(index, 1);
                }
                useStatsStore().fetchStats();

            } catch (error) {
                console.error('Error while removing book', error);
            }
        },

        updateBook(book: Book) {


            try {
                const config = useRuntimeConfig();
                book.rating = book.rating ? parseInt(String(book.rating)) : 0;
                fetch(config.public.BACKEND_URL + '/api/books/' + book.id, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    credentials: 'include',
                    body: JSON.stringify({
                        book: book
                    })
                })
                const index = this.books.findIndex((b) => b.id === book.id);
                if (index !== -1) {
                    this.books.splice(index, 1, book);
                }
                useStatsStore().fetchStats();
            }
            catch (error) {
                console.error('Error while updating book', error);
            }
        },

        getBook(bookId: string) {
            return this.books.find((book) => book.id === bookId);
        },

        async fetchBooks() {

            try {
                const config = useRuntimeConfig();
                const booksData = await fetch(config.public.BACKEND_URL + '/api/books', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    credentials: 'include'
                }).then(res => res.json());
                this.books = booksData.books;
   
            } catch (error) {
                console.error('Error while fetching books', error);
            }


        },
        async fetchPublicBooks(publicid : string) {

            try {
                const config = useRuntimeConfig();
                const booksData = await fetch(config.public.BACKEND_URL + '/api/publicuser', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    body: JSON.stringify({
                        publicid: publicid
                    })
                }).then(res => res.json());
                this.books = booksData.books;
                return booksData;
            } catch (error) {
                console.error('Error while fetching books', error);
            }


        },

        clearBooks() {
            this.books = [];
        }
    },
    getters: {
        getBooks(): Book[] {
            return this.books;
        }
    }
});