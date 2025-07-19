import type { ReactNode } from "react";
import { useLocation } from "react-router-dom";
import BottomNavigation from "./BottomNavigation";
import Navbar from "./Navbar";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const isHomePage = location.pathname === "/";

  return (
    <div className="page-container">
      <Navbar />
      <main
        className={`content-area ${
          isHomePage ? "non-scrollable" : "scrollable"
        }`}
      >
        {children}
      </main>
      <BottomNavigation />
    </div>
  );
}
