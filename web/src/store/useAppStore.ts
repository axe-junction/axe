import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

export interface Trip {
  id: string;
  from: string;
  to: string;
  duration: number;
  cost: number;
  provider: string;
  transportModes: string[];
  isFavorite: boolean;
  stops: number;
  transferPoints: string[];
  route: {
    lat: number;
    lng: number;
  }[];
}

export interface User {
  id: string;
  name: string;
  email: string;
  avatar: string;
}

export interface CommunityPost {
  id: string;
  user: User;
  content: string;
  image?: string;
  likes: number;
  comments: number;
  timestamp: Date;
  isLiked: boolean;
}

interface AppState {
  // User state
  user: User | null;

  // Trip state
  trips: Trip[];
  favoriteTrips: Trip[];
  recentSearches: string[];
  recentDestinations: string[];

  // Community state
  communityPosts: CommunityPost[];

  // Map state
  selectedTrip: Trip | null;
  mapCenter: [number, number];
  selectedDestination: { name: string; lat: number; lng: number } | null;

  // Search state
  searchQuery: string;
  searchFilters: {
    transportModes: string[];
    maxDuration: number;
    maxCost: number;
    provider: string;
  };
}

interface AppActions {
  // User actions
  setUser: (user: User | null) => void;
  updateUser: (updates: Partial<User>) => void;

  // Trip actions
  setTrips: (trips: Trip[]) => void;
  addTrip: (trip: Trip) => void;
  toggleFavorite: (tripId: string) => void;
  setSelectedTrip: (trip: Trip | null) => void;

  // Search actions
  setSearchQuery: (query: string) => void;
  addRecentSearch: (query: string) => void;
  setSearchFilters: (filters: Partial<AppState["searchFilters"]>) => void;

  // Community actions
  setCommunityPosts: (posts: CommunityPost[]) => void;
  toggleLike: (postId: string) => void;

  // Map actions
  setMapCenter: (center: [number, number]) => void;

  // Destination actions
  addRecentDestination: (destination: string) => void;
  setSelectedDestination: (
    destination: { name: string; lat: number; lng: number } | null
  ) => void;
}

export const useAppStore = create<AppState & AppActions>()(
  persist(
    (set) => ({
      // Initial state
      user: null,
      trips: [
        // Sample VTC trips with corrected realistic pricing
        {
          id: "vtc-1",
          from: "Current Location",
          to: "Houari Boumediene Airport",
          duration: 35,
          cost: 1200,
          provider: "Yassir",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
        // El Harrach trips (higher prices - 15 min ETA)
        {
          id: "vtc-2",
          from: "Current Location",
          to: "El Harrach",
          duration: 15,
          cost: 341,
          provider: "Yassir",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
        {
          id: "vtc-4",
          from: "Current Location",
          to: "El Harrach",
          duration: 15,
          cost: 352,
          provider: "Heetch",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
        {
          id: "vtc-6",
          from: "Current Location",
          to: "El Harrach",
          duration: 15,
          cost: 341,
          provider: "InDrive",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
        // Alger Centre trips (lower prices - 8 min ETA)
        {
          id: "vtc-3",
          from: "Current Location",
          to: "Alger Centre",
          duration: 8,
          cost: 871,
          provider: "Yassir",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
        {
          id: "vtc-5",
          from: "Current Location",
          to: "Alger Centre",
          duration: 8,
          cost: 842,
          provider: "Heetch",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
        {
          id: "vtc-7",
          from: "Current Location",
          to: "Alger Centre",
          duration: 8,
          cost: 870,
          provider: "InDrive",
          transportModes: ["VTC"],
          isFavorite: false,
          stops: 0,
          transferPoints: [],
          route: [],
        },
      ],
      favoriteTrips: [],
      recentSearches: [],
      recentDestinations: [],
      communityPosts: [],
      selectedTrip: null,
      mapCenter: [36.7538, 3.0588], // Algiers coordinates
      selectedDestination: null,
      searchQuery: "",
      searchFilters: {
        transportModes: [],
        maxDuration: 120,
        maxCost: 1500, // Increased to accommodate VTC prices
        provider: "all",
      },

      // Actions
      setUser: (user) => set({ user }),
      updateUser: (updates) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...updates } : null,
        })),

      setTrips: (trips) => set({ trips }),
      addTrip: (trip) =>
        set((state) => ({
          trips: [...state.trips, trip],
        })),
      toggleFavorite: (tripId) =>
        set((state) => ({
          trips: state.trips.map((trip) =>
            trip.id === tripId
              ? { ...trip, isFavorite: !trip.isFavorite }
              : trip
          ),
          favoriteTrips: state.favoriteTrips.some((t) => t.id === tripId)
            ? state.favoriteTrips.filter((t) => t.id !== tripId)
            : [
                ...state.favoriteTrips,
                state.trips.find((t) => t.id === tripId)!,
              ],
        })),
      setSelectedTrip: (trip) => set({ selectedTrip: trip }),

      setSearchQuery: (query) => set({ searchQuery: query }),
      addRecentSearch: (query) =>
        set((state) => ({
          recentSearches: [
            query,
            ...state.recentSearches.filter((q) => q !== query),
          ].slice(0, 10),
        })),
      setSearchFilters: (filters) =>
        set((state) => ({
          searchFilters: { ...state.searchFilters, ...filters },
        })),

      setCommunityPosts: (posts) => set({ communityPosts: posts }),
      toggleLike: (postId) =>
        set((state) => ({
          communityPosts: state.communityPosts.map((post) =>
            post.id === postId
              ? {
                  ...post,
                  isLiked: !post.isLiked,
                  likes: post.isLiked ? post.likes - 1 : post.likes + 1,
                }
              : post
          ),
        })),

      setMapCenter: (center) => set({ mapCenter: center }),
      addRecentDestination: (destination) =>
        set((state) => ({
          recentDestinations: [
            destination,
            ...state.recentDestinations.filter((d) => d !== destination),
          ].slice(0, 10),
        })),
      setSelectedDestination: (destination) =>
        set({ selectedDestination: destination }),
    }),
    {
      name: "axe-web-storage",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        favoriteTrips: state.favoriteTrips,
        recentSearches: state.recentSearches,
        recentDestinations: state.recentDestinations,
        searchFilters: state.searchFilters,
      }),
    }
  )
);
