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
  Bus,
  Train,
  Zap,
  Navigation2,
} from "lucide-react";
import { useAppStore } from "../store/useAppStore";
import "leaflet/dist/leaflet.css";
import { RoutingService } from "../services/routingService";
import {
  algiersPlaces,
  searchAlgiersPlaces,
  type AlgiersPlace,
} from "../data/algiersPlaces";

// Real destinations in Algiers
const popularDestinations = [
  { name: "Alger Centre", lat: 36.7538, lng: 3.0588, category: "District" },
  { name: "Bab Ezzouar", lat: 36.7167, lng: 3.1833, category: "District" },
  { name: "Hydra", lat: 36.7669, lng: 3.0347, category: "District" },
  { name: "Kouba", lat: 36.7333, lng: 3.0833, category: "District" },
  { name: "Birtouta", lat: 36.6167, lng: 3.0167, category: "District" },
  { name: "El Harrach", lat: 36.7167, lng: 3.1333, category: "District" },
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
  { name: "El Harrach", lat: 36.7167, lng: 3.1333, category: "District" },
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
  shouldFitBounds = true,
}: {
  start: { lat: number; lng: number } | null;
  end: { lat: number; lng: number } | null;
  onRouteFound: (route: RouteData) => void;
  shouldFitBounds?: boolean;
}) {
  const map = useMap();
  const routingControlRef = useRef<L.Routing.Control | null>(null);
  const hasFittedBoundsRef = useRef(false);

  useEffect(() => {
    if (!start || !end) {
      // Remove existing routing control if no start/end points
      if (routingControlRef.current) {
        map.removeControl(routingControlRef.current);
        routingControlRef.current = null;
      }
      hasFittedBoundsRef.current = false;
      return;
    }

    // Remove existing routing control
    if (routingControlRef.current) {
      map.removeControl(routingControlRef.current);
    }

    // Reset the bounds fitting flag when start/end changes
    hasFittedBoundsRef.current = false;

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

        // Only fit bounds once when the route is first calculated, not on subsequent renders
        if (shouldFitBounds && !hasFittedBoundsRef.current) {
          const bounds = L.latLngBounds(route.coordinates);
          map.fitBounds(bounds, { padding: [50, 50] });
          hasFittedBoundsRef.current = true;
        }
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
  }, [map, start, end, onRouteFound, shouldFitBounds]);

  return null;
}

