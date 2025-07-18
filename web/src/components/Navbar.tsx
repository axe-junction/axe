import { Link, useLocation } from "react-router-dom";
import { Map, Search, Users, Heart, User, Menu, X } from "lucide-react";
import { useState } from "react";
import { useAppStore } from "../store/useAppStore";

const navItems = [
  { path: "/", icon: Map, label: "Map" },
  { path: "/search", icon: Search, label: "Search" },
  { path: "/community", icon: Users, label: "Community" },
  { path: "/favorites", icon: Heart, label: "Favorites" },
  { path: "/profile", icon: User, label: "Profile" },
];

export default function Navbar() {
  const location = useLocation();
  const { user } = useAppStore();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <>
      {/* Desktop Navbar */}
      <nav
        className="hidden lg:block shadow-lg w-full"
        style={{ backgroundColor: "#6316DB" }}
      > 
        <div className="w-full px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo */}
            <div className="flex items-center">
              <Link to="/" className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center">
                  <Map size={20} style={{ color: "#6316DB" }} />
                </div>
                <span className="text-xl font-bold text-white">Axe</span>
              </Link>
            </div>

            {/* Navigation Links */}
            <div className="flex items-center space-x-2">
              {navItems.map(({ path, icon: Icon, label }) => {
                const isActive = location.pathname === path;
                return (
                  <Link
                    key={path}
                    to={path}
                    className={`flex items-center space-x-2 px-4 py-2 rounded-lg transition-colors ${
                      isActive
                        ? "bg-white shadow-sm"
                        : "text-white hover:bg-white/10"
                    }`}
                    style={isActive ? { color: "#6316DB" } : { color: "white" }}
                  >
                    <Icon size={20} />
                    <span className="font-medium">{label}</span>
                  </Link>
                );
              })}
            </div>

            {/* User Profile */}
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-3 text-white">
                <img
                  src={
                    user?.avatar ||
                    "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&h=150&fit=crop&crop=face"
                  }
                  alt="Profile"
                  className="w-8 h-8 rounded-full object-cover border-2 border-white/20"
                />
                <span className="text-sm font-medium">
                  {user?.name || "Guest"}
                </span>
              </div>
            </div>
          </div>
        </div>
      </nav>

      {/* Mobile Navbar */}
      <nav
        className="lg:hidden shadow-lg w-full"
        style={{ backgroundColor: "#6316DB" }}
      >
        <div className="px-4">
          <div className="flex justify-between items-center h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center">
                <Map size={20} style={{ color: "#6316DB" }} />
              </div>
              <span className="text-xl font-bold text-white">Axe</span>
            </Link>

            {/* Mobile Menu Button */}
            <button
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              className="p-2 rounded-lg text-white hover:bg-white/10 transition-colors"
            >
              {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMobileMenuOpen && (
          <div
            className="border-t border-white/20"
            style={{ backgroundColor: "#6316DB" }}
          >
            <div className="px-4 py-2 space-y-1">
              {navItems.map(({ path, icon: Icon, label }) => {
                const isActive = location.pathname === path;
                return (
                  <Link
                    key={path}
                    to={path}
                    onClick={() => setIsMobileMenuOpen(false)}
                    className={`flex items-center space-x-3 px-3 py-3 rounded-lg transition-colors ${
                      isActive
                        ? "bg-white shadow-sm"
                        : "text-white hover:bg-white/10"
                    }`}
                    style={isActive ? { color: "#6316DB" } : { color: "white" }}
                  >
                    <Icon size={20} />
                    <span className="font-medium">{label}</span>
                  </Link>
                );
              })}

              {/* User Profile in Mobile Menu */}
              <div className="flex items-center space-x-3 px-3 py-3 border-t border-white/20 mt-2">
                <img
                  src={
                    user?.avatar ||
                    "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&h=150&fit=crop&crop=face"
                  }
                  alt="Profile"
                  className="w-8 h-8 rounded-full object-cover border-2 border-white/20"
                />
                <span className="text-sm font-medium text-white">
                  {user?.name || "Guest"}
                </span>
              </div>
            </div>
          </div>
        )}
      </nav>
    </>
  );
}
