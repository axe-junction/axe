import { Link, useLocation } from "react-router-dom";
import { Map, Search, Users, Heart, User } from "lucide-react";

const navItems = [
  { path: "/", icon: Map, label: "Map" },
  { path: "/search", icon: Search, label: "Search" },
  { path: "/community", icon: Users, label: "Community" },
  { path: "/favorites", icon: Heart, label: "Favorites" },
  { path: "/profile", icon: User, label: "Profile" },
];

export default function BottomNavigation() {
  const location = useLocation();

  return (
    <nav className="bottom-nav bg-white border-t border-gray-200 lg:hidden">
      <div className="flex justify-around items-center py-2">
        {navItems.map(({ path, icon: Icon, label }) => {
          const isActive = location.pathname === path;
          return (
            <Link
              key={path}
              to={path}
              className={`flex flex-col items-center py-2 px-3 rounded-lg transition-all min-w-0 ${
                isActive
                  ? "text-primary bg-primary/10 scale-105"
                  : "text-gray-500 hover:text-primary hover:bg-primary/5"
              }`}
            >
              <Icon size={22} className={isActive ? "animate-pulse" : ""} />
              <span
                className={`text-xs mt-1 truncate ${
                  isActive ? "font-medium" : ""
                }`}
              >
                {label}
              </span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
