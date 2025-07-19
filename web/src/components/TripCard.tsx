import { useState } from "react";
import {
  Heart,
  Clock,
  MapPin,
  ArrowRight,
  ChevronDown,
  ChevronUp,
} from "lucide-react";
import type { Trip } from "../store/useAppStore";
import { useAppStore } from "../store/useAppStore";

interface TripCardProps {
  trip: Trip;
  onSelect?: (trip: Trip) => void;
}

export default function TripCard({ trip, onSelect }: TripCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const { toggleFavorite } = useAppStore();

  const handleFavoriteToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    toggleFavorite(trip.id);
  };

  const handleCardClick = () => {
    if (onSelect) {
      onSelect(trip);
    }
  };

  return (
    <div
      className="bg-white rounded-lg shadow-md border border-gray-200 p-4 cursor-pointer hover:shadow-lg transition-shadow"
      onClick={handleCardClick}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <MapPin size={16} className="text-gray-500" />
            <span className="text-sm font-medium text-gray-900">
              {trip.from}
            </span>
            <ArrowRight size={16} className="text-gray-400" />
            <span className="text-sm font-medium text-gray-900">{trip.to}</span>
          </div>

          <div className="flex items-center gap-4 mb-2">
            <div className="flex items-center gap-1">
              <Clock size={14} className="text-gray-400" />
              <span className="text-sm text-gray-600">{trip.duration} min</span>
            </div>
            <div className="text-sm font-semibold text-primary">
              {trip.cost} DA
            </div>
          </div>

          <div className="flex items-center gap-2 mb-2">
            <span className="text-xs text-gray-500">Via {trip.provider}</span>
            <span className="text-xs text-gray-400">•</span>
            <span className="text-xs text-gray-500">{trip.stops} stops</span>
          </div>

          <div className="flex gap-1 flex-wrap">
            {trip.transportModes.map((mode, index) => (
              <span
                key={index}
                className="px-2 py-1 bg-gray-100 text-xs rounded-full text-gray-600"
              >
                {mode}
              </span>
            ))}
          </div>
        </div>

        <div className="flex flex-col items-end gap-2">
          <button
            onClick={handleFavoriteToggle}
            className={`p-2 rounded-full transition-colors ${
              trip.isFavorite
                ? "text-red-500 bg-red-50"
                : "text-gray-400 hover:text-red-500 hover:bg-red-50"
            }`}
          >
            <Heart size={18} fill={trip.isFavorite ? "currentColor" : "none"} />
          </button>

          <button
            onClick={(e) => {
              e.stopPropagation();
              setIsExpanded(!isExpanded);
            }}
            className="p-1 rounded-full hover:bg-gray-100 transition-colors"
          >
            {isExpanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </button>
        </div>
      </div>

      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-gray-100">
          <h4 className="text-sm font-medium text-gray-900 mb-2">
            Transfer Points
          </h4>
          <div className="space-y-1">
            {trip.transferPoints.map((point, index) => (
              <div key={index} className="text-sm text-gray-600">
                {index + 1}. {point}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
