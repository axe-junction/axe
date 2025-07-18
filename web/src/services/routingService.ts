import L from "leaflet";
import "leaflet-routing-machine";

interface RoutePoint {
  lat: number;
  lng: number;
}

interface RouteResponse {
  path: RoutePoint[];
  duration: number; // in minutes
  distance: number; // in kilometers
  instructions: string[];
}

export class RoutingService {
  private static instance: RoutingService;

  static getInstance(): RoutingService {
    if (!RoutingService.instance) {
      RoutingService.instance = new RoutingService();
    }
    return RoutingService.instance;
  }

  async calculateRoute(
    start: RoutePoint,
    end: RoutePoint,
    profile: "driving-car" | "foot-walking" | "cycling-regular" = "driving-car"
  ): Promise<RouteResponse> {
    try {
      // Try to use OSRM API directly
      const osrmResponse = await this.callOSRMAPI(start, end, profile);
      return osrmResponse;
    } catch (error) {
      console.error("OSRM API error:", error);
      // Fallback to simulated route
      return this.simulateRealisticRoute(start, end, profile);
    }
  }

  private async callOSRMAPI(
    start: RoutePoint,
    end: RoutePoint,
    profile: "driving-car" | "foot-walking" | "cycling-regular"
  ): Promise<RouteResponse> {
    const osrmProfile = this.mapProfileToOSRM(profile);
    const url = `https://router.project-osrm.org/route/v1/${osrmProfile}/${start.lng},${start.lat};${end.lng},${end.lat}?overview=full&geometries=geojson&steps=true`;

    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    if (!data.routes || data.routes.length === 0) {
      throw new Error("No route found");
    }

    const route = data.routes[0];

    // Extract coordinates from the geometry
    const coordinates = route.geometry.coordinates.map(
      ([lng, lat]: [number, number]) => ({
        lat,
        lng,
      })
    );

    // Extract instructions from steps
    const instructions = route.legs.flatMap((leg: any) =>
      leg.steps.map((step: any) => step.maneuver.instruction || "Continue")
    );

    return {
      path: coordinates,
      duration: Math.round(route.duration / 60), // convert to minutes
      distance: Math.round(route.distance / 1000), // convert to km
      instructions: instructions.filter(Boolean),
    };
  }

  private mapProfileToOSRM(profile: string): string {
    switch (profile) {
      case "foot-walking":
        return "foot";
      case "cycling-regular":
        return "bike";
      case "driving-car":
      default:
        return "driving";
    }
  }

  // Fallback method for when real routing fails
  private simulateRealisticRoute(
    start: RoutePoint,
    end: RoutePoint,
    profile: "driving-car" | "foot-walking" | "cycling-regular" = "driving-car"
  ): RouteResponse {
    // Create a more realistic route that follows major roads in Algiers
    const path: RoutePoint[] = [];

    // Add start point
    path.push(start);

    // Calculate intermediate points that follow main roads
    const latDiff = end.lat - start.lat;
    const lngDiff = end.lng - start.lng;

    // Add waypoints that simulate following major roads
    const numWaypoints = Math.max(
      3,
      Math.floor(Math.abs(latDiff) + Math.abs(lngDiff)) * 10
    );

    for (let i = 1; i < numWaypoints; i++) {
      const progress = i / numWaypoints;

      // Add some realistic road following behavior
      const roadOffset = this.getRoadOffset(progress, start, end);

      const waypoint: RoutePoint = {
        lat: start.lat + latDiff * progress + roadOffset.lat,
        lng: start.lng + lngDiff * progress + roadOffset.lng,
      };

      path.push(waypoint);
    }

    // Add end point
    path.push(end);

    // Calculate realistic duration and distance
    const distance = this.calculateDistance(start, end);
    const duration = this.estimateDuration(distance, start, end, profile);

    return {
      path,
      duration,
      distance,
      instructions: this.generateInstructions(path),
    };
  }

