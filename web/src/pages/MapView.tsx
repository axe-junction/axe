import { useEffect, useState, useRef } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  useMap,
  Polyline,
} from "react-leaflet";
import L from "leaflet";
import "leaflet-routing-machine";
import {
  Search,
  Filter,
  Navigation,
  MapPin,
  Clock,
  Target,
  Route,
  Bus,
  Train,
  Zap,
  Navigation2,
  DollarSign,
  ArrowRight,
} from "lucide-react";
import { useAppStore } from "../store/useAppStore";
import "leaflet/dist/leaflet.css";
import { RoutingService } from "../services/routingService";

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

interface RouteOption {
  id: string;
  name: string;
  duration: number;
  cost: number;
  transportModes: string[];
  description: string;
  icon: typeof Bus | typeof Train | typeof Zap | typeof Navigation2;
  color: string;
  steps: string[];
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

// Custom component to handle routing on the map
function RoutingMachine({
  start,
  end,
  onRouteFound,
}: {
  start: { lat: number; lng: number } | null;
  end: { lat: number; lng: number } | null;
  onRouteFound: (route: RouteData) => void;
}) {
  const map = useMap();
  const routingControlRef = useRef<L.Routing.Control | null>(null);

  useEffect(() => {
    if (!start || !end) {
      // Remove existing routing control if no start/end points
      if (routingControlRef.current) {
        map.removeControl(routingControlRef.current);
        routingControlRef.current = null;
      }
      return;
    }

    // Remove existing routing control
    if (routingControlRef.current) {
      map.removeControl(routingControlRef.current);
    }

    // Create new routing control with proper styling
    const routingControl = L.Routing.control({
      waypoints: [L.latLng(start.lat, start.lng), L.latLng(end.lat, end.lng)],
      routeWhileDragging: false,
      addWaypoints: false,
      createMarker: () => null, // Don't create markers (we handle them separately)
      lineOptions: {
        styles: [
          {
            color: "#3B82F6",
            weight: 6,
            opacity: 1,
            dashArray: "",
            lineCap: "round",
            lineJoin: "round",
          },
        ],
        extendToWaypoints: true,
        missingRouteTolerance: 0,
      },
      show: false, // Don't show the control panel
      router: L.Routing.osrmv1({
        serviceUrl: "https://router.project-osrm.org/route/v1",
        profile: "driving",
      }),
    });

    // Add the control to the map
    routingControl.addTo(map);

    // Listen for route found event
    routingControl.on("routesfound", (e: any) => {
      const routes = e.routes;
      if (routes && routes.length > 0) {
        const route = routes[0];

        const routeData: RouteData = {
          from: start,
          to: end,
          path: route.coordinates.map((coord: L.LatLng) => ({
            lat: coord.lat,
            lng: coord.lng,
          })),
          duration: Math.round(route.summary.totalTime / 60), // convert to minutes
          distance: Math.round(route.summary.totalDistance / 1000), // convert to km
          instructions: route.instructions.map(
            (instruction: any) =>
              instruction.text || instruction.instruction || "Continue"
          ),
        };

        onRouteFound(routeData);

        // Fit map to route bounds with padding
        const bounds = L.latLngBounds(route.coordinates);
        map.fitBounds(bounds, { padding: [50, 50] });
      }
    });

    // Handle routing errors
    routingControl.on("routingerror", (e: any) => {
      console.error("Routing error:", e);
    });

    routingControlRef.current = routingControl;

    return () => {
      if (routingControlRef.current) {
        map.removeControl(routingControlRef.current);
        routingControlRef.current = null;
      }
    };
  }, [map, start, end, onRouteFound]);

  return null;
}

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
  const [routingService] = useState(() => RoutingService.getInstance());
  const [routeOptions, setRouteOptions] = useState<RouteOption[]>([]);
  const [selectedRouteOption, setSelectedRouteOption] =
    useState<RouteOption | null>(null);

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

