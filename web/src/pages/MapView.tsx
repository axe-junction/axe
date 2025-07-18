import { useEffect, useState } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  Polyline,
} from "react-leaflet";
import {
  Search,
  Filter,
  Navigation,
  MapPin,
  Clock,
  Target,
} from "lucide-react";
import TripCard from "../components/TripCard";
import { useAppStore } from "../store/useAppStore";
import type { Trip } from "../store/useAppStore";
import "leaflet/dist/leaflet.css";

// Mock data for demonstration
const mockTrips: Trip[] = [
  {
    id: "1",
    from: "Bab Ezzouar",
    to: "Alger Centre",
    duration: 35,
    cost: 150,
    provider: "ETUSA",
    transportModes: ["Bus", "Metro"],
    isFavorite: false,
    stops: 8,
    transferPoints: ["Station Tafourah", "Place des Martyrs"],
    route: [
      { lat: 36.7167, lng: 3.1833 },
      { lat: 36.7538, lng: 3.0588 },
    ],
  },
  {
    id: "2",
    from: "Bab Ezzouar",
    to: "Alger Centre",
    duration: 42,
    cost: 120,
    provider: "TRANSTU",
    transportModes: ["Tram"],
    isFavorite: true,
    stops: 12,
    transferPoints: ["Hussein Dey", "Place des Martyrs"],
    route: [
      { lat: 36.7167, lng: 3.1833 },
      { lat: 36.7538, lng: 3.0588 },
    ],
  },
];

