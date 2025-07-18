import type { ReactNode } from "react";
import BottomNavigation from "./BottomNavigation";
import Navbar from "./Navbar";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="h-screen flex flex-col bg-white">
      <Navbar />
      <main className="flex-1 overflow-hidden">{children}</main>
      <BottomNavigation />
    </div>
  );
}
