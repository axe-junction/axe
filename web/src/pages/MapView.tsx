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
  Route,
} from "lucide-react";
import { useAppStore } from "../store/useAppStore";
import "leaflet/dist/leaflet.css";

// Real destinations in Algiers
const popularDestinations = [
  { name: "Alger Centre", lat: 36.7538, lng: 3.0588, category: "District" },
  { name: "Bab Ezzouar", lat: 36.7167, lng: 3.1833, category: "District" },
  { name: "Hydra", lat: 36.7669, lng: 3.0347, category: "District" },
  { name: "Kouba", lat: 36.7333, lng: 3.0833, category: "District" },
  { name: "Birtouta", lat: 36.6167, lng: 3.0167, category: "District" },
  {
    name: "University of Algiers",
    lat: 36.7167,
    lng: 3.1833,
    category: "University",
  },
  {
    name: "Houari Boumediene Airport",
    lat: 36.691,
    lng: 3.2154,
    category: "Airport",
  },
  { name: "Algiers Port", lat: 36.7697, lng: 3.0611, category: "Port" },
  { name: "Maqam Echahid", lat: 36.7525, lng: 3.0519, category: "Monument" },
  { name: "Grande Poste", lat: 36.7697, lng: 3.0597, category: "Landmark" },
  { name: "Casbah", lat: 36.7833, lng: 3.0597, category: "Historic" },
  { name: "Riadh El Feth", lat: 36.7403, lng: 3.0508, category: "Shopping" },
];

interface Destination {
  name: string;
  lat: number;
  lng: number;
  category: string;
}

interface RouteData {
  from: { lat: number; lng: number };
  to: { lat: number; lng: number };
  path: { lat: number; lng: number }[];
  duration: number;
  distance: number;
  instructions: string[];
}

// Recent destinations with example data
const exampleRecentDestinations: Destination[] = [
  { name: "Alger Centre", lat: 36.7538, lng: 3.0588, category: "District" },
  {
    name: "Houari Boumediene Airport",
    lat: 36.691,
    lng: 3.2154,
    category: "Airport",
  },
  {
    name: "University of Algiers",
    lat: 36.7167,
    lng: 3.1833,
    category: "University",
  },
];

