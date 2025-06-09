export type BookStatus = 'reading' | 'finished' | 'to-read';

export interface Book {
  id: string;
  googleBooksId: string;
  status: BookStatus; // Changed from String to BookStatus type
  genres?: string[];
  progress?: number;
  startDate?: string;
  endDate?: string;
  rating?: number | null;
  comment?: string | null;
  abandoned?: boolean;
  favorite?: boolean;
  title: string;
  authors: string[];
  imageUrl?: string;
  description?: string;
  pageCount?: number;
  publishedDate?: string;
}


export interface GoogleBook {
  id: string
  volumeInfo: {
    title: string
    authors?: string[]
    description?: string
    imageLinks?: {
      thumbnail?: string
      [key: string]: string | undefined
    }
    pageCount?: number
    categories?: string[]
    publishedDate?: string
  }
  selectedImage?: string
}