export default function MapView() {
  const {
    mapCenter,
    selectedTrip,
    setSelectedTrip,
    trips,
    setTrips,
    searchQuery,
    setSearchQuery,
    recentSearches,
    addRecentSearch,
  } = useAppStore();

  const [showFilters, setShowFilters] = useState(false);
  const [searchResults, setSearchResults] = useState<Trip[]>([]);
  const [currentPosition, setCurrentPosition] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const [searchFocused, setSearchFocused] = useState(false);
  const [destinationQuery, setDestinationQuery] = useState("");

  useEffect(() => {
    // Initialize with mock data
    setTrips(mockTrips);
    setSearchResults(mockTrips);
  }, [setTrips]);

  // Get current position
  useEffect(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setCurrentPosition({
            lat: position.coords.latitude,
            lng: position.coords.longitude,
          });
        },
        (error) => {
          console.error("Error getting location:", error);
        }
      );
    }
  }, []);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    if (query.trim()) {
      const filtered = trips.filter(
        (trip) =>
          trip.from.toLowerCase().includes(query.toLowerCase()) ||
          trip.to.toLowerCase().includes(query.toLowerCase())
      );
      setSearchResults(filtered);
    } else {
      setSearchResults(trips);
    }
  };

  const handleDestinationSearch = (query: string) => {
    setDestinationQuery(query);
    if (query.trim()) {
      const filtered = trips.filter(
        (trip) =>
          trip.to.toLowerCase().includes(query.toLowerCase()) ||
          trip.from.toLowerCase().includes(query.toLowerCase())
      );
      setSearchResults(filtered);
    } else {
      setSearchResults([]);
    }
  };

  const handleRecentDestinationSelect = (destination: string) => {
    setDestinationQuery(destination);
    handleDestinationSearch(destination);
    setSearchFocused(false);
  };

  const handleDestinationSelect = (destination: string) => {
    addRecentSearch(destination);
    setDestinationQuery(destination);
    setSearchFocused(false);
  };

  const handleTripSelect = (trip: Trip) => {
    setSelectedTrip(trip);
  };

  return (
    <div className="h-screen flex flex-col overflow-hidden">
      {/* Header with Search */}
      <div className="bg-white shadow-sm border-b p-4 flex-shrink-0">
        <div className="space-y-3">
          {/* Current Location */}
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <Target size={16} className="text-green-500" />
            <span>
              {currentPosition
                ? "Current location detected"
                : "Detecting location..."}
            </span>
          </div>

          {/* Destination Search */}
          <div className="relative">
            <Search
              size={20}
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
            />
            <input
              type="text"
              placeholder="Where do you want to go?"
              value={destinationQuery}
              onChange={(e) => handleDestinationSearch(e.target.value)}
              onFocus={() => setSearchFocused(true)}
              onBlur={() => setTimeout(() => setSearchFocused(false), 200)}
              className="w-full pl-10 pr-16 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-base"
            />
            <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex gap-1">
              <button
                onClick={() => setShowFilters(!showFilters)}
                className="p-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <Filter size={16} className="text-gray-600" />
              </button>
              <button className="p-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors">
                <Navigation size={16} />
              </button>
            </div>
          </div>

          {/* Recent Destinations Dropdown */}
          {(searchFocused || destinationQuery) && (
            <div className="absolute top-full left-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg z-10 max-h-64 overflow-y-auto">
              {destinationQuery && searchResults.length > 0 && (
                <div className="p-2">
                  <div className="text-xs text-gray-500 mb-2 px-2">
                    Search Results
                  </div>
                  {searchResults.slice(0, 3).map((trip) => (
                    <button
                      key={trip.id}
                      onClick={() => handleDestinationSelect(trip.to)}
                      className="w-full text-left p-3 hover:bg-gray-50 rounded-lg transition-colors"
                    >
                      <div className="flex items-center gap-3">
                        <MapPin size={16} className="text-gray-400" />
                        <div>
                          <div className="font-medium text-gray-900">
                            {trip.to}
                          </div>
                          <div className="text-sm text-gray-500">
                            from {trip.from}
                          </div>
                        </div>
                      </div>
                    </button>
                  ))}
                </div>
              )}

              {!destinationQuery && recentSearches.length > 0 && (
                <div className="p-2">
                  <div className="text-xs text-gray-500 mb-2 px-2">
                    Recent Destinations
                  </div>
                  {recentSearches.slice(0, 5).map((search, index) => (
                    <button
                      key={index}
                      onClick={() =>
                        handleRecentDestinationSelect(
                          search.split(" → ")[1] || search
                        )
                      }
                      className="w-full text-left p-3 hover:bg-gray-50 rounded-lg transition-colors"
                    >
                      <div className="flex items-center gap-3">
                        <Clock size={16} className="text-gray-400" />
                        <span className="text-gray-700">{search}</span>
                      </div>
                    </button>
                  ))}
                </div>
              )}

              {!destinationQuery && recentSearches.length === 0 && (
                <div className="p-4 text-center text-gray-500">
                  <MapPin size={24} className="mx-auto mb-2 text-gray-400" />
                  <div className="text-sm">No recent destinations</div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="bg-white border-b p-4 flex-shrink-0">
          <div className="flex gap-2 flex-wrap">
            <button className="px-3 py-1 bg-gray-100 rounded-full text-sm">
              Bus
            </button>
            <button className="px-3 py-1 bg-gray-100 rounded-full text-sm">
              Metro
            </button>
            <button className="px-3 py-1 bg-gray-100 rounded-full text-sm">
              Tram
            </button>
            <button className="px-3 py-1 bg-gray-100 rounded-full text-sm">
              Walk
            </button>
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Map */}
        <div className="flex-1 relative">
          <MapContainer center={mapCenter} zoom={12} className="h-full w-full">
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />

            {/* Current Position Marker */}
            {currentPosition && (
              <Marker position={[currentPosition.lat, currentPosition.lng]}>
                <Popup>Your current location</Popup>
              </Marker>
            )}

            {selectedTrip && (
              <>
                <Marker
                  position={[
                    selectedTrip.route[0].lat,
                    selectedTrip.route[0].lng,
                  ]}
                >
                  <Popup>{selectedTrip.from}</Popup>
                </Marker>
                <Marker
                  position={[
                    selectedTrip.route[selectedTrip.route.length - 1].lat,
                    selectedTrip.route[selectedTrip.route.length - 1].lng,
                  ]}
                >
                  <Popup>{selectedTrip.to}</Popup>
                </Marker>
                <Polyline
                  positions={selectedTrip.route.map((point) => [
                    point.lat,
                    point.lng,
                  ])}
                  color="#6316DB"
                  weight={4}
                />
              </>
            )}
          </MapContainer>
        </div>

        {/* Results Panel - Desktop */}
        <div className="w-80 bg-white border-l border-gray-200 flex flex-col lg:flex hidden">
          <div className="p-4 border-b flex-shrink-0">
            <h2 className="text-lg font-semibold text-gray-900">
              {destinationQuery ? `Routes to ${destinationQuery}` : "Routes"} (
              {searchResults.length})
            </h2>
          </div>
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {searchResults.length > 0 ? (
              searchResults.map((trip) => (
                <TripCard
                  key={trip.id}
                  trip={trip}
                  onSelect={handleTripSelect}
                />
              ))
            ) : (
              <div className="text-center py-8">
                <MapPin size={48} className="text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">
                  {destinationQuery
                    ? "No routes found for this destination"
                    : "Search for a destination to see routes"}
                </p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mobile Results - Bottom Sheet */}
      <div
        className="lg:hidden bg-white border-t flex-shrink-0"
        style={{ maxHeight: "40vh" }}
      >
        <div className="p-4 border-b">
          <h2 className="text-lg font-semibold text-gray-900">
            {destinationQuery ? `Routes to ${destinationQuery}` : "Routes"} (
            {searchResults.length})
          </h2>
        </div>
        <div className="overflow-y-auto p-4 space-y-4">
          {searchResults.length > 0 ? (
            searchResults
              .slice(0, 3)
              .map((trip) => (
                <TripCard
                  key={trip.id}
                  trip={trip}
                  onSelect={handleTripSelect}
                />
              ))
          ) : (
            <div className="text-center py-6">
              <MapPin size={32} className="text-gray-400 mx-auto mb-2" />
              <p className="text-gray-500 text-sm">
                {destinationQuery
                  ? "No routes found"
                  : "Search for a destination"}
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