export default function MapView() {
  const { mapCenter, recentSearches, addRecentSearch } = useAppStore();

  const [showFilters, setShowFilters] = useState(false);
  const [currentPosition, setCurrentPosition] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const [destinationQuery, setDestinationQuery] = useState("");
  const [destinationSuggestions, setDestinationSuggestions] = useState<
    AlgiersPlace[]
  >([]);
  const [showDestinationSuggestions, setShowDestinationSuggestions] =
    useState(false);
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
  const [shouldFitBounds, setShouldFitBounds] = useState(true);

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

    if (query.trim().length >= 2) {
      const suggestions = searchAlgiersPlaces(query, 8);
      setDestinationSuggestions(suggestions);
    } else {
      setDestinationSuggestions([]);
    }
  };

  const handleDestinationSuggestionSelect = (place: AlgiersPlace) => {
    const destination: Destination = {
      name: place.name,
      lat: place.lat,
      lng: place.lng,
      category: place.category,
    };

    setDestinationQuery(place.name);
    setDestinationSuggestions([]);
    setShowDestinationSuggestions(false);
    handleDestinationSelect(destination);
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
    const isAirport = destination.name.includes("Airport");
    const isUniversity = destination.category === "University";
    const isCentre = destination.name.includes("Centre");
    const isElHarrach = destination.name.includes("El Harrach");

    // Calculate realistic base distance
    let baseDistance = 5; // Default km
    if (currentPosition) {
      const R = 6371; // Earth's radius in km
      const dLat = ((destination.lat - currentPosition.lat) * Math.PI) / 180;
      const dLng = ((destination.lng - currentPosition.lng) * Math.PI) / 180;
      const a =
        Math.sin(dLat / 2) * Math.sin(dLat / 2) +
        Math.cos((currentPosition.lat * Math.PI) / 180) *
          Math.cos((destination.lat * Math.PI) / 180) *
          Math.sin(dLng / 2) *
          Math.sin(dLng / 2);
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
      baseDistance = R * c;
    }

    const options: RouteOption[] = [];

    // Metro + Walking option (realistic for metro-connected areas)
    if (isCentre || isUniversity || isElHarrach) {
      const metroTime = isCentre ? 15 : isElHarrach ? 22 : 25; // Realistic metro times
      const walkingTime = isCentre ? 8 : isElHarrach ? 6 : 10;
      const totalTime = metroTime + walkingTime + 5; // +5 for waiting/transfers

      options.push({
        id: "metro-walk",
        name: "Metro + Walking",
        duration: totalTime,
        cost: 50, // Standard metro fare in Algiers
        transportModes: ["Metro", "Walk"],
        description: isElHarrach
          ? "Metro Line 1 directly serves El Harrach station"
          : isCentre
          ? "Quick metro access to city center"
          : "Fastest route using metro system",
        icon: Train,
        color: "bg-blue-500",
        steps: isElHarrach
          ? [
              "Walk to nearest metro station (4 min)",
              "Take Metro Line 1 towards El Harrach (22 min)",
              "Exit at El Harrach metro station (1 min)",
              "Walk to final destination (6 min)",
            ]
          : isCentre
          ? [
              "Walk to nearest metro station (3 min)",
              "Take Metro Line 1 to city center (15 min)",
              "Walk to destination (8 min)",
            ]
          : [
              "Walk to nearest metro station (5 min)",
              "Take Metro Line 1 (25 min)",
              "Walk to destination (10 min)",
            ],
      });
    }

    // Bus option (enhanced with realistic times for Algiers)
    const busTime = Math.max(20, Math.floor(baseDistance * 2.8)); // Realistic bus speed in Algiers traffic
    const busWalkTime = isCentre ? 5 : isElHarrach ? 4 : 7;

    options.push({
      id: "bus",
      name: "Bus Direct",
      duration: busTime + busWalkTime,
      cost: 25, // Standard bus fare in Algiers
      transportModes: ["Bus"],
      description: isElHarrach
        ? "Bus lines 15, 22, and 35 serve El Harrach"
        : isCentre
        ? "Multiple bus lines to city center"
        : "Direct bus route, affordable option",
      icon: Bus,
      color: "bg-green-500",
      steps: isElHarrach
        ? [
            "Walk to bus stop (2 min)",
            "Take Bus 15, 22, or 35 to El Harrach (28 min)",
            "Walk to destination (4 min)",
          ]
        : isCentre
        ? [
            "Walk to bus stop (3 min)",
            "Take Bus 23 or 64 to Alger Centre (18 min)",
            "Walk to destination (5 min)",
          ]
        : [
            "Walk to bus stop (4 min)",
            "Take direct bus route (35 min)",
            "Walk from bus stop (7 min)",
          ],
    });

    // Tramway option (realistic for Algiers tram network)
    if (
      isCentre ||
      destination.name.includes("Hydra") ||
      destination.name.includes("Kouba") ||
      isElHarrach
    ) {
      const tramTime = isCentre ? 12 : isElHarrach ? 35 : 25; // Tram doesn't directly serve El Harrach
      const tramWalkTime = isCentre ? 6 : isElHarrach ? 12 : 8; // Longer walk for El Harrach
      const transferTime = isElHarrach ? 8 : 0; // Transfer needed for El Harrach

      options.push({
        id: "tram-walk",
        name: isElHarrach ? "Tram + Bus" : "Tramway + Walking",
        duration: tramTime + tramWalkTime + transferTime,
        cost: isElHarrach ? 65 : 40, // Tram + bus for El Harrach
        transportModes: isElHarrach
          ? ["Tram", "Bus", "Walk"]
          : ["Tram", "Walk"],
        description: isElHarrach
          ? "Tram to city center, then bus to El Harrach"
          : isCentre
          ? "Scenic tram route to city center"
          : "Comfortable tram journey",
        icon: Zap,
        color: "bg-yellow-500",
        steps: isElHarrach
          ? [
              "Walk to Tram Line A (5 min)",
              "Take tram to Tafourah (12 min)",
              "Transfer to Bus 22 (3 min)",
              "Take bus to El Harrach (20 min)",
              "Walk to destination (8 min)",
            ]
          : isCentre
          ? [
              "Walk to Tram Line A (4 min)",
              "Take tram to city center (12 min)",
              "Walk to destination (6 min)",
            ]
          : [
              "Walk to tram station (5 min)",
              "Take Tram Line A (25 min)",
              "Walk to destination (8 min)",
            ],
      });
    }

    // Combined transport option (realistic multi-modal for Algiers)
    if (!isAirport) {
      const combinedTime = isCentre ? 25 : isElHarrach ? 32 : 35;

      options.push({
        id: "combined",
        name: "Bus + Metro",
        duration: combinedTime,
        cost: 75, // Combined fare
        transportModes: ["Bus", "Metro"],
        description: isElHarrach
          ? "Bus to metro, then Metro Line 1 to El Harrach"
          : isCentre
          ? "Quick bus connection to metro system"
          : "Optimized multi-modal route",
        icon: Navigation2,
        color: "bg-purple-500",
        steps: isElHarrach
          ? [
              "Walk to bus stop (3 min)",
              "Take Bus 12 to Tafourah Metro (12 min)",
              "Transfer to Metro Line 1 (2 min)",
              "Take Metro to El Harrach station (18 min)",
              "Walk to destination (5 min)",
            ]
          : isCentre
          ? [
              "Walk to bus stop (2 min)",
              "Take Bus 11 to 1er Mai Metro (8 min)",
              "Transfer to Metro Line 1 (2 min)",
              "Take Metro to city center (10 min)",
              "Walk to destination (5 min)",
            ]
          : [
              "Walk to bus stop (3 min)",
              "Take Bus 12 to metro connection (15 min)",
              "Transfer to Metro Line 1 (2 min)",
              "Take Metro to destination area (18 min)",
              "Walk to destination (7 min)",
            ],
      });
    }

    // VTC Options - Grouped into one card with realistic pricing
    const vtcTime = Math.max(12, Math.floor(baseDistance * 2.0)); // VTC faster than bus due to direct route

    // Realistic VTC pricing based on actual data for specific destinations
    let yassirCost: number;
    let heetchCost: number;
    let inDriveCost: number;
    let vtcEta: number;

    if (isElHarrach) {
      // Pricing data for El Harrach (the more expensive destination)
      yassirCost = 871;
      heetchCost = 900;
      inDriveCost = 870; // Using middle of range (725-1015)
      vtcEta = 15;
    } else if (isCentre) {
      // Pricing data for Alger Centre (the less expensive destination)
      yassirCost = 341;
      heetchCost = 352;
      inDriveCost = 340; // Using middle of range (284-397)
      vtcEta = 8;
    } else {
      // Fallback calculation for other destinations
      const baseFare = 150; // Starting fare in DA
      const perKmRate = 60; // Rate per km in DA
      const surgeMultiplier = isAirport ? 1.4 : 1.0; // Airport surge
      const timeOfDay = new Date().getHours();
      const rushHourMultiplier =
        (timeOfDay >= 7 && timeOfDay <= 9) ||
        (timeOfDay >= 17 && timeOfDay <= 19)
          ? 1.3
          : 1.0;

      yassirCost = Math.round(
        (baseFare + baseDistance * perKmRate) *
          surgeMultiplier *
          rushHourMultiplier *
          1.1
      );
      heetchCost = Math.round(
        (baseFare + baseDistance * perKmRate) *
          surgeMultiplier *
          rushHourMultiplier *
          0.95
      );
      inDriveCost = Math.round(
        (baseFare + baseDistance * perKmRate) *
          surgeMultiplier *
          rushHourMultiplier *
          0.85
      );
      vtcEta = vtcTime + 3;
    }

    // Use the lowest price for the main display
    const minVtcCost = Math.min(yassirCost, heetchCost, inDriveCost);

    options.push({
      id: "vtc-options",
      name: "VTC Services",
      duration: vtcEta, // Use realistic ETA
      cost: minVtcCost,
      transportModes: ["VTC"],
      description: isAirport
        ? "Private rides to airport with multiple app options"
        : isElHarrach
        ? "Door-to-door service to El Harrach with competitive pricing"
        : isCentre
        ? "Quick rides to city center with multiple options"
        : "Choose from Yassir, Heetch, or InDrive for door-to-door service",
      icon: Navigation2,
      color: "bg-emerald-500",
      steps: [
        `Choose your preferred VTC app: Yassir: ${yassirCost} DA / Heetch: ${heetchCost} DA / InDrive: ${inDriveCost} DA`,
        `Wait for pickup (${vtcEta < 10 ? "3-5" : "5-8"} min)`,
        `Direct ride to ${destination.name} (${vtcEta - 3} min)`,
        "Pay through app or cash (depending on service)",
      ],
    });

    // Walking option (for nearby destinations only)
    if (baseDistance < 5) {
      const walkingTime = Math.floor(baseDistance * 12); // 5 km/h walking speed

      options.push({
        id: "walking",
        name: "Walking",
        duration: walkingTime,
        cost: 0,
        transportModes: ["Walk"],
        description: "Free, healthy option for short distances",
        icon: Navigation2,
        color: "bg-gray-500",
        steps: [
          "Walk directly via main streets",
          `Estimated walking time: ${walkingTime} minutes`,
        ],
      });
    }

    // Airport express (only for airport destinations)
    if (isAirport) {
      options.push({
        id: "airport-express",
        name: "Airport Express Bus",
        duration: 45,
        cost: 200, // Premium airport service
        transportModes: ["Airport Bus"],
        description: "Direct airport shuttle from city center",
        icon: Navigation2,
        color: "bg-red-500",
        steps: [
          "Walk to airport shuttle stop (5 min)",
          "Take Airport Express Bus (35 min)",
          "Arrive at terminal (5 min)",
        ],
      });
    }

    return options.sort((a, b) => a.duration - b.duration);
  };

  const handleDestinationSelect = (destination: Destination) => {
    setSelectedDestination(destination);
    addRecentSearch(destination.name);
    setShouldFitBounds(true); // Allow bounds fitting for new destination

    // Generate route options
    const options = generateRouteOptions(destination);
    setRouteOptions(options);
    setSelectedRouteOption(options[0]); // Auto-select fastest option

    // Calculate route for the first option
    if (currentPosition) {
      generateRouteForOption(options[0], destination);
    }
  };

  const calculateRouteForOption = async (destination: Destination) => {
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

  const handleRouteOptionSelect = (option: RouteOption) => {
    setSelectedRouteOption(option);
    setShouldFitBounds(false); // Don't fit bounds when switching between route options

    // Generate new route for the selected option
    if (selectedDestination) {
      generateRouteForOption(selectedDestination);
    }
  };

  const clearRoute = () => {
    setSelectedDestination(null);
    setRouteData(null);
    setDestinationQuery("");
    setDestinationSuggestions([]);
    setShowDestinationSuggestions(false);
    setRouteOptions([]);
    setSelectedRouteOption(null);
    setShouldFitBounds(true); // Reset for next route
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
      <div className="flex-1 flex flex-col gap-x-4 overflow-hidden max-w-none lg:max-w-[70%]">
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

        {/* Map Container - Full height on desktop, adjusted on mobile */}
        <div className="h-[calc(50vh-2rem)] lg:h-full relative mb-4 lg:mb-0">
          <MapContainer
            center={mapCenter}
            zoom={12}
            className="h-full w-full rounded-lg lg:rounded-none"
          >
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
              shouldFitBounds={shouldFitBounds}
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
      <div className="flex-1 lg:min-w-[30%] lg:max-w-[400px] bg-white border-l border-gray-200 flex flex-col mb-4 lg:mb-0">
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

            {/* Destination Search - Exactly like SearchPage */}
            <div className="relative">
              <MapPin
                size={20}
                className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              />
              <input
                type="text"
                placeholder="Where do you want to go?"
                value={destinationQuery}
                onChange={(e) => handleDestinationSearch(e.target.value)}
                onFocus={() => setShowDestinationSuggestions(true)}
                onBlur={() =>
                  setTimeout(() => setShowDestinationSuggestions(false), 200)
                }
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

              {/* Destination Suggestions - Exactly like SearchPage */}
              {showDestinationSuggestions &&
                destinationSuggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 bg-white border border-gray-200 rounded-lg shadow-lg z-10 mt-1">
                    {destinationSuggestions.map((place) => (
                      <button
                        key={place.id}
                        onClick={() => handleDestinationSuggestionSelect(place)}
                        className="w-full text-left p-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0 transition-colors"
                      >
                        <div className="flex items-center gap-3">
                          <MapPin size={16} className="text-gray-400" />
                          <div className="flex-1">
                            <div className="font-medium text-gray-900">
                              {place.name}
                            </div>
                            <div className="text-sm text-gray-500 flex items-center gap-2">
                              <span>{place.category}</span>
                              <span>•</span>
                              <span>{place.district}</span>
                              {place.arabicName && (
                                <>
                                  <span>•</span>
                                  <span className="text-xs">
                                    {place.arabicName}
                                  </span>
                                </>
                              )}
                            </div>
                          </div>
                        </div>
                      </button>
                    ))}
                  </div>
                )}

              {/* Show popular destinations when focused but empty - like SearchPage */}
              {showDestinationSuggestions && !destinationQuery.trim() && (
                <div className="absolute top-full left-0 right-0 bg-white border border-gray-200 rounded-lg shadow-lg z-10 mt-1">
                  <div className="p-3 border-b border-gray-100">
                    <div className="text-xs text-gray-500 mb-2">
                      Popular Destinations
                    </div>
                  </div>
                  {algiersPlaces
                    .sort((a, b) => b.popularity - a.popularity)
                    .slice(0, 5)
                    .map((place) => (
                      <button
                        key={place.id}
                        onClick={() => handleDestinationSuggestionSelect(place)}
                        className="w-full text-left p-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0 transition-colors"
                      >
                        <div className="flex items-center gap-3">
                          <MapPin size={16} className="text-gray-400" />
                          <div>
                            <div className="font-medium text-gray-900">
                              {place.name}
                            </div>
                            <div className="text-sm text-gray-500">
                              {place.category} • {place.district}
                            </div>
                          </div>
                        </div>
                      </button>
                    ))}
                </div>
              )}

              {/* No results message - like SearchPage */}
              {showDestinationSuggestions &&
                destinationQuery.trim() &&
                destinationQuery.length >= 2 &&
                destinationSuggestions.length === 0 && (
                  <div className="absolute top-full left-0 right-0 bg-white border border-gray-200 rounded-lg shadow-lg z-10 mt-1">
                    <div className="p-4 text-center text-gray-500">
                      <div className="text-sm">
                        No places found for "{destinationQuery}"
                      </div>
                      <div className="text-xs mt-1">
                        Try searching for districts, landmarks, or popular
                        places
                      </div>
                    </div>
                  </div>
                )}
            </div>
          </div>
        </div>

        {/* Route Options or Recent Destinations */}
        <div className="flex-1 overflow-y-auto pb-4">
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
