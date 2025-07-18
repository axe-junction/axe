import { useEffect, useState } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  Polyline,
} from "react-leaflet";
import { Search, Filter, Navigation } from "lucide-react";
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
  } = useAppStore();

  const [showFilters, setShowFilters] = useState(false);
  const [searchResults, setSearchResults] = useState<Trip[]>([]);

  useEffect(() => {
    // Initialize with mock data
    setTrips(mockTrips);
    setSearchResults(mockTrips);
  }, [setTrips]);

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

  const handleTripSelect = (trip: Trip) => {
    setSelectedTrip(trip);
  };

  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <div className="bg-white shadow-sm border-b p-4">
        <div className="flex items-center gap-3">
          <div className="flex-1 relative">
            <Search
              size={20}
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
            />
            <input
              type="text"
              placeholder="Search routes..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
            />
          </div>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="p-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <Filter size={20} className="text-gray-600" />
          </button>
          <button className="p-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors">
            <Navigation size={20} />
          </button>
        </div>
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="bg-white border-b p-4">
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

      <div className="flex-1 flex">
        {/* Map */}
        <div className="flex-1 relative">
          <MapContainer center={mapCenter} zoom={12} className="h-full w-full">
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />

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

        {/* Results Panel */}
        <div className="w-80 bg-white border-l border-gray-200 flex flex-col lg:block hidden">
          <div className="p-4 border-b">
            <h2 className="text-lg font-semibold text-gray-900">
              Routes ({searchResults.length})
            </h2>
          </div>
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {searchResults.map((trip) => (
              <TripCard key={trip.id} trip={trip} onSelect={handleTripSelect} />
            ))}
          </div>
        </div>
      </div>

      {/* Mobile Results */}
      <div className="lg:hidden bg-white border-t max-h-64 overflow-y-auto">
        <div className="p-4">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Routes ({searchResults.length})
          </h2>
          <div className="space-y-4">
            {searchResults.slice(0, 3).map((trip) => (
              <TripCard key={trip.id} trip={trip} onSelect={handleTripSelect} />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
