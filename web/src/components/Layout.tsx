import type { ReactNode } from "react";
import BottomNavigation from "./BottomNavigation";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-white">
      <main className="pb-20 lg:pb-0">{children}</main>
      <BottomNavigation />
    </div>
  );
}
