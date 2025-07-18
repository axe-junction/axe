import { useState, useEffect } from "react";
import { Search, MapPin, Clock, Filter, X } from "lucide-react";
import TripCard from "../components/TripCard";
import { useAppStore } from "../store/useAppStore";
import type { Trip } from "../store/useAppStore";

export default function SearchPage() {
  const {
    setSearchQuery,
    recentSearches,
    addRecentSearch,
    searchFilters,
    setSearchFilters,
    trips,
  } = useAppStore();

  const [fromLocation, setFromLocation] = useState("");
  const [toLocation, setToLocation] = useState("");
  const [searchResults, setSearchResults] = useState<Trip[]>([]);
  const [showFilters, setShowFilters] = useState(false);
  const [activeFilters, setActiveFilters] = useState(searchFilters);

  useEffect(() => {
    if (fromLocation && toLocation) {
      const filtered = trips.filter((trip) => {
        const matchesRoute =
          trip.from.toLowerCase().includes(fromLocation.toLowerCase()) &&
          trip.to.toLowerCase().includes(toLocation.toLowerCase());
        const matchesTransport =
          activeFilters.transportModes.length === 0 ||
          trip.transportModes.some((mode) =>
            activeFilters.transportModes.includes(mode)
          );
        const matchesDuration = trip.duration <= activeFilters.maxDuration;
        const matchesCost = trip.cost <= activeFilters.maxCost;
        const matchesProvider =
          activeFilters.provider === "all" ||
          trip.provider === activeFilters.provider;

        return (
          matchesRoute &&
          matchesTransport &&
          matchesDuration &&
          matchesCost &&
          matchesProvider
        );
      });
      setSearchResults(filtered);
    } else {
      setSearchResults([]);
    }
  }, [fromLocation, toLocation, activeFilters, trips]);

  const handleSearch = () => {
    if (fromLocation && toLocation) {
      const searchQuery = `${fromLocation} → ${toLocation}`;
      addRecentSearch(searchQuery);
      setSearchQuery(searchQuery);
    }
  };

  const handleFilterToggle = (filterType: string, value: string) => {
    if (filterType === "transportModes") {
      const newModes = activeFilters.transportModes.includes(value)
        ? activeFilters.transportModes.filter((mode) => mode !== value)
        : [...activeFilters.transportModes, value];
      setActiveFilters((prev) => ({ ...prev, transportModes: newModes }));
    }
  };

  const applyFilters = () => {
    setSearchFilters(activeFilters);
    setShowFilters(false);
  };

  const clearFilters = () => {
    const defaultFilters = {
      transportModes: [],
      maxDuration: 120,
      maxCost: 500,
      provider: "all",
    };
    setActiveFilters(defaultFilters);
    setSearchFilters(defaultFilters);
  };

  const transportModes = ["Bus", "Metro", "Tram", "Walk"];
  const providers = ["all", "ETUSA", "TRANSTU", "SETRAM"];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="p-4">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Search Routes
          </h1>

          {/* Search Form */}
          <div className="space-y-3">
            <div className="relative">
              <MapPin
                size={20}
                className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              />
              <input
                type="text"
                placeholder="From"
                value={fromLocation}
                onChange={(e) => setFromLocation(e.target.value)}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              />
            </div>

            <div className="relative">
              <MapPin
                size={20}
                className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              />
              <input
                type="text"
                placeholder="To"
                value={toLocation}
                onChange={(e) => setToLocation(e.target.value)}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              />
            </div>

            <div className="flex gap-2">
              <button
                onClick={handleSearch}
                className="flex-1 bg-primary text-white py-3 px-4 rounded-lg hover:bg-primary/90 transition-colors flex items-center justify-center gap-2"
              >
                <Search size={20} />
                Search Routes
              </button>

              <button
                onClick={() => setShowFilters(!showFilters)}
                className="px-4 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <Filter size={20} className="text-gray-600" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Filters Modal */}
      {showFilters && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end lg:items-center justify-center">
          <div className="bg-white rounded-t-lg lg:rounded-lg w-full lg:w-96 max-h-[80vh] overflow-y-auto">
            <div className="p-4 border-b flex items-center justify-between">
              <h2 className="text-lg font-semibold">Filters</h2>
              <button
                onClick={() => setShowFilters(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X size={20} />
              </button>
            </div>

            <div className="p-4 space-y-6">
              {/* Transport Modes */}
              <div>
                <h3 className="font-medium mb-3">Transport Modes</h3>
                <div className="flex flex-wrap gap-2">
                  {transportModes.map((mode) => (
                    <button
                      key={mode}
                      onClick={() => handleFilterToggle("transportModes", mode)}
                      className={`px-3 py-2 rounded-full text-sm transition-colors ${
                        activeFilters.transportModes.includes(mode)
                          ? "bg-primary text-white"
                          : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                      }`}
                    >
                      {mode}
                    </button>
                  ))}
                </div>
              </div>

              {/* Max Duration */}
              <div>
                <h3 className="font-medium mb-3">Max Duration</h3>
                <div className="flex items-center gap-2">
                  <Clock size={16} className="text-gray-400" />
                  <input
                    type="range"
                    min="15"
                    max="180"
                    step="15"
                    value={activeFilters.maxDuration}
                    onChange={(e) =>
                      setActiveFilters((prev) => ({
                        ...prev,
                        maxDuration: parseInt(e.target.value),
                      }))
                    }
                    className="flex-1"
                  />
                  <span className="text-sm text-gray-600">
                    {activeFilters.maxDuration}min
                  </span>
                </div>
              </div>

              {/* Max Cost */}
              <div>
                <h3 className="font-medium mb-3">Max Cost</h3>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-600">DA</span>
                  <input
                    type="range"
                    min="50"
                    max="1000"
                    step="50"
                    value={activeFilters.maxCost}
                    onChange={(e) =>
                      setActiveFilters((prev) => ({
                        ...prev,
                        maxCost: parseInt(e.target.value),
                      }))
                    }
                    className="flex-1"
                  />
                  <span className="text-sm text-gray-600">
                    {activeFilters.maxCost} DA
                  </span>
                </div>
              </div>

              {/* Provider */}
              <div>
                <h3 className="font-medium mb-3">Provider</h3>
                <select
                  value={activeFilters.provider}
                  onChange={(e) =>
                    setActiveFilters((prev) => ({
                      ...prev,
                      provider: e.target.value,
                    }))
                  }
                  className="w-full p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                >
                  {providers.map((provider) => (
                    <option key={provider} value={provider}>
                      {provider === "all" ? "All Providers" : provider}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="p-4 border-t flex gap-2">
              <button
                onClick={clearFilters}
                className="flex-1 py-2 px-4 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Clear All
              </button>
              <button
                onClick={applyFilters}
                className="flex-1 py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors"
              >
                Apply Filters
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Recent Searches */}
      {recentSearches.length > 0 && searchResults.length === 0 && (
        <div className="bg-white m-4 rounded-lg shadow-sm border">
          <div className="p-4 border-b">
            <h2 className="font-semibold text-gray-900">Recent Searches</h2>
          </div>
          <div className="p-4 space-y-2">
            {recentSearches.slice(0, 5).map((search, index) => (
              <button
                key={index}
                onClick={() => {
                  const [from, to] = search.split(" → ");
                  setFromLocation(from);
                  setToLocation(to);
                }}
                className="w-full text-left p-3 hover:bg-gray-50 rounded-lg transition-colors"
              >
                <div className="flex items-center gap-2">
                  <Clock size={16} className="text-gray-400" />
                  <span className="text-gray-700">{search}</span>
                </div>
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Search Results */}
      {searchResults.length > 0 && (
        <div className="p-4">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Found {searchResults.length} routes
          </h2>
          <div className="space-y-4">
            {searchResults.map((trip) => (
              <TripCard key={trip.id} trip={trip} />
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {fromLocation && toLocation && searchResults.length === 0 && (
        <div className="flex flex-col items-center justify-center py-12">
          <Search size={48} className="text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            No routes found
          </h3>
          <p className="text-gray-600 text-center">
            Try adjusting your search criteria or filters
          </p>
        </div>
      )}
    </div>
  );
}