    try {
      const route = await routingService.calculateRoute(currentPosition, {
        lat: destination.lat,
        lng: destination.lng,
      });

      const routeData: RouteData = {
        from: currentPosition,
        to: { lat: destination.lat, lng: destination.lng },
        path: route.path,
        duration: route.duration,
        distance: route.distance,
        instructions: route.instructions,
      };

      setRouteData(routeData);
    } catch (error) {
      console.error("Failed to calculate route:", error);

      // Create a simple fallback route
      const fallbackRoute: RouteData = {
        from: currentPosition,
        to: { lat: destination.lat, lng: destination.lng },
        path: [currentPosition, { lat: destination.lat, lng: destination.lng }],
        duration: Math.floor(Math.random() * 30) + 15,
        distance: Math.floor(Math.random() * 20) + 5,
        instructions: ["Head towards your destination", "You have arrived"],
      };

      setRouteData(fallbackRoute);
    } finally {
      setIsCalculatingRoute(false);
    }
  };

  const handleRouteFound = (route: RouteData) => {
    setRouteData(route);
    setIsCalculatingRoute(false);
  };

  const generateRouteOptions = (destination: Destination): RouteOption[] => {
    const baseDistance = Math.random() * 20 + 5; // 5-25 km
    const isAirport = destination.name.includes("Airport");
    const isUniversity = destination.category === "University";
    const isCentre = destination.name.includes("Centre");

    const options: RouteOption[] = [];

    // Metro + Walking option
    if (isCentre || isUniversity) {
      options.push({
        id: "metro-walk",
        name: "Metro + Walking",
        duration: Math.floor(baseDistance * 2 + 10),
        cost: 50,
        transportModes: ["Metro", "Walk"],
        description: "Fastest route using metro system",
        icon: Train,
        color: "bg-blue-500",
        steps: [
          "Walk to nearest metro station (5 min)",
          "Take Metro Line 1 (25 min)",
          "Walk to destination (8 min)",
        ],
      });
    }

    // Bus option
    options.push({
      id: "bus",
      name: "Bus Direct",
      duration: Math.floor(baseDistance * 3 + 15),
      cost: 30,
      transportModes: ["Bus"],
      description: "Direct bus route, more affordable",
      icon: Bus,
      color: "bg-green-500",
      steps: [
        "Walk to bus stop (3 min)",
        "Take Bus 23 direct route (35 min)",
        "Walk from bus stop (5 min)",
      ],
    });

    // Tramway option (for certain destinations)
    if (
      isCentre ||
      destination.name.includes("Hydra") ||
      destination.name.includes("Kouba")
    ) {
      options.push({
        id: "tram-walk",
        name: "Tramway + Walking",
        duration: Math.floor(baseDistance * 2.5 + 8),
        cost: 40,
        transportModes: ["Tram", "Walk"],
        description: "Scenic route via tramway",
        icon: Zap,
        color: "bg-yellow-500",
        steps: [
          "Walk to tram station (4 min)",
          "Take Tram Line A (20 min)",
          "Walk to destination (6 min)",
        ],
      });
    }

    // Combined transport option
    if (!isAirport) {
      options.push({
        id: "combined",
        name: "Bus + Metro",
        duration: Math.floor(baseDistance * 2.2 + 12),
        cost: 70,
        transportModes: ["Bus", "Metro"],
        description: "Optimized multi-modal route",
        icon: Navigation2,
        color: "bg-purple-500",
        steps: [
          "Walk to bus stop (3 min)",
          "Take Bus 12 (15 min)",
          "Transfer to Metro at Tafourah (2 min)",
          "Take Metro Line 1 (18 min)",
          "Walk to destination (7 min)",
        ],
      });
    }

    // Airport express (for airport destinations)
    if (isAirport) {
      options.push({
        id: "airport-express",
        name: "Airport Express",
        duration: Math.floor(baseDistance * 1.5 + 5),
        cost: 150,
        transportModes: ["Bus"],
        description: "Direct airport shuttle service",
        icon: Navigation2,
        color: "bg-red-500",
        steps: [
          "Walk to shuttle stop (2 min)",
          "Take Airport Express (35 min)",
          "Arrive at terminal (3 min)",
        ],
      });
    }

    // Walking option (for nearby destinations)
    if (baseDistance < 8) {
      options.push({
        id: "walking",
        name: "Walking",
        duration: Math.floor(baseDistance * 12),
        cost: 0,
        transportModes: ["Walk"],
        description: "Free, healthy option",
        icon: Navigation2,
        color: "bg-gray-500",
        steps: [
          "Walk directly to destination",
          "Estimated walking time based on distance",
        ],
      });
    }

    return options.sort((a, b) => a.duration - b.duration);
  };

  const handleDestinationSelect = (destination: Destination) => {
    setSelectedDestination(destination);
    setDestinationQuery(destination.name);
    setSearchFocused(false);
    addRecentSearch(destination.name);

    // Generate route options
    const options = generateRouteOptions(destination);
    setRouteOptions(options);
    setSelectedRouteOption(options[0]); // Auto-select fastest option

    // Calculate route for the first option
    if (currentPosition) {
      generateRouteForOption(options[0], destination);
    }
  };

  const generateRouteForOption = async (
    option: RouteOption,
    destination: Destination
  ) => {
    if (!currentPosition) return;

    setIsCalculatingRoute(true);

    try {
      let routePath: { lat: number; lng: number }[] = [];

      // Generate different routes based on transport mode
      if (option.id === "metro-walk") {
        // Route via metro stations
        const metroStations = [
          { lat: 36.7538, lng: 3.0588 }, // Alger Centre (metro hub)
          { lat: 36.7403, lng: 3.0508 }, // Tafourah station
        ];
        routePath = [
          currentPosition,
          metroStations[0],
          metroStations[1],
          { lat: destination.lat, lng: destination.lng },
        ];
      } else if (option.id === "tram-walk") {
        // Route via tram lines
        const tramStations = [
          { lat: 36.7669, lng: 3.0347 }, // Hydra tram
          { lat: 36.7525, lng: 3.0519 }, // Maqam Echahid
        ];
        routePath = [
          currentPosition,
          tramStations[0],
          tramStations[1],
          { lat: destination.lat, lng: destination.lng },
        ];
      } else if (option.id === "bus") {
        // Direct bus route with fewer stops
        const busStops = [
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
        ];
        routePath = [
          currentPosition,
          ...busStops,
          { lat: destination.lat, lng: destination.lng },
        ];
      } else if (option.id === "combined") {
        // Bus + Metro combination
        const transferPoint = { lat: 36.7403, lng: 3.0508 }; // Tafourah transfer
        const busStop = {
          lat:
            currentPosition.lat +
            (transferPoint.lat - currentPosition.lat) * 0.5,
          lng:
            currentPosition.lng +
            (transferPoint.lng - currentPosition.lng) * 0.5,
        };
        routePath = [
          currentPosition,
          busStop,
          transferPoint,
          { lat: destination.lat, lng: destination.lng },
        ];
      } else if (option.id === "airport-express") {
        // Direct highway route to airport
        const highwayPoint = { lat: 36.7, lng: 3.15 }; // Highway junction
        routePath = [
          currentPosition,
          highwayPoint,
          { lat: destination.lat, lng: destination.lng },
        ];
      } else if (option.id === "walking") {
        // Walking route - more direct path
        const midPoint = {
          lat: (currentPosition.lat + destination.lat) / 2,
          lng: (currentPosition.lng + destination.lng) / 2,
        };
        routePath = [
          currentPosition,
          midPoint,
          { lat: destination.lat, lng: destination.lng },
        ];
      } else {
        // Fallback to original routing service
        const route = await routingService.calculateRoute(currentPosition, {
          lat: destination.lat,
          lng: destination.lng,
        });
        routePath = route.path;
      }

      const routeData: RouteData = {
        from: currentPosition,
        to: { lat: destination.lat, lng: destination.lng },
        path: routePath,
        duration: option.duration,
        distance: option.cost / 10, // Approximate distance from cost
        instructions: option.steps,
      };

      setRouteData(routeData);
    } catch (error) {
      console.error("Failed to calculate route:", error);

      // Create a simple fallback route specific to the option
      const fallbackRoute: RouteData = {
        from: currentPosition,
        to: { lat: destination.lat, lng: destination.lng },
        path: [currentPosition, { lat: destination.lat, lng: destination.lng }],
        duration: option.duration,
        distance: Math.floor(Math.random() * 20) + 5,
        instructions: option.steps,
      };

      setRouteData(fallbackRoute);
    } finally {
      setIsCalculatingRoute(false);
    }
  };

  const handleRouteOptionSelect = (option: RouteOption) => {
    setSelectedRouteOption(option);

    // Generate new route for the selected option
    if (selectedDestination) {
      generateRouteForOption(option, selectedDestination);
    }
  };

  const clearRoute = () => {
    setSelectedDestination(null);
    setRouteData(null);
    setDestinationQuery("");
    setRouteOptions([]);
    setSelectedRouteOption(null);
  };

  const getRouteColor = (routeOption: RouteOption | null): string => {
    if (!routeOption) return "#3B82F6";

    const colorMap: { [key: string]: string } = {
      "bg-blue-500": "#3B82F6",
      "bg-green-500": "#10B981",
      "bg-yellow-500": "#F59E0B",
      "bg-purple-500": "#8B5CF6",
      "bg-red-500": "#EF4444",
      "bg-gray-500": "#6B7280",
    };

    return colorMap[routeOption.color] || "#3B82F6";
  };

  return (
    <div className="h-full flex flex-col lg:flex-row overflow-hidden">
      {/* Main Content - Map */}
      <div className="flex-1 flex flex-col overflow-hidden max-w-none lg:max-w-[70%]">
        {/* Mobile Header - only on mobile */}
        <div className="lg:hidden bg-white shadow-sm border-b p-4 flex-shrink-0">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-900">Navigation</h1>
            {selectedDestination && (
              <button
                onClick={clearRoute}
                className="text-sm hover:text-gray-600"
                style={{ color: "#6316DB" }}
              >
                Clear route
              </button>
            )}
          </div>
        </div>

        {/* Route Summary - shows on both mobile and desktop */}
        {routeData && selectedRouteOption && (
          <div className="bg-blue-50 border-b p-4 flex-shrink-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div
                  className={`p-2 rounded-lg ${selectedRouteOption.color} text-white`}
                >
                  <selectedRouteOption.icon size={20} />
                </div>
                <div>
                  <div className="font-medium text-gray-900">
                    {selectedRouteOption.duration} min •{" "}
                    {selectedRouteOption.cost} DA
                  </div>
                  <div className="text-sm text-gray-600">
                    {selectedRouteOption.name} to {selectedDestination?.name}
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-2">
                {isCalculatingRoute && (
                  <div className="text-sm text-blue-600">Calculating...</div>
                )}
                <div className="flex gap-1">
                  {selectedRouteOption.transportModes.map((mode, index) => (
                    <span
                      key={index}
                      className="px-2 py-1 bg-white rounded-full text-xs text-gray-600"
                    >
                      {mode}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Map Container - Fixed height on mobile, flexible on desktop */}
        <div className="h-[50vh] lg:flex-1 relative">
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

            {/* Routing Machine Component */}
            <RoutingMachine
              start={currentPosition}
              end={
                selectedDestination
                  ? {
                      lat: selectedDestination.lat,
                      lng: selectedDestination.lng,
                    }
                  : null
              }
              onRouteFound={handleRouteFound}
            />

            {/* Route Polyline - Always visible when route data exists */}
            {routeData && routeData.path.length > 0 && (
              <Polyline
                positions={routeData.path.map((point) => [
                  point.lat,
                  point.lng,
                ])}
                pathOptions={{
                  color: getRouteColor(selectedRouteOption),
                  weight: 6,
                  opacity: 1,
                  dashArray: "",
                  lineCap: "round",
                  lineJoin: "round",
                }}
              />
            )}
          </MapContainer>
        </div>
      </div>

      {/* Right Panel - Search and Route Options */}
      <div className="flex-1 lg:min-w-[30%] lg:max-w-[400px] bg-white border-l border-gray-200 flex flex-col">
        {/* Search Section */}
        <div className="p-4 border-b flex-shrink-0">
          <div className="space-y-3">
            {/* Desktop Header */}
            <div className="hidden lg:flex items-center justify-between">
              <h1 className="text-xl font-bold text-gray-900">Navigation</h1>
              {selectedDestination && (
                <button
                  onClick={clearRoute}
                  className="text-sm text-primary hover:text-primary/80"
                >
                  Clear route
                </button>
              )}
            </div>

            {/* Current Location Status */}
            <div className="flex items-center gap-2 text-sm text-gray-600 bg-gray-50 p-3 rounded-lg">
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

            {/* Search Results Dropdown */}
            {(searchFocused || destinationQuery) && !selectedDestination && (
              <div className="absolute top-full left-0 right-0 bg-white border border-gray-200 rounded-lg shadow-lg z-10 max-h-64 overflow-y-auto mx-4">
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
                    {popularDestinations
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
              </div>
            )}
          </div>
        </div>

        {/* Route Options or Recent Destinations */}
        <div className="flex-1 overflow-y-auto">
          {selectedDestination && routeOptions.length > 0 ? (
            // Route Options
            <>
              <div className="p-4 border-b">
                <h2 className="text-lg font-semibold text-gray-900">
                  Route Options
                </h2>
                <p className="text-sm text-gray-600">
                  Choose your preferred route to {selectedDestination.name}
                </p>
              </div>
              <div className="p-4 space-y-3">
                {routeOptions.map((option) => (
                  <button
                    key={option.id}
                    onClick={() => handleRouteOptionSelect(option)}
                    className={`w-full p-4 rounded-lg border-2 transition-all text-left ${
                      selectedRouteOption?.id === option.id
                        ? "border-primary bg-primary/5 ring-2 ring-primary/20"
                        : "border-gray-200 hover:border-gray-300 hover:bg-gray-50"
                    }`}
                  >
                    <div className="flex items-start gap-3">
                      <div
                        className={`p-2 rounded-lg ${option.color} text-white flex-shrink-0`}
                      >
                        <option.icon size={16} />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-1">
                          <h3 className="font-medium text-gray-900">
                            {option.name}
                          </h3>
                          <div className="flex items-center gap-2">
                            <div className="flex items-center gap-1 text-sm text-gray-600">
                              <Clock size={12} />
                              {option.duration}m
                            </div>
                            <div className="flex items-center gap-1 text-sm text-gray-600">
                              <DollarSign size={12} />
                              {option.cost} DA
                            </div>
                          </div>
                        </div>
                        <p className="text-sm text-gray-600 mb-2">
                          {option.description}
                        </p>
                        <div className="flex gap-1 flex-wrap">
                          {option.transportModes.map((mode, index) => (
                            <span
                              key={index}
                              className="px-2 py-1 bg-gray-100 text-xs rounded-full text-gray-600"
                            >
                              {mode}
                            </span>
                          ))}
                        </div>
                        {selectedRouteOption?.id === option.id && (
                          <div className="mt-3 pt-3 border-t border-gray-200">
                            <h4 className="text-sm font-medium text-gray-900 mb-2">
                              Steps:
                            </h4>
                            <div className="space-y-1">
                              {option.steps.map((step, index) => (
                                <div
                                  key={index}
                                  className="flex items-start gap-2 text-sm text-gray-600"
                                >
                                  <span className="text-primary font-medium min-w-[20px]">
                                    {index + 1}.
                                  </span>
                                  <span>{step}</span>
                                </div>
                              ))}
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            </>
          ) : (
            // Recent Destinations
            <>
              <div className="p-4 border-b">
                <h2 className="text-lg font-semibold text-gray-900">
                  Recent Destinations
                </h2>
              </div>
              <div className="p-4 space-y-3">
                {recentDestinations.length > 0 ? (
                  recentDestinations.map((destination, index) => (
                    <button
                      key={index}
                      onClick={() => handleDestinationSelect(destination)}
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
            </>
          )}
        </div>
      </div>
    </div>
  );
}
