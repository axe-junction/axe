import { useState } from "react";
import { Heart, Clock, Trash2, Search } from "lucide-react";
import TripCard from "../components/TripCard";
import { useAppStore } from "../store/useAppStore";

export default function FavoritesPage() {
  const { favoriteTrips, recentSearches, toggleFavorite } = useAppStore();
  const [activeTab, setActiveTab] = useState<"favorites" | "recent">(
    "favorites"
  );
  const [searchQuery, setSearchQuery] = useState("");

  const filteredFavorites = favoriteTrips.filter(
    (trip) =>
      trip.from.toLowerCase().includes(searchQuery.toLowerCase()) ||
      trip.to.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const filteredRecentSearches = recentSearches.filter((search) =>
    search.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleRemoveFavorite = (tripId: string) => {
    toggleFavorite(tripId);
  };

  return (
    <div className="min-h-full bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="p-4">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Favorites & Recent
          </h1>

          {/* Search Bar */}
          <div className="relative mb-4">
            <Search
              size={20}
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
            />
            <input
              type="text"
              placeholder="Search favorites and recent..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
            />
          </div>

          {/* Tabs */}
          <div className="flex border-b border-gray-200">
            <button
              onClick={() => setActiveTab("favorites")}
              className={`flex-1 py-2 px-4 text-center font-medium border-b-2 transition-colors ${
                activeTab === "favorites"
                  ? "border-primary text-primary"
                  : "border-transparent text-gray-500 hover:text-gray-700"
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <Heart size={16} />
                Favorites ({filteredFavorites.length})
              </div>
            </button>
            <button
              onClick={() => setActiveTab("recent")}
              className={`flex-1 py-2 px-4 text-center font-medium border-b-2 transition-colors ${
                activeTab === "recent"
                  ? "border-primary text-primary"
                  : "border-transparent text-gray-500 hover:text-gray-700"
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <Clock size={16} />
                Recent ({filteredRecentSearches.length})
              </div>
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="p-4">
        {activeTab === "favorites" && (
          <div className="space-y-4">
            {filteredFavorites.length > 0 ? (
              filteredFavorites.map((trip) => (
                <div key={trip.id} className="relative">
                  <TripCard trip={trip} />
                  <button
                    onClick={() => handleRemoveFavorite(trip.id)}
                    className="absolute top-4 right-4 p-2 bg-red-50 text-red-500 rounded-full hover:bg-red-100 transition-colors"
                    title="Remove from favorites"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              ))
            ) : (
              <div className="text-center py-12">
                <Heart size={48} className="text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  {searchQuery
                    ? "No matching favorites"
                    : "No favorite trips yet"}
                </h3>
                <p className="text-gray-600">
                  {searchQuery
                    ? "Try adjusting your search terms"
                    : "Add trips to favorites by tapping the heart icon"}
                </p>
              </div>
            )}
          </div>
        )}

        {activeTab === "recent" && (
          <div className="space-y-2">
            {filteredRecentSearches.length > 0 ? (
              filteredRecentSearches.map((search, index) => (
                <div
                  key={index}
                  className="bg-white rounded-lg shadow-sm border p-4 hover:shadow-md transition-shadow cursor-pointer"
                >
                  <div className="flex items-center gap-3">
                    <Clock size={16} className="text-gray-400" />
                    <span className="text-gray-900">{search}</span>
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-12">
                <Clock size={48} className="text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  {searchQuery
                    ? "No matching recent searches"
                    : "No recent searches"}
                </h3>
                <p className="text-gray-600">
                  {searchQuery
                    ? "Try adjusting your search terms"
                    : "Your recent searches will appear here"}
                </p>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="fixed bottom-20 right-4 lg:bottom-4">
        <div className="bg-white rounded-lg shadow-lg border p-4 space-y-2">
          <div className="text-sm text-gray-600">Quick Stats</div>
          <div className="flex items-center gap-2">
            <Heart size={16} className="text-red-500" />
            <span className="text-sm font-medium">
              {favoriteTrips.length} favorites
            </span>
          </div>
          <div className="flex items-center gap-2">
            <Clock size={16} className="text-blue-500" />
            <span className="text-sm font-medium">
              {recentSearches.length} recent searches
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
