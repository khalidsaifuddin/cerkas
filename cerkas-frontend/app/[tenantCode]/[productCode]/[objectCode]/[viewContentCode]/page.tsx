"use client"; // Mark as client component

import { useParams } from "next/navigation";

// Define the expected params type with an index signature
interface RouteParams {
  tenantCode: string;
  productCode: string;
  objectCode: string;
  viewContentCode: string;
  [key: string]: string | string[]; // Index signature for Params compatibility
}

export default function DynamicPage() {
  const params = useParams<RouteParams>();
  const { tenantCode, productCode, objectCode, viewContentCode } = params;

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold text-indigo-700 mb-6">
          Dynamic Route Page
        </h1>
        <div className="space-y-3 text-gray-700">
          <p>
            <span className="font-semibold">Tenant Code:</span>{" "}
            <span className="text-indigo-600">{tenantCode}</span>
          </p>
          <p>
            <span className="font-semibold">Product Code:</span>{" "}
            <span className="text-indigo-600">{productCode}</span>
          </p>
          <p>
            <span className="font-semibold">Object Code:</span>{" "}
            <span className="text-indigo-600">{objectCode}</span>
          </p>
          <p>
            <span className="font-semibold">View Content Code:</span>{" "}
            <span className="text-indigo-600">{viewContentCode}</span>
          </p>
        </div>
      </div>
    </div>
  );
}