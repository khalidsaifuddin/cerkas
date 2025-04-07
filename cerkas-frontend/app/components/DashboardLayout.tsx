"use client"; // Client component for interactivity

import Link from "next/link";
import { ReactNode, useState } from "react";
import { dashboardConfig } from "../appConfig";

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
        className={`bg-gradient-to-b from-teal-900 to-teal-700 text-white shadow-lg transition-all duration-300 ${isSidebarOpen ? "w-48" : "w-14"
          }`}
      >
        <div className="p-2 pl-4 flex items-center justify-between border-b border-teal-600">
          {isSidebarOpen && (
            <h2 className="text-2xl font-semibold tracking-wide">{dashboardConfig.title}</h2>
          )}
          <button
            onClick={toggleSidebar}
            className="p-2 rounded-lg bg-teal-600 hover:bg-teal-500 focus:outline-none focus:ring-2 focus:ring-teal-400 transition-colors"
          >
            <div className="text-lg font-bold">
              {isSidebarOpen ? "<" : ">"}
            </div>
          </button>
        </div>
        <nav className="mt-4">
          <ul>
            <li>
              <Link
                href="/"
                className={`block p-4 hover:bg-teal-600 transition-colors ${isSidebarOpen ? "text-left" : "text-center"
                  }`}
              >
                {isSidebarOpen ? "Home" : "H"}
              </Link>
            </li>
            <li>
              <Link
                href="/about"
                className={`block p-4 hover:bg-teal-600 transition-colors ${isSidebarOpen ? "text-left" : "text-center"
                  }`}
              >
                {isSidebarOpen ? "About" : "A"}
              </Link>
            </li>
            <li>
              <Link
                href="/fetchly/delivery/user/default"
                className={`block p-4 hover:bg-teal-600 transition-colors ${isSidebarOpen ? "text-left" : "text-center"
                  }`}
              >
                {isSidebarOpen ? "Dynamic Page" : "D"}
              </Link>
            </li>
          </ul>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="flex-1 bg-gray-50 h-screen flex flex-col">
        {/* Optional fixed header inside main */}
        <div className="p-4 border-b bg-gradient-to-r from-teal-900 to-teal-700 text-white shadow z-10">
          <h1 className="text-xl font-semibold">Page Title</h1>
        </div>

        {/* Scrollable content area */}
        <div className="flex-1 overflow-auto">
          {children}
        </div>
      </main>
    </div>
  );
}