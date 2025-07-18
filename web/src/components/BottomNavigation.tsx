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
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 lg:hidden z-40">
      <div className="flex justify-around items-center py-2">
        {navItems.map(({ path, icon: Icon, label }) => {
          const isActive = location.pathname === path;
          return (
            <Link
              key={path}
              to={path}
              className={`flex flex-col items-center py-2 px-3 rounded-lg transition-colors min-w-0 ${
                isActive
                  ? "text-primary bg-primary/10"
                  : "text-gray-500 hover:text-gray-700"
              }`}
            >
              <Icon size={22} />
              <span className="text-xs mt-1 truncate">{label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