export default function MapView() {
  const { mapCenter, recentSearches, addRecentSearch } = useAppStore();

  const [showFilters, setShowFilters] = useState(false);
  const [currentPosition, setCurrentPosition] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const [searchFocused, setSearchFocused] = useState(false);
  const [destinationQuery, setDestinationQuery] = useState("");
  const [filteredDestinations, setFilteredDestinations] = useState<
    Destination[]
  >([]);
  const [selectedDestination, setSelectedDestination] =
    useState<Destination | null>(null);
  const [routeData, setRouteData] = useState<RouteData | null>(null);
  const [recentDestinations, setRecentDestinations] = useState<Destination[]>(
    []
  );
  const [isCalculatingRoute, setIsCalculatingRoute] = useState(false);

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

  useEffect(() => {
    // Load recent destinations from recent searches or use examples
    if (recentSearches.length > 0) {
      const recent = recentSearches
        .slice(0, 3)
        .map((search) => {
          const destination = popularDestinations.find((dest) =>
            dest.name.toLowerCase().includes(search.toLowerCase())
          );
          return destination || popularDestinations[0];
        })
        .filter(Boolean);
      setRecentDestinations(recent);
    } else {
      // Show example destinations when no recent searches
      setRecentDestinations(exampleRecentDestinations);
    }
  }, [recentSearches]);

  const handleDestinationSearch = (query: string) => {
    setDestinationQuery(query);
    if (query.trim()) {
      const filtered = popularDestinations.filter(
        (dest) =>
          dest.name.toLowerCase().includes(query.toLowerCase()) ||
          dest.category.toLowerCase().includes(query.toLowerCase())
      );
      setFilteredDestinations(filtered);
    } else {
      setFilteredDestinations([]);
    }
  };

  const calculateRoute = async (destination: Destination) => {
    if (!currentPosition) return;

    setIsCalculatingRoute(true);

    // Simulate route calculation (in real app, you'd use a routing service)
    setTimeout(() => {
      const route: RouteData = {
        from: currentPosition,
        to: { lat: destination.lat, lng: destination.lng },
        path: [
          currentPosition,
          {
            lat:
              currentPosition.lat +
              (destination.lat - currentPosition.lat) * 0.3,
            lng:
              currentPosition.lng +
              (destination.lng - currentPosition.lng) * 0.3,
          },
          {
            lat:
              currentPosition.lat +
              (destination.lat - currentPosition.lat) * 0.7,
            lng:
              currentPosition.lng +
              (destination.lng - currentPosition.lng) * 0.7,
          },
          { lat: destination.lat, lng: destination.lng },
        ],
        duration: Math.floor(Math.random() * 30) + 15, // 15-45 minutes
        distance: Math.floor(Math.random() * 20) + 5, // 5-25 km
        instructions: [
          "Head north on your current street",
          "Turn right onto main road",
          "Continue straight for 2.5 km",
          "Turn left at the roundabout",
          "Your destination will be on the right",
        ],
      };
      setRouteData(route);
      setIsCalculatingRoute(false);
    }, 1500);
  };

  const handleDestinationSelect = (destination: Destination) => {
    setSelectedDestination(destination);
    setDestinationQuery(destination.name);
    setSearchFocused(false);
    addRecentSearch(destination.name);

    // Calculate route from current position
    if (currentPosition) {
      calculateRoute(destination);
    }
  };

  const handleRecentDestinationSelect = (destination: Destination) => {
    handleDestinationSelect(destination);
  };

  const clearRoute = () => {
    setSelectedDestination(null);
    setRouteData(null);
    setDestinationQuery("");
  };

  return (
    <div className="h-screen flex flex-col overflow-hidden">
      {/* Header with Search */}
      <div className="bg-white shadow-sm border-b p-4 flex-shrink-0">
        <div className="space-y-3">
          {/* Current Location */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <Target size={16} className="text-green-500" />
              <span>
                {currentPosition
                  ? "Current location detected"
                  : "Detecting location..."}
              </span>
            </div>
            {selectedDestination && (
              <button
                onClick={clearRoute}
                className="text-sm text-primary hover:text-primary/80"
              >
                Clear route
              </button>
            )}
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

          {/* Search Results Dropdown */}
          {(searchFocused || destinationQuery) && (
            <div className="absolute top-full left-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg z-10 max-h-64 overflow-y-auto">
              {destinationQuery && filteredDestinations.length > 0 && (
                <div className="p-2">
                  <div className="text-xs text-gray-500 mb-2 px-2">
                    Search Results
                  </div>
                  {filteredDestinations
                    .slice(0, 5)
                    .map((destination, index) => (
                      <button
                        key={index}
                        onClick={() => handleDestinationSelect(destination)}
                        className="w-full text-left p-3 hover:bg-gray-50 rounded-lg transition-colors"
                      >
                        <div className="flex items-center gap-3">
                          <MapPin size={16} className="text-gray-400" />
                          <div>
                            <div className="font-medium text-gray-900">
                              {destination.name}
                            </div>
                            <div className="text-sm text-gray-500">
                              {destination.category}
                            </div>
                          </div>
                        </div>
                      </button>
                    ))}
                </div>
              )}

              {!destinationQuery && (
                <div className="p-2">
                  <div className="text-xs text-gray-500 mb-2 px-2">
                    Popular Destinations
                  </div>
                  {popularDestinations.slice(0, 5).map((destination, index) => (
                    <button
                      key={index}
                      onClick={() => handleDestinationSelect(destination)}
                      className="w-full text-left p-3 hover:bg-gray-50 rounded-lg transition-colors"
                    >
                      <div className="flex items-center gap-3">
                        <MapPin size={16} className="text-gray-400" />
                        <div>
                          <div className="font-medium text-gray-900">
                            {destination.name}
                          </div>
                          <div className="text-sm text-gray-500">
                            {destination.category}
                          </div>
                        </div>
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Route Info */}
      {routeData && (
        <div className="bg-blue-50 border-b p-4 flex-shrink-0">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Route size={20} className="text-blue-600" />
              <div>
                <div className="font-medium text-gray-900">
                  {routeData.duration} min • {routeData.distance} km
                </div>
                <div className="text-sm text-gray-600">
                  Route to {selectedDestination?.name}
                </div>
              </div>
            </div>
            {isCalculatingRoute && (
              <div className="text-sm text-blue-600">Calculating...</div>
            )}
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

            {/* Destination Marker */}
            {selectedDestination && (
              <Marker
                position={[selectedDestination.lat, selectedDestination.lng]}
              >
                <Popup>{selectedDestination.name}</Popup>
              </Marker>
            )}

            {/* Route Path */}
            {routeData && (
              <Polyline
                positions={routeData.path.map((point) => [
                  point.lat,
                  point.lng,
                ])}
                color="#3B82F6"
                weight={4}
                opacity={0.8}
              />
            )}
          </MapContainer>
        </div>

        {/* Recent Destinations Panel - Desktop */}
        <div className="w-80 bg-white border-l border-gray-200 flex flex-col lg:flex hidden">
          <div className="p-4 border-b flex-shrink-0">
            <h2 className="text-lg font-semibold text-gray-900">
              Recent Destinations
            </h2>
          </div>
          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {recentDestinations.length > 0 ? (
              recentDestinations.map((destination, index) => (
                <button
                  key={index}
                  onClick={() => handleRecentDestinationSelect(destination)}
                  className="w-full p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors text-left"
                >
                  <div className="flex items-center gap-3">
                    <Clock size={16} className="text-gray-400" />
                    <div>
                      <div className="font-medium text-gray-900">
                        {destination.name}
                      </div>
                      <div className="text-sm text-gray-500">
                        {destination.category}
                      </div>
                    </div>
                  </div>
                </button>
              ))
            ) : (
              <div className="text-center py-8">
                <Clock size={48} className="text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">No recent destinations</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mobile Recent Destinations - Bottom Sheet */}
      <div
        className="lg:hidden bg-white border-t flex-shrink-0"
        style={{ maxHeight: "40vh" }}
      >
        <div className="p-4 border-b">
          <h2 className="text-lg font-semibold text-gray-900">
            Recent Destinations
          </h2>
        </div>
        <div className="overflow-y-auto p-4 space-y-3">
          {recentDestinations.length > 0 ? (
            recentDestinations.slice(0, 3).map((destination, index) => (
              <button
                key={index}
                onClick={() => handleRecentDestinationSelect(destination)}
                className="w-full p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors text-left"
              >
                <div className="flex items-center gap-3">
                  <Clock size={16} className="text-gray-400" />
                  <div>
                    <div className="font-medium text-gray-900">
                      {destination.name}
                    </div>
                    <div className="text-sm text-gray-500">
                      {destination.category}
                    </div>
                  </div>
                </div>
              </button>
            ))
          ) : (
            <div className="text-center py-6">
              <Clock size={32} className="text-gray-400 mx-auto mb-2" />
              <p className="text-gray-500 text-sm">No recent destinations</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
