"use client"; // Client component for interactivity

import Link from "next/link";
import { ReactNode, useState } from "react";

// Props for the DashboardLayout component
interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside
        className={`bg-gradient-to-b from-indigo-900 to-indigo-700 text-white shadow-lg transition-all duration-300 ${isSidebarOpen ? "w-64" : "w-16"
          }`}
      >
        <div className="p-6 flex items-center justify-between border-b border-indigo-600">
          {isSidebarOpen && (
            <h2 className="text-2xl font-semibold tracking-wide">Dashboard</h2>
          )}
          <button
            onClick={toggleSidebar}
            className="p-2 rounded-lg bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-400 transition-colors"
          >
            <div className="text-lg font-bold">
              {isSidebarOpen ? "<<" : ">>"}
            </div>
          </button>
        </div>
        <nav className="mt-4">
          <ul>
            <li>
              <Link
                href="/"
                className={`block p-4 hover:bg-indigo-600 transition-colors ${isSidebarOpen ? "text-left" : "text-center"
                  }`}
              >
                {isSidebarOpen ? "Home" : "H"}
              </Link>
            </li>
            <li>
              <Link
                href="/about"
                className={`block p-4 hover:bg-indigo-600 transition-colors ${isSidebarOpen ? "text-left" : "text-center"
                  }`}
              >
                {isSidebarOpen ? "About" : "A"}
              </Link>
            </li>
            <li>
              <Link
                href="/tenant1/prodA/obj123/viewX"
                className={`block p-4 hover:bg-indigo-600 transition-colors ${isSidebarOpen ? "text-left" : "text-center"
                  }`}
              >
                {isSidebarOpen ? "Dynamic Page" : "D"}
              </Link>
            </li>
          </ul>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-8 bg-gray-50">
        {children}
      </main>
    </div>
  );
}