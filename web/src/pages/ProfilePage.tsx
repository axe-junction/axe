import { useState } from "react";
import {
  User,
  Mail,
  Phone,
  MapPin,
  Edit3,
  Save,
  X,
  Settings,
  Bell,
  Shield,
  HelpCircle,
  LogOut,
  Heart,
  Clock,
  Navigation,
} from "lucide-react";
import { useAppStore } from "../store/useAppStore";

export default function ProfilePage() {
  const { user, setUser, favoriteTrips, recentSearches } = useAppStore();
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({
    name: user?.name || "John Doe",
    email: user?.email || "john.doe@example.com",
    phone: "+213 555 0123",
    location: "Algiers, Algeria",
  });

  const handleSave = () => {
    setUser({
      id: user?.id || "1",
      name: editForm.name,
      email: editForm.email,
      avatar:
        user?.avatar ||
        "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&h=150&fit=crop&crop=face",
    });
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditForm({
      name: user?.name || "John Doe",
      email: user?.email || "john.doe@example.com",
      phone: "+213 555 0123",
      location: "Algiers, Algeria",
    });
    setIsEditing(false);
  };

  const menuItems = [
    { icon: Settings, label: "Settings", action: () => {} },
    { icon: Bell, label: "Notifications", action: () => {} },
    { icon: Shield, label: "Privacy", action: () => {} },
    { icon: HelpCircle, label: "Help & Support", action: () => {} },
    { icon: LogOut, label: "Sign Out", action: () => {}, danger: true },
  ];

  return (
    <div className="min-h-full bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="p-4">
          <h1 className="text-2xl font-bold text-gray-900">Profile</h1>
        </div>
      </div>

      {/* Profile Card */}
      <div className="p-4">
        <div className="bg-white rounded-lg shadow-sm border p-6">
          <div className="flex items-start justify-between mb-6">
            <div className="flex items-center gap-4">
              <img
                src={
                  user?.avatar ||
                  "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&h=150&fit=crop&crop=face"
                }
                alt="Profile"
                className="w-20 h-20 rounded-full object-cover"
              />
              <div>
                <h2 className="text-xl font-semibold text-gray-900">
                  {user?.name || "John Doe"}
                </h2>
                <p className="text-gray-600">
                  {user?.email || "john.doe@example.com"}
                </p>
              </div>
            </div>
            <button
              onClick={() => setIsEditing(true)}
              className="p-2 text-gray-600 hover:text-primary hover:bg-primary/10 rounded-full transition-colors"
            >
              <Edit3 size={20} />
            </button>
          </div>

          {/* Profile Details */}
          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <User size={16} className="text-gray-400" />
              <span className="text-gray-900">{editForm.name}</span>
            </div>
            <div className="flex items-center gap-3">
              <Mail size={16} className="text-gray-400" />
              <span className="text-gray-900">{editForm.email}</span>
            </div>
            <div className="flex items-center gap-3">
              <Phone size={16} className="text-gray-400" />
              <span className="text-gray-900">{editForm.phone}</span>
            </div>
            <div className="flex items-center gap-3">
              <MapPin size={16} className="text-gray-400" />
              <span className="text-gray-900">{editForm.location}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="px-4 pb-4">
        <div className="grid grid-cols-3 gap-4">
          <div className="bg-white rounded-lg shadow-sm border p-4 text-center">
            <Heart size={24} className="text-red-500 mx-auto mb-2" />
            <div className="text-lg font-semibold text-gray-900">
              {favoriteTrips.length}
            </div>
            <div className="text-sm text-gray-600">Favorites</div>
          </div>
          <div className="bg-white rounded-lg shadow-sm border p-4 text-center">
            <Clock size={24} className="text-blue-500 mx-auto mb-2" />
            <div className="text-lg font-semibold text-gray-900">
              {recentSearches.length}
            </div>
            <div className="text-sm text-gray-600">Recent</div>
          </div>
          <div className="bg-white rounded-lg shadow-sm border p-4 text-center">
            <Navigation size={24} className="text-green-500 mx-auto mb-2" />
            <div className="text-lg font-semibold text-gray-900">128</div>
            <div className="text-sm text-gray-600">Total Trips</div>
          </div>
        </div>
      </div>

      {/* Menu */}
      <div className="px-4 pb-4">
        <div className="bg-white rounded-lg shadow-sm border">
          {menuItems.map((item, index) => (
            <button
              key={index}
              onClick={item.action}
              className={`w-full flex items-center gap-3 p-4 hover:bg-gray-50 transition-colors ${
                index !== menuItems.length - 1 ? "border-b border-gray-100" : ""
              } ${
                index === 0
                  ? "rounded-t-lg"
                  : index === menuItems.length - 1
                  ? "rounded-b-lg"
                  : ""
              }`}
            >
              <item.icon
                size={20}
                className={item.danger ? "text-red-500" : "text-gray-400"}
              />
              <span
                className={`flex-1 text-left ${
                  item.danger ? "text-red-500" : "text-gray-900"
                }`}
              >
                {item.label}
              </span>
            </button>
          ))}
        </div>
      </div>

      {/* Edit Modal */}
      {isEditing && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg w-full max-w-md">
            <div className="p-4 border-b flex items-center justify-between">
              <h2 className="text-lg font-semibold">Edit Profile</h2>
              <button
                onClick={handleCancel}
                className="p-2 hover:bg-gray-100 rounded-full"
              >
                <X size={20} />
              </button>
            </div>

            <div className="p-4 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Name
                </label>
                <input
                  type="text"
                  value={editForm.name}
                  onChange={(e) =>
                    setEditForm({ ...editForm, name: e.target.value })
                  }
                  className="w-full p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  value={editForm.email}
                  onChange={(e) =>
                    setEditForm({ ...editForm, email: e.target.value })
                  }
                  className="w-full p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Phone
                </label>
                <input
                  type="tel"
                  value={editForm.phone}
                  onChange={(e) =>
                    setEditForm({ ...editForm, phone: e.target.value })
                  }
                  className="w-full p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Location
                </label>
                <input
                  type="text"
                  value={editForm.location}
                  onChange={(e) =>
                    setEditForm({ ...editForm, location: e.target.value })
                  }
                  className="w-full p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                />
              </div>
            </div>

            <div className="p-4 border-t flex gap-2">
              <button
                onClick={handleCancel}
                className="flex-1 py-2 px-4 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                className="flex-1 py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors flex items-center justify-center gap-2"
              >
                <Save size={16} />
                Save
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