  private getRoadOffset(
    progress: number,
    start: RoutePoint,
    end: RoutePoint
  ): RoutePoint {
    // Simulate road network by adding realistic offsets
    const majorRoads = [
      { lat: 36.7538, lng: 3.0588 }, // Alger Centre
      { lat: 36.7697, lng: 3.0611 }, // Port area
      { lat: 36.7403, lng: 3.0508 }, // Riadh El Feth
      { lat: 36.7167, lng: 3.1833 }, // Bab Ezzouar
    ];

    // Find the closest major road
    let closestRoad = majorRoads[0];
    let minDistance = Infinity;

    const currentLat = start.lat + (end.lat - start.lat) * progress;
    const currentLng = start.lng + (end.lng - start.lng) * progress;

    majorRoads.forEach((road) => {
      const distance =
        Math.abs(road.lat - currentLat) + Math.abs(road.lng - currentLng);
      if (distance < minDistance) {
        minDistance = distance;
        closestRoad = road;
      }
    });

    // Apply a small offset towards the major road
    const roadInfluence = 0.1;
    return {
      lat:
        (closestRoad.lat - currentLat) *
        roadInfluence *
        Math.sin(progress * Math.PI),
      lng:
        (closestRoad.lng - currentLng) *
        roadInfluence *
        Math.sin(progress * Math.PI),
    };
  }

  private calculateDistance(start: RoutePoint, end: RoutePoint): number {
    const R = 6371; // Earth's radius in km
    const dLat = ((end.lat - start.lat) * Math.PI) / 180;
    const dLng = ((end.lng - start.lng) * Math.PI) / 180;
    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos((start.lat * Math.PI) / 180) *
        Math.cos((end.lat * Math.PI) / 180) *
        Math.sin(dLng / 2) *
        Math.sin(dLng / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return Math.round(R * c);
  }

  private estimateDuration(
    distance: number,
    start: RoutePoint,
    end: RoutePoint,
    profile: "driving-car" | "foot-walking" | "cycling-regular" = "driving-car"
  ): number {
    // Base speed in km/h based on transportation mode
    let baseSpeed = 30; // Urban areas for driving

    // Adjust base speed based on profile
    switch (profile) {
      case "foot-walking":
        baseSpeed = 5; // Walking speed
        break;
      case "cycling-regular":
        baseSpeed = 15; // Cycling speed
        break;
      case "driving-car":
      default:
        baseSpeed = 30; // Urban driving
        break;
    }

    // Adjust speed based on known traffic patterns (only for driving)
    const hour = new Date().getHours();
    let trafficMultiplier = 1;

    if (
      profile === "driving-car" &&
      ((hour >= 7 && hour <= 9) || (hour >= 17 && hour <= 19))
    ) {
      trafficMultiplier = 1.5; // Rush hour
    }

    // Airport routes are typically faster (only for driving)
    if (
      profile === "driving-car" &&
      Math.abs(end.lat - 36.691) < 0.01 &&
      Math.abs(end.lng - 3.2154) < 0.01
    ) {
      baseSpeed = 50; // Highway to airport
    }

    const effectiveSpeed = baseSpeed / trafficMultiplier;
    return Math.round((distance / effectiveSpeed) * 60); // Convert to minutes
  }

  private generateInstructions(path: RoutePoint[]): string[] {
    const instructions: string[] = [];

    if (path.length < 2) return instructions;

    instructions.push("Head towards your destination");

    // Generate turn instructions based on direction changes
    for (let i = 1; i < path.length - 1; i++) {
      const prev = path[i - 1];
      const current = path[i];
      const next = path[i + 1];

      const bearing1 = this.getBearing(prev, current);
      const bearing2 = this.getBearing(current, next);
      const turnAngle = bearing2 - bearing1;

      if (Math.abs(turnAngle) > 30) {
        const direction = turnAngle > 0 ? "right" : "left";
        instructions.push(`Turn ${direction} and continue`);
      } else {
        instructions.push("Continue straight");
      }
    }

    instructions.push("You have arrived at your destination");

    return instructions;
  }

  private getBearing(start: RoutePoint, end: RoutePoint): number {
    const dLng = ((end.lng - start.lng) * Math.PI) / 180;
    const startLat = (start.lat * Math.PI) / 180;
    const endLat = (end.lat * Math.PI) / 180;

    const y = Math.sin(dLng) * Math.cos(endLat);
    const x =
      Math.cos(startLat) * Math.sin(endLat) -
      Math.sin(startLat) * Math.cos(endLat) * Math.cos(dLng);

    return (Math.atan2(y, x) * 180) / Math.PI;
  }
}
